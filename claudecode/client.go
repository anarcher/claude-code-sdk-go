package claudecode

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

// InternalClient handles message processing and parsing
type InternalClient struct {
	transport *Transport
	options   *ClaudeCodeOptions
}

// NewInternalClient creates a new internal client
func NewInternalClient(ctx context.Context, options *ClaudeCodeOptions) (*InternalClient, error) {
	if options == nil {
		options = DefaultOptions()
	}
	
	// Find CLI
	cliPath, err := findCLI()
	if err != nil {
		return nil, err
	}
	
	// Build CLI arguments
	args := buildCLIArgs(options)
	
	// Create transport
	transport, err := NewTransport(ctx, cliPath, args)
	if err != nil {
		return nil, err
	}
	
	return &InternalClient{
		transport: transport,
		options:   options,
	}, nil
}

// buildCLIArgs builds command line arguments from options
func buildCLIArgs(options *ClaudeCodeOptions) []string {
	var args []string
	
	// Add streaming JSON output for programmatic use and skip permissions
	args = append(args, "--verbose", "--output-format", "stream-json", "--dangerously-skip-permissions")
	
	// Add options
	if len(options.AllowedTools) > 0 {
		args = append(args, "--allowedTools", strings.Join(options.AllowedTools, ","))
	}
	
	// Note: MaxThinkingTokens is not supported by the current CLI version
	// if options.MaxThinkingTokens != nil {
	//     args = append(args, "--max-thinking-tokens", fmt.Sprintf("%d", *options.MaxThinkingTokens))
	// }
	
	// Note: CLI only supports --append-system-prompt, not --system-prompt
	// if options.SystemPrompt != nil {
	//     args = append(args, "--system-prompt", *options.SystemPrompt)
	// }
	
	if options.AppendSystemPrompt != nil {
		args = append(args, "--append-system-prompt", *options.AppendSystemPrompt)
	}
	
	if options.PermissionMode != nil {
		args = append(args, "--permission-mode", string(*options.PermissionMode))
	}
	
	if options.ContinueConversation {
		args = append(args, "--continue")
	}
	
	if options.Resume != nil {
		args = append(args, "--resume", *options.Resume)
	}
	
	if options.MaxTurns != nil {
		args = append(args, "--max-turns", fmt.Sprintf("%d", *options.MaxTurns))
	}
	
	if len(options.DisallowedTools) > 0 {
		args = append(args, "--disallowedTools", strings.Join(options.DisallowedTools, ","))
	}
	
	if options.Model != nil {
		args = append(args, "--model", *options.Model)
	}
	
	if options.PermissionPromptToolName != nil {
		args = append(args, "--permission-prompt-tool-name", *options.PermissionPromptToolName)
	}
	
	if options.CWD != nil {
		args = append(args, "--cwd", *options.CWD)
	}
	
	// Handle MCP servers
	for _, server := range options.MCPServers {
		serverJSON, err := json.Marshal(server)
		if err != nil {
			continue
		}
		args = append(args, "--mcp-server", string(serverJSON))
	}
	
	// Handle MCP tools
	for _, tool := range options.MCPTools {
		args = append(args, "--mcp-tool", string(tool))
	}
	
	return args
}

// SendPrompt sends a prompt to the CLI
func (c *InternalClient) SendPrompt(prompt string) error {
	return c.transport.Send(prompt)
}

// ReceiveMessage receives and parses the next message
func (c *InternalClient) ReceiveMessage() (Message, error) {
	raw, err := c.transport.Receive()
	if err != nil {
		return nil, err
	}
	
	// Parse the message type first
	var msgType struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(raw, &msgType); err != nil {
		return nil, &ParseError{Message: "failed to parse message type", Data: string(raw)}
	}
	
	// Parse based on type
	switch MessageType(msgType.Type) {
	case MessageTypeUser:
		var msg UserMessage
		if err := json.Unmarshal(raw, &msg); err != nil {
			return nil, &ParseError{Message: "failed to parse user message", Data: string(raw)}
		}
		return msg, nil
		
	case MessageTypeAssistant:
		var msg AssistantMessage
		if err := json.Unmarshal(raw, &msg); err != nil {
			return nil, &ParseError{Message: "failed to parse assistant message", Data: string(raw)}
		}
		return msg, nil
		
	case MessageTypeSystem:
		var msg SystemMessage
		if err := json.Unmarshal(raw, &msg); err != nil {
			return nil, &ParseError{Message: "failed to parse system message", Data: string(raw)}
		}
		return msg, nil
		
	case MessageTypeResult:
		var msg ResultMessage
		if err := json.Unmarshal(raw, &msg); err != nil {
			return nil, &ParseError{Message: "failed to parse result message", Data: string(raw)}
		}
		return msg, nil
		
	default:
		return nil, &ParseError{Message: fmt.Sprintf("unknown message type: %s", msgType.Type), Data: string(raw)}
	}
}

// Close closes the client and cleans up resources
func (c *InternalClient) Close() error {
	return c.transport.Close()
}