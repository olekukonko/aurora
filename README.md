# Aurora - Pretty Print Terminal Logging Library

Aurora is a lightweight Go library designed for beautiful, colorized, and structured terminal logging. It offers leveled logging with customizable symbols, automatic timestamps, and JSON pretty-printing capabilities, making it ideal for both development and production use.

## Features

- 🌈 **Colorized Output**: Distinct colors for each log level
- ⏱ **Timestamps**: Automatic inclusion of timestamps
- 📛 **Contextual Prefixes**: Add module or component-specific prefixes
- 📦 **JSON Pretty Printing**: Built-in support for structured data
- 🔒 **Thread-Safe**: Safe concurrent logging
- 🍬 **Sugar Functions**: Convenient global logging methods
- 🎨 **Customizable**: Adjust symbols and colors to your liking

## Installation

Install Aurora using:

```bash
go get github.com/olekukonko/aurora
```

## Usage

### Basic Logging

```go
package main

import "github.com/olekukonko/aurora"

func main() {
	// Simple logging with sugar functions
	aurora.Info("System initialized")
	aurora.Warn("Disk space running low")
	aurora.Error("Failed to connect to database")

	// Formatted logging
	aurora.Debug("User %s logged in with %d attempts", "alice", 3)
}
```

### Contextual Logging

```go
func main() {
	// Create a logger with a prefix for context
	dbLogger := aurora.With("database")

	dbLogger.Info("Connection pool created")
	dbLogger.Warn("Slow query detected")
}
```

### JSON Output

```go
func main() {
    aurora.Br(2)
    aurora.Robot(aurora.InfoLevel)
    data := map[string]interface{}{
        "user": "bob",
        "age":  32,
        "tags": []string{"admin", "premium"},
    }
    
    aurora.Br(2)
    aurora.JSON("User Profile", data)
}
```

### Advanced Usage

```go
package main

import (
	"bytes"
	"github.com/olekukonko/aurora"
)

func main() {
	// Custom logger with a buffer
	var buf bytes.Buffer
	logger := aurora.New(&buf)

	// Full log with timestamp
	logger.Logf(aurora.NoticeLevel, "System check completed")

	// Inline log (no timestamp)
	logger.Inlinef(aurora.DebugLevel, "Processing item %d", 42)

	// Plain output (no symbol or timestamp)
	logger.Printf(aurora.InfoLevel, "Simple message")
}
```

## Log Levels

Aurora supports the following log levels with default symbols and colors:

| Level      | Symbol | Color   | Usage                  |
|------------|--------|---------|------------------------|
| `Alert`    | ✭      | Blue    | Important alerts       |
| `Info`     | ✔      | Green   | Informational messages |
| `Error`    | ✘      | Red     | Error conditions       |
| `Notice`   | ⚑      | Yellow  | Notable events         |
| `Debug`    | ⧳      | Cyan    | Debug information      |
| `Warn`     | ⚠      | Magenta | Warning conditions     |
| `Critical` | ‼      | White   | Critical conditions    |

## Configuration

### Custom Output

The default logger writes to `os.Stdout`. Create custom loggers for different outputs:

```go
// Write to a file
file, _ := os.Create("app.log")
fileLogger := aurora.New(file)

// Write to multiple outputs
multiWriter := io.MultiWriter(os.Stdout, file)
multiLogger := aurora.New(multiWriter)
```

### Custom Symbols and Colors

Customize symbols and colors for specific log levels:

```go
// Change symbol and color for Debug level
aurora.SetSymbol(aurora.DebugLevel, "[DEBUG]")
aurora.SetColor(aurora.DebugLevel, color.New(color.FgHiWhite))
```

## Best Practices

- **Contextual Logging**: Use `With()` to create prefixed loggers for different components.
- **Level Usage**: Reserve `Critical` for severe issues and `Debug` for development verbosity.
- **Structured Data**: Use `JSON()` for logging complex data structures.
- **Production**: Combine with file or network writers for persistent logging.
- **Testing**: Capture output in tests using `bytes.Buffer`.

## Example Output

Here’s what Aurora’s output might look like in your terminal:

```

    .=._,=.
   ' ([o]) `
     _)e(_
.'c ."|_|". n`.
'--'  /_\  `--'
     (   )
    __) (__


[✔] 2025-03-25 01:23:45 PM System initialized
[⚠] 2025-03-25 01:23:46 PM Disk space running low
[⧳] TestJSON: JSON ↴↴
{
    "user": "bob",
    "age": 32,
    "tags": ["admin", "premium"]
}
```

*(Note: Colors not shown here; see actual terminal for full effect.)*

For a real screenshot, consider adding one to your repo:

```markdown
![Example terminal output](https://example.com/aurora-screenshot.png)
```

## Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/xyz`)
3. Commit your changes (`git commit -m "Add xyz"`)
4. Push to the branch (`git push origin feature/xyz`)
5. Open a Pull Request

## Dependencies

- [fatih/color](https://github.com/fatih/color) - Terminal color support
- [mattes/go-asciibot](https://github.com/mattes/go-asciibot) - ASCII art generation
- [nwidger/jsoncolor](https://github.com/nwidger/jsoncolor) - JSON colorization

## License

MIT License - see [LICENSE](LICENSE) for details.