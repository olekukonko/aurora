package aurora

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/mattes/go-asciibot"
	"github.com/nwidger/jsoncolor"
	"io"
	"os"
	"strings"
	"sync"
	"time"
)

// LogLevel defines the severity of the log message.
// It uses an integer type to allow comparison between levels.
// Use the constants below to specify different log levels.
type LogLevel int

const (
	IconSuccess = "✓" // Success icon used in Success/Failure methods
	IconError   = "✗" // Error icon used in Success/Failure methods
)

// Indentation constants for consistent JSON formatting across the application.
// These provide standardized ways to format JSON output while maintaining readability.
const (
	// IndentNone provides compact JSON output without any indentation or line breaks.
	// Ideal for network transmission or minimal storage requirements.
	IndentNone = ""

	// IndentSpace uses a single space for indentation.
	// Provides minimal readability improvement while keeping output compact.
	IndentSpace = " "

	// IndentSpace2 uses two spaces per indentation level (most common standard).
	// Recommended for general use as it balances readability and compactness.
	// Complies with many style guides including Google JSON Style Guide.
	IndentSpace2 = "  "

	// IndentSpace4 uses four spaces per indentation level.
	// Provides clearer nesting visibility at the cost of more horizontal space.
	// Used in some corporate coding standards.
	IndentSpace4 = "    "

	// IndentTab uses tab characters for indentation.
	// Preferred by developers who use tab-based indentation in their editors.
	// Note: Tab width varies across editors/viewers.
	IndentTab = "\t"

	// IndentDebug uses a visually distinct pattern (bullet + space).
	// Particularly useful for debugging as it makes indentation levels very visible.
	// Not recommended for production use.
	IndentDebug = "• "
)

// Log level constants in order of increasing severity
// These define the available logging levels from least to most severe
const (
	DebugLevel LogLevel = iota
	InfoLevel
	NoticeLevel
	WarnLevel
	ErrorLevel
	AlertLevel
	CriticalLevel
	NoLevel
)

// Default symbols for each log level
// These provide visual indicators for different log severities
var defaultSymbols = map[LogLevel]string{
	AlertLevel:    "[✭]", // Alert symbol for attention-grabbing messages
	InfoLevel:     "[✔]", // Info symbol for general information
	ErrorLevel:    "[✘]", // Error symbol for error conditions
	NoticeLevel:   "[⚑]", // Notice symbol for notable events
	DebugLevel:    "[⧳]", // Debug symbol for debugging output
	WarnLevel:     "[⚠]", // Warning symbol for potential issues
	CriticalLevel: "[‼]", // Critical symbol for severe problems
	NoLevel:       " ",   // No symbol for plain messages
}

// Default colors for each log level
// These assign distinct colors to make log levels easily distinguishable
var defaultColors = map[LogLevel]*color.Color{
	AlertLevel:    color.New(color.FgHiBlue),    // Blue for alerts stands out
	InfoLevel:     color.New(color.FgHiGreen),   // Green for info indicates normalcy
	ErrorLevel:    color.New(color.FgHiRed),     // Red for errors signals problems
	NoticeLevel:   color.New(color.FgHiYellow),  // Yellow for notices draws attention
	DebugLevel:    color.New(color.FgHiCyan),    // Cyan for debug aids developers
	WarnLevel:     color.New(color.FgHiMagenta), // Magenta for warnings is distinct
	CriticalLevel: color.New(color.FgHiWhite),   // White for critical is highly visible
	NoLevel:       color.New(color.FgHiBlack),   // Gray for no level is unobtrusive
}

// Package-level customization
// These variables allow global configuration of logging appearance
var (
	symbols = make(map[LogLevel]string)
	colors  = make(map[LogLevel]*color.Color)
	mu      sync.RWMutex
)

// Formater defines a custom formatting function signature
// It allows users to specify their own formatting logic
// that can be used with the Format method
type Formater func(format string, a ...interface{}) string

// Default is a global Notifier instance that writes to os.Stdout
// It provides a convenient way to log without creating a new instance
var Default = New(os.Stdout)

func init() {
	ResetSymbols() // Initialize symbols to default values
	ResetColors()  // Initialize colors to default values
}

// Notifier provides structured, colorful logging capabilities
// It handles synchronization and output formatting
type Notifier struct {
	mu     *sync.Mutex // Protects concurrent access
	output io.Writer   // Destination for log messages
	prefix string      // Optional prefix for all messages
}

// New creates Notifier that writes to given io.Writer
// Uses os.Stdout if writer is nil for convenience
// Returns a pointer to the new Notifier instance
func New(w io.Writer) *Notifier {
	if w == nil {
		w = os.Stdout
	}
	return &Notifier{
		mu:     &sync.Mutex{},
		output: w,
		prefix: "",
	}
}

// Alert logs a message at Alert level
// Useful for important but non-critical notifications
func (n *Notifier) Alert(f string, a ...any) { n.Inlinef(AlertLevel, f, a...) }

// Br inserts a single blank line in the output
// Helps with visual separation of log entries
func (n *Notifier) Br() { n.Line(1) }

// Color writes a message with specific color, ignoring log level colors
// Useful for special messages that need distinct coloring
// Bypasses the default level-based coloring system
func (n *Notifier) Color(c *color.Color, format string, args ...any) {
	n.mu.Lock()
	defer n.mu.Unlock()
	c.Fprint(n.output, fmt.Sprintf(format, args...))
}

// Critical logs a message at Critical level
// Used for severe issues requiring immediate attention
func (n *Notifier) Critical(f string, a ...any) { n.Inlinef(CriticalLevel, f, a...) }

// Debug logs a message at Debug level
// Intended for developer-facing diagnostic information
func (n *Notifier) Debug(f string, a ...any) { n.Inlinef(DebugLevel, f, a...) }

// Error logs a message at Error level
// Indicates problems that need attention
func (n *Notifier) Error(f string, a ...any) { n.Inlinef(ErrorLevel, f, a...) }

// Failure prints error message with red color and crossmark prefix
// Provides consistent error message formatting across application
// Uses the ErrorLevel for consistency
func (n *Notifier) Failure(format string, args ...any) {
	n.Inlinef(ErrorLevel, n.f(IconError, " ", format), args...)
}

// Format writes message using custom formatter function
// Allows complete control over message formatting while maintaining
// thread safety and output consistency through mutex locking
func (n *Notifier) Format(formatter Formater, format string, args ...any) {
	n.mu.Lock()
	defer n.mu.Unlock()
	colors[NoLevel].Fprint(n.output, formatter(format, args...))
}

// Func executes function and writes output with specified log level color
// The function is only called when actually writing to output
// Useful for expensive computations that should only run when logged
func (n *Notifier) Func(level LogLevel, fn func() string) {
	n.mu.Lock()
	defer n.mu.Unlock()
	colors[level].Fprint(n.output, fn())
}

// Highlight writes text with yellow background highlight
// Excellent for drawing attention to important log messages
// Uses a distinct background color for emphasis
func (n *Notifier) Highlight(format string, args ...any) {
	n.Color(color.New(color.BgYellow, color.FgBlack), format, args...)
}

// If conditionally logs message based on boolean condition
// Simplifies conditional logging without cluttering code
// Reduces need for external if statements
func (n *Notifier) If(condition bool, level LogLevel, format string, args ...any) {
	if condition {
		n.Inlinef(level, format, args...)
	}
}

// Info logs a message at Info level
// Used for general operational information
func (n *Notifier) Info(f string, a ...any) { n.Inlinef(InfoLevel, f, a...) }

// JSON logs JSON data without title (no indentation)
func (n *Notifier) JSON(values ...any) {
	n.JSONIndent("", IndentNone, values...)
}

// JSONTitle logs JSON data with title (no indentation)
func (n *Notifier) JSONTitle(title string, values ...any) {
	n.JSONIndent(title, IndentSpace2, values...)
}

// JSONIndent logs JSON data with custom indentation
func (n *Notifier) JSONIndent(title string, indent string, values ...any) {
	if title != "" {
		n.Inlinef(DebugLevel, "%s: JSON ↴↴", title)
	}
	n.mu.Lock()
	defer n.mu.Unlock()

	formatter := jsoncolor.NewFormatter()
	formatter.Indent = indent
	for _, v := range values {
		data, err := jsoncolor.MarshalIndent(v, "", indent)
		if err != nil {
			n.Logf(ErrorLevel, "failed to marshal JSON: %v", err)
			continue
		}
		n.output.Write(data)
		n.output.Write([]byte{'\n'})
	}
	n.output.Write([]byte{'\n'})
}

// Inlinef writes single-line log without timestamp
// Ideal for compact output where timestamps aren't needed
// Includes level symbol and color
func (n *Notifier) Inlinef(level LogLevel, format string, args ...any) {
	n.mu.Lock()
	defer n.mu.Unlock()

	symbol := symbols[level]
	msg := fmt.Sprintf(format, args...)
	msg = n.formatWithPrefix(msg)
	line := fmt.Sprintf("%s %s\n", symbol, msg)

	colors[level].Fprint(n.output, line)
}

// Line inserts specified number of blank lines
// Useful for visually separating log sections
// Helps organize output readability
func (n *Notifier) Line(count int) {
	n.mu.Lock()
	defer n.mu.Unlock()
	colors[NoLevel].Fprint(n.output, fmt.Sprintf("%s", strings.Repeat("\n", count)))
}

// Logf writes formatted log with timestamp and level symbol
// Provides complete log message with all standard fields
// Includes timestamp for temporal context
func (n *Notifier) Logf(level LogLevel, format string, args ...any) {
	n.mu.Lock()
	defer n.mu.Unlock()

	timestamp := time.Now().Format("2006-01-02 03:04:05 PM")
	symbol := symbols[level]
	msg := fmt.Sprintf(format, args...)
	msg = n.formatWithPrefix(msg)
	line := fmt.Sprintf("%s %s %s\n", symbol, timestamp, msg)

	colors[level].Fprint(n.output, line)
}

// Notice logs a message at Notice level
// For events that should be noted but aren't problems
func (n *Notifier) Notice(f string, a ...any) { n.Inlinef(NoticeLevel, f, a...) }

// Panic logs a message at Critical level and then panics with the same message
// Used for unrecoverable errors that should halt program execution
func (n *Notifier) Panic(f string, a ...any) {
	msg := fmt.Sprintf(f, a...)
	n.Inlinef(CriticalLevel, msg)
	panic(msg)
}

// Printf writes plain message without timestamp or symbol
// Maintains prefix and color while being more minimal
// Useful for simple formatted output
func (n *Notifier) Printf(level LogLevel, format string, args ...any) {
	n.mu.Lock()
	defer n.mu.Unlock()

	msg := fmt.Sprintf(format, args...)
	msg = n.formatWithPrefix(msg)
	line := fmt.Sprintf("%s\n", msg)

	colors[level].Fprint(n.output, line)
}

// Robot displays random ASCII robot art
// Adds fun visual element to console output
// Makes logs more engaging
func (n *Notifier) Robot(level LogLevel) {
	n.mu.Lock()
	defer n.mu.Unlock()
	colors[level].Fprint(n.output, fmt.Sprintf("%s\n", asciibot.Random()))
}

// Success prints success message with green color and checkmark
// Standardized way to indicate successful operations
// Uses InfoLevel for positive feedback
func (n *Notifier) Success(format string, args ...any) {
	n.Inlinef(InfoLevel, n.f(IconSuccess, " ", format), args...)
}

// Warn logs a message at Warn level
// Indicates potential issues that aren't errors
func (n *Notifier) Warn(f string, a ...any) { n.Inlinef(WarnLevel, f, a...) }

// With creates new Notifier with additional prefix
// Enables contextual logging with shared configuration
// Maintains original Notifier's output and synchronization
func (n *Notifier) With(prefix string) *Notifier {
	newPrefix := prefix
	if n.prefix != "" {
		newPrefix = fmt.Sprintf("%s %s", n.prefix, prefix)
	}
	return &Notifier{
		mu:     n.mu,
		output: n.output,
		prefix: newPrefix,
	}
}

// formatWithPrefix adds the configured prefix to messages
// Internal helper method for consistent prefix handling
func (n *Notifier) formatWithPrefix(msg string) string {
	if n.prefix != "" {
		return fmt.Sprintf("[%s] %s", n.prefix, msg)
	}
	return msg
}

// f concatenates multiple arguments into a single string
// Internal helper for building formatted messages
func (n *Notifier) f(args ...any) string {
	s := strings.Builder{}
	for _, arg := range args {
		s.WriteString(fmt.Sprint(arg))
	}
	return s.String()
}

// Alert logs a message at Alert level using the default Notifier
// Convenience function for quick alerting
func Alert(f string, a ...any) { Default.Alert(f, a...) }

// Br inserts a single blank line using the default Notifier
// Shortcut for adding visual separation
func Br() { Line(1) }

// Color writes a message with specific color using default Notifier
// Quick way to add custom-colored output
func Color(c *color.Color, format string, args ...any) {
	Default.Color(c, format, args...)
}

// Critical logs a message at Critical level using default Notifier
// Convenient access to critical logging
func Critical(f string, a ...any) { Default.Critical(f, a...) }

// Debug logs a message at Debug level using default Notifier
// Quick debugging output
func Debug(f string, a ...any) { Default.Debug(f, a...) }

// Error logs a message at Error level using default Notifier
// Simple error reporting
func Error(f string, a ...any) { Default.Error(f, a...) }

// Failure logs an error message with error icon using default Notifier
// Standardized error formatting shortcut
func Failure(format string, args ...any) {
	Default.Failure(format, args...)
}

// Format uses custom formatter with default Notifier
// Flexible formatting with default instance
func Format(formatter Formater, format string, args ...any) {
	Default.Format(formatter, format, args...)
}

// Func executes function and logs result with default Notifier
// Lazy evaluation with default instance
func Func(level LogLevel, fn func() string) {
	Default.Func(level, fn)
}

// Highlight writes highlighted text using default Notifier
// Quick attention-grabbing output
func Highlight(format string, args ...any) {
	Default.Highlight(format, args...)
}

// If conditionally logs message using default Notifier
// Simplified conditional logging
func If(condition bool, level LogLevel, format string, args ...any) {
	Default.If(condition, level, format, args...)
}

// Info logs a message at Info level using default Notifier
// Standard informational logging
func Info(f string, a ...any) { Default.Info(f, a...) }

// Inlinef writes single-line log without timestamp using default Notifier
// Compact logging shortcut
func Inlinef(level LogLevel, f string, a ...any) { Default.Inlinef(level, f, a...) }

// JSON logs JSON data without title using default Notifier (no indentation)
// Structured data logging shortcut for compact output
func JSON(v ...any) { Default.JSON(v...) }

// JSONTitle logs JSON data with title using default Notifier (no indentation)
// Structured data logging with title for context
func JSONTitle(title string, v ...any) { Default.JSONTitle(title, v...) }

// JSONIndent logs JSON data with custom indentation using default Notifier
// Full control over indentation style
func JSONIndent(title string, indent string, values ...any) {
	Default.JSONIndent(title, indent, values...)
}

// Line inserts blank lines using default Notifier
// Visual separation utility
func Line(no int) { Default.Line(no) }

// Logf writes formatted log with timestamp using default Notifier
// Full-featured logging shortcut
func Logf(level LogLevel, f string, a ...any) { Default.Logf(level, f, a...) }

// Notice logs a message at Notice level using default Notifier
// Notable event reporting
func Notice(f string, a ...any) { Default.Notice(f, a...) }

// Printf writes plain message using default Notifier
// Minimal formatted output
func Printf(level LogLevel, f string, a ...any) { Default.Printf(level, f, a...) }

// Panic logs a message at Critical level using default Notifier and panics
// Convenience function for critical errors that should stop execution
func Panic(f string, a ...any) { Default.Panic(f, a...) }

// Robot displays ASCII robot using default Notifier
// Fun visual addition
func Robot(l LogLevel) { Default.Robot(l) }

// Success logs success message with checkmark using default Notifier
// Positive feedback shortcut
func Success(format string, args ...any) {
	Default.Success(format, args...)
}

// Warn logs a message at Warn level using default Notifier
// Warning notification shortcut
func Warn(f string, a ...any) { Default.Warn(f, a...) }

// With creates new Notifier with prefix using default Notifier
// Contextual logging setup
func With(prefix string) *Notifier { return Default.With(prefix) }

/* ========== Package Configuration ========== */

// ResetColors resets all colors to their default values
// Useful for restoring original color scheme
func ResetColors() {
	mu.Lock()
	defer mu.Unlock()
	for k, v := range defaultColors {
		colors[k] = v
	}
}

// ResetSymbols resets all symbols to their default values
// Restores original symbol set
func ResetSymbols() {
	mu.Lock()
	defer mu.Unlock()
	for k, v := range defaultSymbols {
		symbols[k] = v
	}
}

// SetColor sets custom color for specific log level
// Allows customization of level appearance
func SetColor(level LogLevel, color *color.Color) {
	mu.Lock()
	defer mu.Unlock()
	colors[level] = color
}

// SetSymbol sets custom symbol for specific log level
// Enables custom visual indicators
func SetSymbol(level LogLevel, symbol string) {
	mu.Lock()
	defer mu.Unlock()
	symbols[level] = symbol
}
