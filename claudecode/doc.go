// Package claudecode provides a Go SDK for interacting with Claude Code CLI.
//
// This package allows you to programmatically interact with Claude through
// the Claude Code CLI, sending prompts and receiving responses including
// text, tool usage, and results.
//
// Basic usage:
//
//	ctx := context.Background()
//	ch := claudecode.Query(ctx, "What's 2+2?", nil)
//
//	for result := range ch {
//	    if result.Error != nil {
//	        log.Fatal(result.Error)
//	    }
//	    // Process result.Message
//	}
//
// With options:
//
//	options := &claudecode.ClaudeCodeOptions{
//	    AllowedTools: []string{"Edit", "Read"},
//	    Model: stringPtr("claude-3-opus-20240229"),
//	}
//	ch := claudecode.Query(ctx, prompt, options)
//
// For a simpler interface that collects all messages:
//
//	result, messages, err := claudecode.QuerySimple(ctx, prompt, options)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(result.Content)
package claudecode