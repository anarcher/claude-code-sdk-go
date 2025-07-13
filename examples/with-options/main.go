package main

import (
	"context"
	"fmt"
	"log"

	claudecode "github.com/anarcher/claude-code-sdk-go"
)

func main() {
	ctx := context.Background()

	// Configure options
	maxThinking := 10000
	systemPrompt := "You are a helpful coding assistant."
	model := "claude-3-opus-20240229"
	cwd := "/tmp"

	options := &claudecode.ClaudeCodeOptions{
		AllowedTools:      []string{"Edit", "Read", "Write", "Bash"},
		MaxThinkingTokens: &maxThinking,
		SystemPrompt:      &systemPrompt,
		Model:             &model,
		CWD:               &cwd,
		DisallowedTools:   []string{"WebSearch"}, // Disable web search
	}

	prompt := "Create a simple Python script that generates the Fibonacci sequence"

	// Use QuerySimple for a simpler interface
	fmt.Println("Sending query with custom options...")
	result, messages, err := claudecode.QuerySimple(ctx, prompt, options)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	// Process all messages
	fmt.Printf("Received %d messages\n\n", len(messages))

	for _, msg := range messages {
		switch m := msg.(type) {
		case claudecode.AssistantMessage:
			for _, rawBlock := range m.Content {
				block, err := claudecode.ParseContentBlock(rawBlock)
				if err != nil {
					continue
				}
				switch b := block.(type) {
				case claudecode.TextBlock:
					fmt.Println(b.Text)
				case claudecode.ToolUseBlock:
					fmt.Printf("\n[Using tool: %s]\n", b.Name)
				}
			}
		}
	}

	// Show final result
	fmt.Printf("\n--- Final Result ---\n%s\n", result.Content)
	if result.Cost != nil {
		fmt.Printf("\nTotal cost: $%.4f\n", result.Cost.TotalCost)
	}
}