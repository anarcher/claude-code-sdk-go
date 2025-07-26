package main

import (
	"context"
	"fmt"
	"log"

	"github.com/anarcher/claude-code-sdk-go/claudecode"
)

func main() {
	ctx := context.Background()

	// Configure options (using only supported CLI options)
	appendPrompt := "Please be concise and clear in your explanations."

	options := &claudecode.ClaudeCodeOptions{
		AppendSystemPrompt: &appendPrompt,
	}

	prompt := "Explain how binary search works"

	// Use QuerySimple for a simpler interface
	fmt.Println("Sending query with custom options...")
	result, messages, err := claudecode.QuerySimple(ctx, prompt, options)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	// Process all messages
	fmt.Printf("Received %d messages\n\n", len(messages))

	for _, msg := range messages {
		switch msg.Type() {
		case claudecode.MessageTypeAssistant:
			m := msg.(*claudecode.AssistantMessage)
			for _, rawBlock := range m.Content() {
				block, err := claudecode.ParseContentBlock(rawBlock)
				if err != nil {
					continue
				}
				switch b := block.(type) {
				case *claudecode.TextBlock:
					fmt.Println(b.Text)
				case *claudecode.ToolUseBlock:
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