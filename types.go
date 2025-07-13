package claudecode

import (
	"encoding/json"
	"fmt"
)

// MessageType represents the type of message
type MessageType string

const (
	MessageTypeUser      MessageType = "user"
	MessageTypeAssistant MessageType = "assistant"
	MessageTypeSystem    MessageType = "system"
	MessageTypeResult    MessageType = "result"
)

// Message is an interface implemented by all message types
type Message interface {
	Type() MessageType
}

// UserMessage represents a message from the user
type UserMessage struct {
	Content string `json:"content"`
}

func (m UserMessage) Type() MessageType {
	return MessageTypeUser
}

// AssistantMessage represents a message from the assistant
type AssistantMessage struct {
	Content []json.RawMessage `json:"content"`
}

func (m AssistantMessage) Type() MessageType {
	return MessageTypeAssistant
}

// SystemMessage represents a system message
type SystemMessage struct {
	Subtype string          `json:"subtype"`
	Data    json.RawMessage `json:"data"`
}

func (m SystemMessage) Type() MessageType {
	return MessageTypeSystem
}

// ResultMessage represents the final result
type ResultMessage struct {
	Content  string          `json:"content"`
	Cost     *Cost           `json:"cost,omitempty"`
	Usage    *Usage          `json:"usage,omitempty"`
	Session  *SessionInfo    `json:"session,omitempty"`
	Metadata json.RawMessage `json:"metadata,omitempty"`
}

func (m ResultMessage) Type() MessageType {
	return MessageTypeResult
}

// ContentBlock is an interface for different types of content blocks
type ContentBlock interface {
	BlockType() string
}

// TextBlock represents text content
type TextBlock struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

func (b TextBlock) BlockType() string {
	return "text"
}

// ToolUseBlock represents a tool invocation
type ToolUseBlock struct {
	Type  string          `json:"type"`
	ID    string          `json:"id"`
	Name  string          `json:"name"`
	Input json.RawMessage `json:"input"`
}

func (b ToolUseBlock) BlockType() string {
	return "tool_use"
}

// ToolResultBlock represents the result of a tool execution
type ToolResultBlock struct {
	Type       string            `json:"type"`
	ToolUseID  string            `json:"tool_use_id"`
	Content    []json.RawMessage `json:"content,omitempty"`
	IsError    bool              `json:"is_error,omitempty"`
	Output     *string           `json:"output,omitempty"`
	StatusCode *int              `json:"status_code,omitempty"`
	Warnings   []string          `json:"warnings,omitempty"`
	Metadata   json.RawMessage   `json:"metadata,omitempty"`
}

func (b ToolResultBlock) BlockType() string {
	return "tool_result"
}

// Cost represents the cost information
type Cost struct {
	InputCached          int     `json:"input_cached"`
	InputUncached        int     `json:"input_uncached"`
	Output               int     `json:"output"`
	InputCachedCost      float64 `json:"input_cached_cost"`
	InputUncachedCost    float64 `json:"input_uncached_cost"`
	OutputCost           float64 `json:"output_cost"`
	TotalCost            float64 `json:"total_cost"`
	CustomizationWeights int     `json:"customization_weights,omitempty"`
	CustomizationCost    float64 `json:"customization_cost,omitempty"`
}

// Usage represents token usage information
type Usage struct {
	InputTokens         int `json:"input_tokens"`
	OutputTokens        int `json:"output_tokens"`
	CacheCreationTokens int `json:"cache_creation_tokens,omitempty"`
	CacheReadTokens     int `json:"cache_read_tokens,omitempty"`
	ThinkingInputTokens int `json:"thinking_input_tokens,omitempty"`
	TotalTokens         int `json:"total_tokens"`
}

// SessionInfo represents session information
type SessionInfo struct {
	ID            string          `json:"id"`
	Memory        json.RawMessage `json:"memory,omitempty"`
	ContextWindow int             `json:"context_window,omitempty"`
}

// ParseContentBlock parses a JSON RawMessage into a ContentBlock
func ParseContentBlock(data json.RawMessage) (ContentBlock, error) {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}

	typeData, ok := raw["type"]
	if !ok {
		return nil, fmt.Errorf("missing 'type' field in content block")
	}

	var blockType string
	if err := json.Unmarshal(typeData, &blockType); err != nil {
		return nil, err
	}

	switch blockType {
	case "text":
		var block TextBlock
		if err := json.Unmarshal(data, &block); err != nil {
			return nil, err
		}
		return block, nil
	case "tool_use":
		var block ToolUseBlock
		if err := json.Unmarshal(data, &block); err != nil {
			return nil, err
		}
		return block, nil
	case "tool_result":
		var block ToolResultBlock
		if err := json.Unmarshal(data, &block); err != nil {
			return nil, err
		}
		return block, nil
	default:
		return nil, fmt.Errorf("unknown content block type: %s", blockType)
	}
}