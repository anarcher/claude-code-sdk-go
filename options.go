package claudecode

import "encoding/json"

// PermissionMode controls how tools are executed
type PermissionMode string

const (
	PermissionModeAsk  PermissionMode = "ask"
	PermissionModeAuto PermissionMode = "auto"
)

// ClaudeCodeOptions represents configuration options for the Claude Code SDK
type ClaudeCodeOptions struct {
	// AllowedTools is a list of tools that are allowed to be used
	AllowedTools []string `json:"allowed_tools,omitempty"`

	// MaxThinkingTokens is the maximum number of thinking tokens (default: 8000)
	MaxThinkingTokens *int `json:"max_thinking_tokens,omitempty"`

	// SystemPrompt is a custom system prompt to use
	SystemPrompt *string `json:"system_prompt,omitempty"`

	// AppendSystemPrompt is additional content to append to the system prompt
	AppendSystemPrompt *string `json:"append_system_prompt,omitempty"`

	// MCPTools is a list of MCP tools
	MCPTools []json.RawMessage `json:"mcp_tools,omitempty"`

	// MCPServers is a list of MCP server configurations
	MCPServers []MCPServerConfig `json:"mcp_servers,omitempty"`

	// PermissionMode controls how tools are executed (ask or auto)
	PermissionMode *PermissionMode `json:"permission_mode,omitempty"`

	// ContinueConversation continues a previous conversation
	ContinueConversation bool `json:"continue_conversation,omitempty"`

	// Resume continues from a specific session ID
	Resume *string `json:"resume,omitempty"`

	// MaxTurns limits the number of conversation turns
	MaxTurns *int `json:"max_turns,omitempty"`

	// DisallowedTools is a list of tools that are not allowed to be used
	DisallowedTools []string `json:"disallowed_tools,omitempty"`

	// Model specifies which model to use
	Model *string `json:"model,omitempty"`

	// PermissionPromptToolName is the name of the tool to use for permission prompts
	PermissionPromptToolName *string `json:"permission_prompt_tool_name,omitempty"`

	// CWD is the working directory
	CWD *string `json:"cwd,omitempty"`
}

// MCPServerConfig represents an MCP server configuration
type MCPServerConfig struct {
	// Type specifies the server type (stdio, sse, or http)
	Type MCPServerType `json:"type"`

	// StdioConfig is used when Type is "stdio"
	StdioConfig *MCPStdioConfig `json:"stdio_config,omitempty"`

	// SSEConfig is used when Type is "sse"
	SSEConfig *MCPSSEConfig `json:"sse_config,omitempty"`

	// HTTPConfig is used when Type is "http"
	HTTPConfig *MCPHTTPConfig `json:"http_config,omitempty"`
}

// MCPServerType represents the type of MCP server
type MCPServerType string

const (
	MCPServerTypeStdio MCPServerType = "stdio"
	MCPServerTypeSSE   MCPServerType = "sse"
	MCPServerTypeHTTP  MCPServerType = "http"
)

// MCPStdioConfig represents configuration for stdio MCP servers
type MCPStdioConfig struct {
	Command string            `json:"command"`
	Args    []string          `json:"args,omitempty"`
	Env     map[string]string `json:"env,omitempty"`
}

// MCPSSEConfig represents configuration for SSE MCP servers
type MCPSSEConfig struct {
	URL       string             `json:"url"`
	APIKey    *string            `json:"api_key,omitempty"`
	Headers   map[string]string  `json:"headers,omitempty"`
	Transports []MCPTransportType `json:"transports,omitempty"`
}

// MCPHTTPConfig represents configuration for HTTP MCP servers
type MCPHTTPConfig struct {
	URL       string             `json:"url"`
	APIKey    *string            `json:"api_key,omitempty"`
	Headers   map[string]string  `json:"headers,omitempty"`
	Transports []MCPTransportType `json:"transports,omitempty"`
}

// MCPTransportType represents the transport type for MCP servers
type MCPTransportType string

const (
	MCPTransportTypeHTTP MCPTransportType = "http"
	MCPTransportTypeSSE  MCPTransportType = "sse"
)

// DefaultOptions returns a new ClaudeCodeOptions with default values
func DefaultOptions() *ClaudeCodeOptions {
	maxThinking := 8000
	return &ClaudeCodeOptions{
		MaxThinkingTokens: &maxThinking,
	}
}