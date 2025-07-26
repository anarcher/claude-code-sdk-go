package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/anarcher/claude-code-sdk-go/claudecode"
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

	// Use options to allow tools for more verbose output
	options := &claudecode.ClaudeCodeOptions{
		AllowedTools: []string{"Read", "Write", "Bash"},
	}

	prompt := "What is 2+2? Think step by step."

	fmt.Println("Starting query with verbose output...")
	fmt.Println("==================================================")
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
			fmt.Println("\n[ASSISTANT MESSAGE]")
			// Stream assistant responses
			for _, rawBlock := range m.Content() {
				block, err := claudecode.ParseContentBlock(rawBlock)
				if err != nil {
					continue
				}
				switch b := block.(type) {
				case claudecode.TextBlock:
					fmt.Printf("  Text: %s", b.Text)
				case claudecode.ToolUseBlock:
					fmt.Printf("\n  🔧 Tool Use: %s\n", b.Name)
					// Pretty print the input
					var prettyInput interface{}
					if err := json.Unmarshal(b.Input, &prettyInput); err == nil {
						inputJSON, _ := json.MarshalIndent(prettyInput, "    ", "  ")
						fmt.Printf("    Input: %s\n", string(inputJSON))
					}
				case claudecode.ToolResultBlock:
					fmt.Printf("\n  📋 Tool Result (ID: %s)\n", b.ToolUseID)
					if b.IsError {
						fmt.Printf("    ❌ Error occurred\n")
					} else if b.Output != nil {
						fmt.Printf("    ✅ Output: %s\n", *b.Output)
					}
					if b.StatusCode != nil {
						fmt.Printf("    Status Code: %d\n", *b.StatusCode)
					}
				}
			}

		case claudecode.SystemMessage:
			fmt.Printf("\n[SYSTEM MESSAGE - %s]\n", m.Subtype)
			
			// Pretty print the data if it exists
			if len(m.Data) > 0 {
				var prettyData interface{}
				if err := json.Unmarshal(m.Data, &prettyData); err == nil {
					dataJSON, _ := json.MarshalIndent(prettyData, "  ", "  ")
					fmt.Printf("  Data: %s\n", string(dataJSON))
				} else {
					fmt.Printf("  Raw Data: %s\n", string(m.Data))
				}
			}

		case claudecode.ResultMessage:
			fmt.Printf("\n[RESULT MESSAGE]\n")
			fmt.Printf("  Content: %s\n", m.Content)
			if m.Cost != nil {
				fmt.Printf("  💰 Cost: $%.4f\n", m.Cost.TotalCost)
				fmt.Printf("     - Input (cached): %d tokens ($%.4f)\n", m.Cost.InputCached, m.Cost.InputCachedCost)
				fmt.Printf("     - Input (uncached): %d tokens ($%.4f)\n", m.Cost.InputUncached, m.Cost.InputUncachedCost)
				fmt.Printf("     - Output: %d tokens ($%.4f)\n", m.Cost.Output, m.Cost.OutputCost)
			}
			if m.Usage != nil {
				fmt.Printf("  📊 Token Usage:\n")
				fmt.Printf("     - Total: %d\n", m.Usage.TotalTokens)
				fmt.Printf("     - Input: %d\n", m.Usage.InputTokens)
				fmt.Printf("     - Output: %d\n", m.Usage.OutputTokens)
			}
		}
	}
	
	fmt.Println("\n==================================================")
	fmt.Println("Query completed!")
}