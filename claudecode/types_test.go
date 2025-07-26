package claudecode

import (
	"encoding/json"
	"testing"
)

func TestMessageTypes(t *testing.T) {
	tests := []struct {
		name     string
		message  Message
		wantType MessageType
	}{
		{
			name:     "UserMessage",
			message:  UserMessage{Content: "Hello"},
			wantType: MessageTypeUser,
		},
		{
			name: "AssistantMessage",
			message: AssistantMessage{
				Message: struct {
					Content []json.RawMessage `json:"content"`
					ID      string            `json:"id"`
					Role    string            `json:"role"`
					Model   string            `json:"model"`
				}{
					Content: []json.RawMessage{json.RawMessage(`{"type":"text","text":"Hi"}`)},
				},
			},
			wantType: MessageTypeAssistant,
		},
		{
			name:     "SystemMessage",
			message:  SystemMessage{Subtype: "test", Data: json.RawMessage(`{}`)},
			wantType: MessageTypeSystem,
		},
		{
			name:     "ResultMessage",
			message:  ResultMessage{Content: "Done"},
			wantType: MessageTypeResult,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.message.Type(); got != tt.wantType {
				t.Errorf("Type() = %v, want %v", got, tt.wantType)
			}
		})
	}
}

func TestParseContentBlock(t *testing.T) {
	tests := []struct {
		name    string
		json    string
		want    ContentBlock
		wantErr bool
	}{
		{
			name: "TextBlock",
			json: `{"type": "text", "text": "Hello"}`,
			want: TextBlock{Type: "text", Text: "Hello"},
		},
		{
			name: "ToolUseBlock",
			json: `{"type": "tool_use", "id": "123", "name": "test", "input": {}}`,
			want: ToolUseBlock{Type: "tool_use", ID: "123", Name: "test", Input: json.RawMessage(`{}`)},
		},
		{
			name: "ToolResultBlock",
			json: `{"type": "tool_result", "tool_use_id": "123", "output": "result"}`,
			want: ToolResultBlock{Type: "tool_result", ToolUseID: "123", Output: stringPtr("result")},
		},
		{
			name:    "InvalidType",
			json:    `{"type": "invalid"}`,
			wantErr: true,
		},
		{
			name:    "MissingType",
			json:    `{}`,
			wantErr: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			block, err := ParseContentBlock(json.RawMessage(tt.json))
			
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseContentBlock() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if !tt.wantErr {
				// Compare type
				if block.BlockType() != tt.want.BlockType() {
					t.Errorf("BlockType() = %v, want %v", block.BlockType(), tt.want.BlockType())
				}
			}
		})
	}
}

func TestCostUsageMarshaling(t *testing.T) {
	cost := Cost{
		InputCached:       100,
		InputUncached:     200,
		Output:            50,
		InputCachedCost:   0.001,
		InputUncachedCost: 0.002,
		OutputCost:        0.005,
		TotalCost:         0.008,
	}
	
	data, err := json.Marshal(cost)
	if err != nil {
		t.Fatalf("Failed to marshal Cost: %v", err)
	}
	
	var decoded Cost
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal Cost: %v", err)
	}
	
	if decoded != cost {
		t.Errorf("Unmarshaled Cost = %+v, want %+v", decoded, cost)
	}
}

func stringPtr(s string) *string {
	return &s
}