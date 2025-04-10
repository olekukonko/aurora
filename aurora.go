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
	NoLevel
)

// Default symbols for each log level
var defaultSymbols = map[LogLevel]string{
	AlertLevel:    "[✭]",
	InfoLevel:     "[✔]",
	ErrorLevel:    "[✘]",
	NoticeLevel:   "[⚑]",
	DebugLevel:    "[⧳]",
	WarnLevel:     "[⚠]",
	CriticalLevel: "[‼]",
	NoLevel:       " ",
}

// Default colors for each log level
var defaultColors = map[LogLevel]*color.Color{
	AlertLevel:    color.New(color.FgHiBlue),
	InfoLevel:     color.New(color.FgHiGreen),
	ErrorLevel:    color.New(color.FgHiRed),
	NoticeLevel:   color.New(color.FgHiYellow),
	DebugLevel:    color.New(color.FgHiCyan),
	WarnLevel:     color.New(color.FgHiMagenta),
	CriticalLevel: color.New(color.FgHiWhite),
	NoLevel:       color.New(color.FgHiBlack),
}

// Package-level customization
var (
	symbols = make(map[LogLevel]string)
	colors  = make(map[LogLevel]*color.Color)
	mu      sync.RWMutex
)

func init() {
	// Initialize with defaults
	ResetSymbols()
	ResetColors()
}

// SetSymbol sets a custom symbol for a specific log level
func SetSymbol(level LogLevel, symbol string) {
	mu.Lock()
	defer mu.Unlock()
	symbols[level] = symbol
}

// SetColor sets a custom color for a specific log level
func SetColor(level LogLevel, color *color.Color) {
	mu.Lock()
	defer mu.Unlock()
	colors[level] = color
}

// ResetSymbols resets all symbols to their default values
func ResetSymbols() {
	mu.Lock()
	defer mu.Unlock()
	for k, v := range defaultSymbols {
		symbols[k] = v
	}
}

// ResetColors resets all colors to their default values
func ResetColors() {
	mu.Lock()
	defer mu.Unlock()
	for k, v := range defaultColors {
		colors[k] = v
	}
}

// Notifier remains the same as before...
type Notifier struct {
	mu     *sync.Mutex
	output io.Writer
	prefix string
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

// Br inserts a single blank line into the output.
// It does not include any timestamp or log level indicator,
// but it includes the configured prefix if one is set.
func (n *Notifier) Br() {
	n.Line(1)
}

// Line inserts the specified number of blank lines into the output.
// It omits timestamps and log level symbols, but includes the prefix if configured.
func (n *Notifier) Line(count int) {
	n.mu.Lock()
	defer n.mu.Unlock()
	colors[NoLevel].Fprint(n.output, "%s", strings.Repeat("\n", count))
}

// Robot writes a plain log message without a timestamp or level symbol,
// but it still prepends the prefix if set.
func (n *Notifier) Robot(level LogLevel) {
	n.mu.Lock()
	defer n.mu.Unlock()
	colors[level].Fprint(n.output, fmt.Sprintf("%s\n", asciibot.Random()))
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

// Robot sends friendly robot
func Robot(l LogLevel) { Default.Robot(l) }

// Br sends line Break
func Br(no int) { Default.Printf(InfoLevel, "%s", strings.Repeat("\n", no)) }

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
