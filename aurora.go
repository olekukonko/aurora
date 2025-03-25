package aurora

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/nwidger/jsoncolor"
	"io"
	"os"
	"sync"
	"time"
)

// LogLevel defines the severity of the log message.
type LogLevel int

// Log level constants.
const (
	AlertLevel LogLevel = iota
	InfoLevel
	ErrorLevel
	NoticeLevel
	DebugLevel
	WarnLevel
	CriticalLevel
)

// symbols maps each LogLevel to a specific symbol.
var symbols = map[LogLevel]string{
	AlertLevel:    "[✭]",
	InfoLevel:     "[✔]",
	ErrorLevel:    "[✘]",
	NoticeLevel:   "[⚑]",
	DebugLevel:    "[⧳]",
	WarnLevel:     "[⚠]",
	CriticalLevel: "[‼]",
}

// colors maps each LogLevel to a color function for terminal output.
var colors = map[LogLevel]*color.Color{
	AlertLevel:    color.New(color.FgHiBlue),
	InfoLevel:     color.New(color.FgHiGreen),
	ErrorLevel:    color.New(color.FgHiRed),
	NoticeLevel:   color.New(color.FgHiYellow),
	DebugLevel:    color.New(color.FgHiCyan),
	WarnLevel:     color.New(color.FgHiMagenta),
	CriticalLevel: color.New(color.FgHiWhite),
}

// Notifier is responsible for synchronized, level-aware, and colorized logging.
// It optionally supports a prefix to add contextual information to each log message.
type Notifier struct {
	mu     *sync.Mutex // Pointer to a mutex for thread-safe writes.
	output io.Writer   // Destination for log output.
	prefix string      // Optional prefix to add context to log messages.
}

// New creates a new Notifier that writes to the given io.Writer.
// If the writer is nil, os.Stdout is used.
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

// With returns a new Notifier that shares the same output and mutex,
// but prepends the provided prefix to each log message. This enables
// contextual or module-specific logging.
func (n *Notifier) With(prefix string) *Notifier {
	newPrefix := prefix
	if n.prefix != "" {
		newPrefix = fmt.Sprintf("%s %s", n.prefix, prefix)
	}
	return &Notifier{
		mu:     n.mu, // Share the same mutex for synchronized writes.
		output: n.output,
		prefix: newPrefix,
	}
}

// formatWithPrefix prepends the logger's prefix to the message if one is set.
func (n *Notifier) formatWithPrefix(msg string) string {
	if n.prefix != "" {
		return fmt.Sprintf("[%s] %s", n.prefix, msg)
	}
	return msg
}

// Logf writes a formatted log line including a level symbol, a timestamp,
// and a message. The message is colorized according to its log level and
// includes any prefix set on the Notifier.
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

// Inlinef writes a single-line log without a timestamp, useful for inline messages.
// The log line will include the level symbol and any prefix.
func (n *Notifier) Inlinef(level LogLevel, format string, args ...any) {
	n.mu.Lock()
	defer n.mu.Unlock()

	symbol := symbols[level]
	msg := fmt.Sprintf(format, args...)
	msg = n.formatWithPrefix(msg)
	line := fmt.Sprintf("%s %s\n", symbol, msg)

	colors[level].Fprint(n.output, line)
}

// Printf writes a plain log message without a timestamp or level symbol,
// but it still prepends the prefix if set.
func (n *Notifier) Printf(level LogLevel, format string, args ...any) {
	n.mu.Lock()
	defer n.mu.Unlock()

	msg := fmt.Sprintf(format, args...)
	msg = n.formatWithPrefix(msg)
	line := fmt.Sprintf("%s\n", msg)

	colors[level].Fprint(n.output, line)
}

// JSON pretty-prints the provided values in JSON format with colorization,
// preceded by a title line. Each value is printed on a separate line.
func (n *Notifier) JSON(title string, values ...any) {
	n.Inlinef(DebugLevel, "%s: JSON ↴↴", title)

	n.mu.Lock()
	defer n.mu.Unlock()

	for _, v := range values {
		data, err := jsoncolor.Marshal(v)
		if err != nil {
			n.Logf(ErrorLevel, "failed to marshal JSON: %v", err)
			continue
		}
		n.output.Write(data)
		n.output.Write([]byte{'\n'})
	}
	n.output.Write([]byte{'\n'})
}

// Convenience methods for logging at specific levels.

// Alert logs a message at Alert level.
func (n *Notifier) Alert(f string, a ...any) { n.Inlinef(AlertLevel, f, a...) }

// Info logs a message at Info level.
func (n *Notifier) Info(f string, a ...any) { n.Inlinef(InfoLevel, f, a...) }

// Error logs a message at Error level.
func (n *Notifier) Error(f string, a ...any) { n.Inlinef(ErrorLevel, f, a...) }

// Notice logs a message at Notice level.
func (n *Notifier) Notice(f string, a ...any) { n.Inlinef(NoticeLevel, f, a...) }

// Debug logs a message at Debug level.
func (n *Notifier) Debug(f string, a ...any) { n.Inlinef(DebugLevel, f, a...) }

// Warn logs a message at Warn level.
func (n *Notifier) Warn(f string, a ...any) { n.Inlinef(WarnLevel, f, a...) }

// Critical logs a message at Critical level.
func (n *Notifier) Critical(f string, a ...any) { n.Inlinef(CriticalLevel, f, a...) }

// Default is a global Notifier instance that writes to os.Stdout for convenience.
var Default = New(os.Stdout)

// Global sugar functions for easier usage with the default Notifier.

// Alert logs a message at Alert level using the default Notifier.
func Alert(f string, a ...any) { Default.Alert(f, a...) }

// Info logs a message at Info level using the default Notifier.
func Info(f string, a ...any) { Default.Info(f, a...) }

// Error logs a message at Error level using the default Notifier.
func Error(f string, a ...any) { Default.Error(f, a...) }

// Notice logs a message at Notice level using the default Notifier.
func Notice(f string, a ...any) { Default.Notice(f, a...) }

// Debug logs a message at Debug level using the default Notifier.
func Debug(f string, a ...any) { Default.Debug(f, a...) }

// Warn logs a message at Warn level using the default Notifier.
func Warn(f string, a ...any) { Default.Warn(f, a...) }

// Critical logs a message at Critical level using the default Notifier.
func Critical(f string, a ...any) { Default.Critical(f, a...) }

// JSON logs JSON formatted values with a title using the default Notifier.
func JSON(t string, v ...any) { Default.JSON(t, v...) }

// With returns a new Notifier with the given prefix.
func With(prefix string) *Notifier { return Default.With(prefix) }

// Printf writes a formatted message using the default Notifier.
func Printf(level LogLevel, f string, a ...any) { Default.Printf(level, f, a...) }

// Inlinef writes a single-line log without a timestamp, useful for inline messages.
func Inlinef(level LogLevel, f string, a ...any) { Default.Inlinef(level, f, a...) }

// Logf writes a formatted log line including a level symbol, a timestamp,
func Logf(level LogLevel, f string, a ...any) { Default.Logf(level, f, a...) }
