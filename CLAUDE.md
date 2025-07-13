# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go SDK for programmatically interacting with Claude through the Claude Code CLI. It provides streaming and simple APIs for sending prompts and receiving structured responses including text, tool usage, and metadata.

## Architecture

The SDK uses a layered architecture with subprocess communication:

- **Public API** (`sdk.go`): `Query()` for streaming, `QuerySimple()` for collected responses
- **Client Layer** (`client.go`): Message processing and CLI argument building
- **Transport Layer** (`transport.go`): Subprocess management and pipe communication
- **Type System** (`types.go`): Message types, content blocks, and JSON parsing
- **Configuration** (`options.go`): Comprehensive options for Claude interaction

Key pattern: The SDK spawns the Claude CLI as a subprocess, communicates via stdin/stdout pipes using JSON streaming protocol, and provides Go channels for async message handling.

## Development Commands

```bash
# Build and test
go build ./...
go test ./...
go test -v ./types_test.go ./options_test.go ./errors_test.go

# Run examples (requires Claude CLI installed)
go run examples/basic/main.go
go run examples/streaming/main.go  
go run examples/with-options/main.go

# Code quality
go fmt ./...
go vet ./...
go mod tidy
```

## CLI Integration Requirements

The SDK requires Claude Code CLI to be installed and accessible. It auto-discovers the CLI in standard locations (`claude` in PATH, `/usr/local/bin/claude`, `/opt/homebrew/bin/claude`) or via `CLAUDE_CLI_PATH` environment variable.

**Critical Implementation Details:**
- Uses `--print --verbose --output-format stream-json --dangerously-skip-permissions` CLI flags
- Sends plain text prompts to stdin (not JSON)
- Must close stdin after sending prompt to signal completion
- Parses streaming JSON responses from stdout

## Message Flow Architecture

1. **Input**: Text prompt sent to CLI stdin
2. **Output**: Three JSON message types received from CLI stdout:
   - `{"type":"system","subtype":"init",...}` - Session initialization
   - `{"type":"assistant","message":{"content":[...],...},...}` - Assistant responses with nested content
   - `{"type":"result",...}` - Final result with cost/usage data

**Important**: `AssistantMessage.Content()` method provides backward compatibility for accessing content blocks, as the CLI nests content inside a `message` field.

## Configuration Limitations

Some options in `ClaudeCodeOptions` are not supported by the current CLI:
- `MaxThinkingTokens` - CLI flag doesn't exist
- `SystemPrompt` - Only `AppendSystemPrompt` is supported (`--append-system-prompt`)

## Error Handling

The SDK provides typed errors:
- `CLIError` - CLI process errors
- `ParseError` - JSON parsing failures  
- `TransportError` - Subprocess communication issues
- `ValidationError` - Input validation errors

## Testing Approach

Tests cover individual components (`types_test.go`, `options_test.go`, `errors_test.go`). Integration testing requires Claude CLI installation and valid API credentials. Examples serve as integration tests.