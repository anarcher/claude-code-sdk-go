package claudecode

import (
	"errors"
	"strings"
	"testing"
)

func TestErrors(t *testing.T) {
	tests := []struct {
		name    string
		err     error
		wantMsg string
	}{
		{
			name:    "CLIError",
			err:     &CLIError{Message: "command failed", Code: 1},
			wantMsg: "CLI error (code 1): command failed",
		},
		{
			name:    "ParseError",
			err:     &ParseError{Message: "invalid JSON", Data: `{"bad": json}`},
			wantMsg: `parse error: invalid JSON (data: {"bad": json})`,
		},
		{
			name:    "TransportError without cause",
			err:     &TransportError{Message: "connection lost"},
			wantMsg: "transport error: connection lost",
		},
		{
			name:    "TransportError with cause",
			err:     &TransportError{Message: "failed to connect", Cause: errors.New("timeout")},
			wantMsg: "transport error: failed to connect: timeout",
		},
		{
			name:    "ValidationError",
			err:     &ValidationError{Field: "prompt", Message: "cannot be empty"},
			wantMsg: "validation error in prompt: cannot be empty",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.wantMsg {
				t.Errorf("Error() = %v, want %v", got, tt.wantMsg)
			}
		})
	}
}

func TestTransportErrorUnwrap(t *testing.T) {
	cause := errors.New("original error")
	err := &TransportError{Message: "wrapped", Cause: cause}
	
	if unwrapped := errors.Unwrap(err); unwrapped != cause {
		t.Errorf("Unwrap() = %v, want %v", unwrapped, cause)
	}
}

func TestSentinelErrors(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want string
	}{
		{"ErrCLINotFound", ErrCLINotFound, "claude CLI not found"},
		{"ErrInvalidMessage", ErrInvalidMessage, "invalid message format"},
		{"ErrTransportClosed", ErrTransportClosed, "transport is closed"},
		{"ErrBufferOverflow", ErrBufferOverflow, "buffer overflow"},
		{"ErrTimeout", ErrTimeout, "operation timed out"},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !strings.Contains(tt.err.Error(), tt.want) {
				t.Errorf("Error message %q doesn't contain %q", tt.err.Error(), tt.want)
			}
		})
	}
}