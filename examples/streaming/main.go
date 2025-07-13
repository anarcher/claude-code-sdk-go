package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	claudecode "github.com/anarcher/claude-code-sdk-go"
)

func main() {
	// Create context that cancels on interrupt
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle interrupts gracefully
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Println("\nReceived interrupt, shutting down...")
		cancel()
	}()

	// Use empty options for simple text response
	options := &claudecode.ClaudeCodeOptions{}

	prompt := "Explain the concept of recursion in programming"

	fmt.Println("Starting streaming query...")
	ch := claudecode.Query(ctx, prompt, options)

	// Stream messages as they arrive
	for result := range ch {
		if result.Error != nil {
			log.Fatalf("Error: %v", result.Error)
		}
		
		msg := result.Message

		// Process different message types
		switch m := msg.(type) {
		case claudecode.AssistantMessage:
			// Stream assistant responses
			for _, rawBlock := range m.Content() {
				block, err := claudecode.ParseContentBlock(rawBlock)
				if err != nil {
					continue
				}
				switch b := block.(type) {
				case claudecode.TextBlock:
					fmt.Print(b.Text)
				case claudecode.ToolUseBlock:
					fmt.Printf("\n[Calling %s...]\n", b.Name)
				case claudecode.ToolResultBlock:
					if b.IsError {
						fmt.Printf("[Tool error]\n")
					} else if b.Output != nil {
						// Tool output - could be file contents, command output, etc.
						fmt.Printf("[Tool completed]\n")
					}
				}
			}

		case claudecode.SystemMessage:
			// System messages (like thinking, etc.)
			if m.Subtype == "thinking" {
				fmt.Print("🤔 ")
			}

		case claudecode.ResultMessage:
			// Final result
			fmt.Printf("\n\n✅ Complete!\n")
			if m.Cost != nil {
				fmt.Printf("Cost: $%.4f\n", m.Cost.TotalCost)
			}
		}
	}
}