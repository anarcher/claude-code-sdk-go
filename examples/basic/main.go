package main

import (
	"context"
	"fmt"
	"log"
	"time"

	claudecode "github.com/anarcher/claude-code-sdk-go"
)

func main() {
	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Query with empty options (no defaults)
	fmt.Println("Sending query to Claude...")
	ch := claudecode.Query(ctx, "What's 2+2?", &claudecode.ClaudeCodeOptions{})

	// Process messages as they arrive
	for result := range ch {
		if result.Error != nil {
			log.Fatalf("Error: %v", result.Error)
		}

		msg := result.Message
		switch m := msg.(type) {
		case claudecode.UserMessage:
			fmt.Printf("User: %s\n", m.Content)

		case claudecode.AssistantMessage:
			fmt.Print("Assistant: ")
			for _, rawBlock := range m.Content() {
				block, err := claudecode.ParseContentBlock(rawBlock)
				if err != nil {
					continue
				}
				switch b := block.(type) {
				case claudecode.TextBlock:
					fmt.Print(b.Text)
				case claudecode.ToolUseBlock:
					fmt.Printf("\n[Tool Use: %s]", b.Name)
				case claudecode.ToolResultBlock:
					fmt.Printf("\n[Tool Result: %s]", b.ToolUseID)
				}
			}
			fmt.Println()

		case claudecode.SystemMessage:
			fmt.Printf("System (%s): %s\n", m.Subtype, string(m.Data))

		case claudecode.ResultMessage:
			fmt.Printf("\nResult: %s\n", m.Content)
			if m.Cost != nil {
				fmt.Printf("Cost: $%.4f\n", m.Cost.TotalCost)
			}
			if m.Usage != nil {
				fmt.Printf("Tokens: %d total (%d input, %d output)\n",
					m.Usage.TotalTokens, m.Usage.InputTokens, m.Usage.OutputTokens)
			}
		}
	}
}