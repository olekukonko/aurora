package aurora

import (
	"bytes"
	"regexp"
	"strings"
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
	n.JSON("TestJSON", testData)

	output := buf.String()
	// Remove ANSI color codes for testing
	cleanOutput := regexp.MustCompile(`\x1b\[[0-9;]*m`).ReplaceAllString(output, "")

	if !strings.Contains(cleanOutput, "TestJSON: JSON") {
		t.Errorf("expected output to contain JSON header, got %q", output)
	}
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
