package claudecode

import (
	"context"
	"io"
)

// MessageChannel represents a channel that yields messages
type MessageChannel <-chan MessageResult

// MessageResult wraps a message with a potential error
type MessageResult struct {
	Message Message
	Error   error
}

// Query sends a prompt to Claude Code and returns a channel that yields messages
func Query(ctx context.Context, prompt string, options *ClaudeCodeOptions) MessageChannel {
	ch := make(chan MessageResult)
	
	go func() {
		defer close(ch)
		
		// Create client
		client, err := NewInternalClient(ctx, options)
		if err != nil {
			ch <- MessageResult{Error: err}
			return
		}
		defer client.Close()
		
		// Send prompt
		if err := client.SendPrompt(prompt); err != nil {
			ch <- MessageResult{Error: err}
			return
		}
		
		// Receive messages until done
		for {
			msg, err := client.ReceiveMessage()
			if err != nil {
				if err == io.EOF {
					// Normal termination
					return
				}
				ch <- MessageResult{Error: err}
				return
			}
			
			// Send message
			select {
			case ch <- MessageResult{Message: msg}:
				// Check if this was a result message (terminal)
				if msg.Type() == MessageTypeResult {
					return
				}
			case <-ctx.Done():
				ch <- MessageResult{Error: ctx.Err()}
				return
			}
		}
	}()
	
	return ch
}

// QuerySimple is a simplified version that collects all messages and returns the final result
func QuerySimple(ctx context.Context, prompt string, options *ClaudeCodeOptions) (*ResultMessage, []Message, error) {
	ch := Query(ctx, prompt, options)
	
	var messages []Message
	var result *ResultMessage
	
	for msgResult := range ch {
		if msgResult.Error != nil {
			return nil, messages, msgResult.Error
		}
		
		messages = append(messages, msgResult.Message)
		
		if res, ok := msgResult.Message.(ResultMessage); ok {
			result = &res
		}
	}
	
	if result == nil {
		return nil, messages, &TransportError{Message: "no result message received"}
	}
	
	return result, messages, nil
}

// Next is a helper method to get the next message from a channel
func (ch MessageChannel) Next() (Message, error) {
	result, ok := <-ch
	if !ok {
		return nil, io.EOF
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return result.Message, nil
}

// Collect collects all messages from the channel into a slice
func (ch MessageChannel) Collect(ctx context.Context) ([]Message, error) {
	var messages []Message
	
	for {
		select {
		case result, ok := <-ch:
			if !ok {
				return messages, nil
			}
			if result.Error != nil {
				return messages, result.Error
			}
			messages = append(messages, result.Message)
		case <-ctx.Done():
			return messages, ctx.Err()
		}
	}
}