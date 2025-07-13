package claudecode

import (
	"testing"
)

func TestDefaultOptions(t *testing.T) {
	opts := DefaultOptions()
	
	if opts == nil {
		t.Fatal("DefaultOptions() returned nil")
	}
	
	if opts.MaxThinkingTokens == nil {
		t.Error("MaxThinkingTokens should not be nil")
	} else if *opts.MaxThinkingTokens != 8000 {
		t.Errorf("MaxThinkingTokens = %d, want 8000", *opts.MaxThinkingTokens)
	}
}

func TestPermissionModeValues(t *testing.T) {
	tests := []struct {
		mode PermissionMode
		want string
	}{
		{PermissionModeAsk, "ask"},
		{PermissionModeAuto, "auto"},
	}
	
	for _, tt := range tests {
		if string(tt.mode) != tt.want {
			t.Errorf("PermissionMode %v = %s, want %s", tt.mode, string(tt.mode), tt.want)
		}
	}
}

func TestMCPServerTypeValues(t *testing.T) {
	tests := []struct {
		serverType MCPServerType
		want       string
	}{
		{MCPServerTypeStdio, "stdio"},
		{MCPServerTypeSSE, "sse"},
		{MCPServerTypeHTTP, "http"},
	}
	
	for _, tt := range tests {
		if string(tt.serverType) != tt.want {
			t.Errorf("MCPServerType %v = %s, want %s", tt.serverType, string(tt.serverType), tt.want)
		}
	}
}

func TestClaudeCodeOptionsFields(t *testing.T) {
	// Test that all fields can be set
	maxThinking := 10000
	systemPrompt := "test prompt"
	appendPrompt := "append this"
	permMode := PermissionModeAuto
	resume := "session-123"
	maxTurns := 5
	model := "claude-3"
	permToolName := "custom-permission"
	cwd := "/home/user"
	
	opts := ClaudeCodeOptions{
		AllowedTools:             []string{"tool1", "tool2"},
		MaxThinkingTokens:        &maxThinking,
		SystemPrompt:             &systemPrompt,
		AppendSystemPrompt:       &appendPrompt,
		PermissionMode:           &permMode,
		ContinueConversation:     true,
		Resume:                   &resume,
		MaxTurns:                 &maxTurns,
		DisallowedTools:          []string{"dangerous-tool"},
		Model:                    &model,
		PermissionPromptToolName: &permToolName,
		CWD:                      &cwd,
		MCPServers: []MCPServerConfig{
			{
				Type: MCPServerTypeStdio,
				StdioConfig: &MCPStdioConfig{
					Command: "mcp-server",
					Args:    []string{"--port", "8080"},
					Env:     map[string]string{"DEBUG": "true"},
				},
			},
		},
	}
	
	// Basic validation
	if len(opts.AllowedTools) != 2 {
		t.Errorf("AllowedTools length = %d, want 2", len(opts.AllowedTools))
	}
	
	if *opts.MaxThinkingTokens != maxThinking {
		t.Errorf("MaxThinkingTokens = %d, want %d", *opts.MaxThinkingTokens, maxThinking)
	}
	
	if !opts.ContinueConversation {
		t.Error("ContinueConversation should be true")
	}
	
	if len(opts.MCPServers) != 1 {
		t.Errorf("MCPServers length = %d, want 1", len(opts.MCPServers))
	}
}