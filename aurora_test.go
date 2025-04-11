package aurora

import (
	"bytes"
	"github.com/fatih/color"
	"regexp"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestNotifier_Logf(t *testing.T) {
	// Use a buffer to capture the log output.
	var buf bytes.Buffer
	n := New(&buf)

	// Log an info-level message.
	n.Logf(InfoLevel, "test log %d", 123)

	output := buf.String()
	if !strings.Contains(output, "test log 123") {
		t.Errorf("expected output to contain %q, got %q", "test log 123", output)
	}

	// Check that a timestamp is present (roughly).
	if !strings.Contains(output, time.Now().Format("2006-01-02")) {
		t.Errorf("expected output to contain a timestamp, got %q", output)
	}
}

func TestNotifier_With(t *testing.T) {
	var buf bytes.Buffer
	n := New(&buf)

	// Create a contextual logger with a prefix.
	sub := n.With("module")
	sub.Info("message")

	output := buf.String()
	if !strings.Contains(output, "[module]") {
		t.Errorf("expected output to contain the prefix [module], got %q", output)
	}
	if !strings.Contains(output, "message") {
		t.Errorf("expected output to contain %q, got %q", "message", output)
	}
}

func TestNotifier_JSON(t *testing.T) {
	var buf bytes.Buffer
	n := New(&buf)

	testData := map[string]interface{}{"key": "value"}
	n.JSON(testData) // Just test raw JSON output

	output := buf.String()
	// Remove ANSI color codes for testing
	cleanOutput := regexp.MustCompile(`\x1b\[[0-9;]*m`).ReplaceAllString(output, "")

	// Expect raw JSON without any title
	if !strings.Contains(cleanOutput, `"key":"value"`) {
		t.Errorf("expected output to contain JSON key-value, got %q", output)
	}
}

func TestNotifier_Printf(t *testing.T) {
	var buf bytes.Buffer
	n := New(&buf)

	// Test Printf which writes a plain message without a level symbol or timestamp.
	n.Printf(InfoLevel, "plain message %s", "test")
	output := buf.String()
	if !strings.Contains(output, "plain message test") {
		t.Errorf("expected output to contain %q, got %q", "plain message test", output)
	}
}

// TestInlinef tests basic Inlinef functionality with different levels
func TestInlinef(t *testing.T) {
	color.NoColor = true // Disable colors for predictable output
	defer func() { color.NoColor = false }()

	tests := []struct {
		name       string
		level      LogLevel
		message    string
		wantOutput string
	}{
		{
			name:       "Info level",
			level:      InfoLevel,
			message:    "Info test",
			wantOutput: "[✔] Info test\n",
		},
		{
			name:       "Error level",
			level:      ErrorLevel,
			message:    "Error test",
			wantOutput: "[✘] Error test\n",
		},
		{
			name:       "Debug level",
			level:      DebugLevel,
			message:    "Debug test",
			wantOutput: "[⧳] Debug test\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			n := New(&buf)

			n.Inlinef(tt.level, tt.message)

			got := buf.String()
			if got != tt.wantOutput {
				t.Errorf("Inlinef() = %q, want %q", got, tt.wantOutput)
			}
		})
	}
}

// TestSuccess tests the Success method
func TestSuccess(t *testing.T) {
	color.NoColor = true
	defer func() { color.NoColor = false }()

	var buf bytes.Buffer
	n := New(&buf)

	n.Success("Operation completed")

	output := buf.String()
	if !strings.Contains(output, "✓ Operation completed") {
		t.Errorf("Success() expected '✓ Operation completed', got: %q", output)
	}
}

// TestFailure tests the Failure method
func TestFailure(t *testing.T) {
	color.NoColor = true
	defer func() { color.NoColor = false }()

	var buf bytes.Buffer
	n := New(&buf)

	n.Failure("Operation failed")

	output := buf.String()
	if !strings.Contains(output, "✗ Operation failed") {
		t.Errorf("Failure() expected '✗ Operation failed', got: %q", output)
	}
}

// TestWithPrefix tests the With method for prefixing
func TestWithPrefix(t *testing.T) {
	color.NoColor = true
	defer func() { color.NoColor = false }()

	var buf bytes.Buffer
	n := New(&buf)
	nWithPrefix := n.With("TEST")

	nWithPrefix.Inlinef(InfoLevel, "Prefixed message")

	output := buf.String()
	if !strings.Contains(output, "[TEST] Prefixed message") {
		t.Errorf("With() expected '[TEST] Prefixed message', got: %q", output)
	}
	if !strings.Contains(output, "[✔]") {
		t.Errorf("With() expected Info symbol '[✔]', got: %q", output)
	}
}

// TestConcurrentInlinef tests thread-safety of Inlinef
func TestConcurrentInlinef(t *testing.T) {
	color.NoColor = true
	defer func() { color.NoColor = false }()

	var buf bytes.Buffer
	n := New(&buf)

	var wg sync.WaitGroup
	levels := []LogLevel{DebugLevel, InfoLevel, ErrorLevel}

	for i, level := range levels {
		wg.Add(1)
		go func(lvl LogLevel, idx int) {
			defer wg.Done()
			n.Inlinef(lvl, "Message %d", idx)
		}(level, i)
	}

	wg.Wait()

	output := buf.String()
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) != len(levels) {
		t.Errorf("Expected %d lines, got %d: %q", len(levels), len(lines), output)
	}
}

// TestRobot tests the Robot method for ASCII art
func TestRobot(t *testing.T) {
	color.NoColor = true
	defer func() { color.NoColor = false }()

	var buf bytes.Buffer
	n := New(&buf)

	n.Robot(InfoLevel)

	output := buf.String()
	if output == "" {
		t.Errorf("Robot() expected non-empty ASCII art, got empty string")
	}
	if !strings.Contains(output, "\n") {
		t.Errorf("Robot() expected multi-line output, got: %q", output)
	}
}

// TestJSON tests the JSON method
func TestJSON(t *testing.T) {
	color.NoColor = true
	defer func() { color.NoColor = false }()

	var buf bytes.Buffer
	n := New(&buf)

	data := map[string]string{"key": "value"}
	n.JSONTitle("Test Data", data)

	output := buf.String()
	expectedPrefix := "[⧳] Test Data: JSON ↴↴\n"
	// Match the actual JSON output format (no space after colon)
	expectedJSON := "{\n  \"key\":\"value\"\n}\n\n"

	if !strings.HasPrefix(output, expectedPrefix) {
		t.Errorf("JSON() expected prefix %q, got: %q", expectedPrefix, output)
	}
	if !strings.Contains(output, expectedJSON) {
		t.Errorf("JSON() expected marshaled data %q, got: %q", expectedJSON, output)
	}
}

// TestDefaultNotifier tests the global Default notifier
func TestDefaultNotifier(t *testing.T) {
	color.NoColor = true
	defer func() { color.NoColor = false }()

	var buf bytes.Buffer
	Default = New(&buf) // Override Default for testing

	Info("Default test")

	output := buf.String()
	if !strings.Contains(output, "[✔] Default test") {
		t.Errorf("Default Info() expected '[✔] Default test', got: %q", output)
	}
}
