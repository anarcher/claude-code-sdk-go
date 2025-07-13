package claudecode

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

const (
	maxBufferSize = 1024 * 1024 // 1MB
	readTimeout   = 5 * time.Minute
)

// Transport handles communication with the Claude CLI subprocess
type Transport struct {
	cmd      *exec.Cmd
	stdin    io.WriteCloser
	stdout   io.ReadCloser
	stderr   io.ReadCloser
	scanner  *bufio.Scanner
	errChan  chan error
	closed   bool
	mu       sync.Mutex
	ctx      context.Context
	cancel   context.CancelFunc
}

// NewTransport creates a new transport with the given CLI path and arguments
func NewTransport(ctx context.Context, cliPath string, args []string) (*Transport, error) {
	ctx, cancel := context.WithCancel(ctx)
	
	cmd := exec.CommandContext(ctx, cliPath, args...)
	
	// Set up pipes
	stdin, err := cmd.StdinPipe()
	if err != nil {
		cancel()
		return nil, &TransportError{Message: "failed to create stdin pipe", Cause: err}
	}
	
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		cancel()
		return nil, &TransportError{Message: "failed to create stdout pipe", Cause: err}
	}
	
	stderr, err := cmd.StderrPipe()
	if err != nil {
		cancel()
		return nil, &TransportError{Message: "failed to create stderr pipe", Cause: err}
	}
	
	// Start the process
	if err := cmd.Start(); err != nil {
		cancel()
		return nil, &TransportError{Message: "failed to start CLI process", Cause: err}
	}
	
	scanner := bufio.NewScanner(stdout)
	scanner.Buffer(make([]byte, 0, maxBufferSize), maxBufferSize)
	
	t := &Transport{
		cmd:     cmd,
		stdin:   stdin,
		stdout:  stdout,
		stderr:  stderr,
		scanner: scanner,
		errChan: make(chan error, 1),
		ctx:     ctx,
		cancel:  cancel,
	}
	
	// Start error monitoring
	go t.monitorErrors()
	
	return t, nil
}

// monitorErrors reads from stderr in the background
func (t *Transport) monitorErrors() {
	scanner := bufio.NewScanner(t.stderr)
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			select {
			case t.errChan <- &CLIError{Message: line, Code: 1}:
			default:
				// Channel is full, drop the error
			}
		}
	}
	if err := scanner.Err(); err != nil && err != io.EOF {
		select {
		case t.errChan <- &TransportError{Message: "error reading stderr", Cause: err}:
		default:
		}
	}
}

// Send sends a message to the CLI
func (t *Transport) Send(prompt string) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	
	if t.closed {
		return ErrTransportClosed
	}
	
	// The CLI expects JSON input
	input := map[string]string{"prompt": prompt}
	data, err := json.Marshal(input)
	if err != nil {
		return &TransportError{Message: "failed to marshal input", Cause: err}
	}
	
	_, err = fmt.Fprintf(t.stdin, "%s\n", data)
	if err != nil {
		return &TransportError{Message: "failed to write to stdin", Cause: err}
	}
	
	return nil
}

// Receive reads the next message from the CLI
func (t *Transport) Receive() (json.RawMessage, error) {
	t.mu.Lock()
	if t.closed {
		t.mu.Unlock()
		return nil, ErrTransportClosed
	}
	t.mu.Unlock()
	
	// Check for errors first
	select {
	case err := <-t.errChan:
		return nil, err
	default:
	}
	
	// Set up timeout
	done := make(chan bool, 1)
	var line string
	var scanErr error
	
	go func() {
		if t.scanner.Scan() {
			line = t.scanner.Text()
		} else {
			scanErr = t.scanner.Err()
			if scanErr == nil {
				scanErr = io.EOF
			}
		}
		done <- true
	}()
	
	select {
	case <-done:
		if scanErr != nil {
			return nil, &TransportError{Message: "failed to read from stdout", Cause: scanErr}
		}
		
		line = strings.TrimSpace(line)
		if line == "" {
			// Empty line, try again
			return t.Receive()
		}
		
		// Validate JSON
		var msg json.RawMessage
		if err := json.Unmarshal([]byte(line), &msg); err != nil {
			return nil, &ParseError{Message: "invalid JSON", Data: line}
		}
		
		return msg, nil
		
	case <-time.After(readTimeout):
		return nil, ErrTimeout
		
	case <-t.ctx.Done():
		return nil, t.ctx.Err()
	}
}

// Close closes the transport and terminates the subprocess
func (t *Transport) Close() error {
	t.mu.Lock()
	defer t.mu.Unlock()
	
	if t.closed {
		return nil
	}
	
	t.closed = true
	t.cancel()
	
	// Close pipes
	if t.stdin != nil {
		t.stdin.Close()
	}
	if t.stdout != nil {
		t.stdout.Close()
	}
	if t.stderr != nil {
		t.stderr.Close()
	}
	
	// Wait for process to exit (with timeout)
	done := make(chan error, 1)
	go func() {
		done <- t.cmd.Wait()
	}()
	
	select {
	case err := <-done:
		// Process exited
		if err != nil && !strings.Contains(err.Error(), "killed") {
			return &TransportError{Message: "process exited with error", Cause: err}
		}
		return nil
	case <-time.After(5 * time.Second):
		// Force kill if it doesn't exit gracefully
		if t.cmd.Process != nil {
			t.cmd.Process.Kill()
		}
		return nil
	}
}

// findCLI attempts to find the Claude CLI executable
func findCLI() (string, error) {
	// Check common locations
	paths := []string{
		"claude",
		"claude-cli",
		"/usr/local/bin/claude",
		"/usr/bin/claude",
		"/opt/homebrew/bin/claude",
	}
	
	// Check PATH
	if path, err := exec.LookPath("claude"); err == nil {
		return path, nil
	}
	
	// Check specific paths
	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}
	
	// Check CLAUDE_CLI_PATH environment variable
	if envPath := os.Getenv("CLAUDE_CLI_PATH"); envPath != "" {
		if _, err := os.Stat(envPath); err == nil {
			return envPath, nil
		}
	}
	
	return "", ErrCLINotFound
}