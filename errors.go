package claudecode

import (
	"errors"
	"fmt"
)

var (
	// ErrCLINotFound is returned when the Claude CLI cannot be found
	ErrCLINotFound = errors.New("claude CLI not found")

	// ErrInvalidMessage is returned when a message cannot be parsed
	ErrInvalidMessage = errors.New("invalid message format")

	// ErrTransportClosed is returned when trying to use a closed transport
	ErrTransportClosed = errors.New("transport is closed")

	// ErrBufferOverflow is returned when the buffer limit is exceeded
	ErrBufferOverflow = errors.New("buffer overflow")

	// ErrTimeout is returned when an operation times out
	ErrTimeout = errors.New("operation timed out")
)

// CLIError represents an error from the Claude CLI
type CLIError struct {
	Message string
	Code    int
}

func (e *CLIError) Error() string {
	return fmt.Sprintf("CLI error (code %d): %s", e.Code, e.Message)
}

// ParseError represents an error parsing a message
type ParseError struct {
	Message string
	Data    string
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("parse error: %s (data: %s)", e.Message, e.Data)
}

// TransportError represents an error in the transport layer
type TransportError struct {
	Message string
	Cause   error
}

func (e *TransportError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("transport error: %s: %v", e.Message, e.Cause)
	}
	return fmt.Sprintf("transport error: %s", e.Message)
}

func (e *TransportError) Unwrap() error {
	return e.Cause
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error in %s: %s", e.Field, e.Message)
}