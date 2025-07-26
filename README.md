# Claude Code SDK for Go

A Go SDK for interacting with Claude through the Claude Code CLI. This SDK provides a programmatic interface to send prompts to Claude and receive structured responses.

## Installation

```bash
go get github.com/anarcher/claude-code-sdk-go/claudecode
```

## Requirements

- Go 1.21 or later
- Claude Code CLI installed and accessible in your PATH

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/anarcher/claude-code-sdk-go/claudecode"
)

func main() {
    ctx := context.Background()
    
    // Send a simple query
    ch := claudecode.Query(ctx, "What's 2+2?", nil)
    
    // Process messages as they arrive
    for result := range ch {
        if result.Error != nil {
            log.Fatal(result.Error)
        }
        
        // Handle different message types
        switch result.Message.Type() {
        case "assistant":
            // Process assistant response
            if msg, ok := result.Message.(claudecode.AssistantMessage); ok {
                for _, block := range msg.Content() {
                    if text, ok := block.(claudecode.TextBlock); ok {
                        fmt.Println(text.Text)
                    }
                }
            }
        case "result":
            if msg, ok := result.Message.(claudecode.ResultMessage); ok {
                fmt.Printf("Result: %s\n", msg.Content)
            }
        }
    }
}
```

## Advanced Usage

### With Options

```go
options := &claudecode.ClaudeCodeOptions{
    AllowedTools:      []string{"Edit", "Read", "Write"},
    MaxThinkingTokens: intPtr(10000),
    Model:             stringPtr("claude-3-opus-20240229"),
    PermissionMode:    permissionModePtr(claudecode.PermissionModeAuto),
}

ch := claudecode.Query(ctx, "Write a hello world program", options)
```

### Simple Interface

For a simpler interface that collects all messages:

```go
result, messages, err := claudecode.QuerySimple(ctx, prompt, options)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Final result: %s\n", result.Content)
fmt.Printf("Total cost: $%.4f\n", result.Cost.TotalCost)
```

### Streaming Responses

```go
ch := claudecode.Query(ctx, prompt, nil)

for result := range ch {
    if result.Error != nil {
        log.Fatal(result.Error)
    }
    
    // Process message in real-time
    switch result.Message.Type() {
    case "assistant":
        // Handle assistant message
    case "system":
        // Handle system message
    case "result":
        // Handle result message
    }
}
```

## Message Types

The SDK supports four main message types:

- `UserMessage`: Messages from the user
- `AssistantMessage`: Messages from Claude containing content blocks
- `SystemMessage`: System messages with metadata
- `ResultMessage`: Final result with cost and usage information

## Content Blocks

Assistant messages contain content blocks:

- `TextBlock`: Plain text content
- `ToolUseBlock`: Tool invocation details
- `ToolResultBlock`: Results from tool execution

## Configuration Options

- `AllowedTools`: List of tools Claude can use
- `MaxThinkingTokens`: Maximum tokens for thinking (default: 8000)
- `SystemPrompt`: Custom system prompt
- `Model`: Specific model to use
- `PermissionMode`: How to handle tool permissions ("ask" or "auto")
- `CWD`: Working directory for tool execution
- And more...

## Error Handling

The SDK provides typed errors for different scenarios:

```go
switch err := err.(type) {
case *claudecode.CLIError:
    // Handle CLI errors
case *claudecode.ParseError:
    // Handle parsing errors
case *claudecode.TransportError:
    // Handle transport errors
}
```

## Examples

See the `examples/` directory for more detailed examples:

- `basic/`: Simple query example
- `with-options/`: Using configuration options
- `streaming/`: Real-time streaming of responses

## License

MIT License - see LICENSE file for details.