package logger

import (
	"io"
	"log/slog"
	"os"
)

var Log *slog.Logger

// InitLogger initializes the global logger with custom output and level
func initLogger(output io.Writer, level slog.Level) {
	opts := &slog.HandlerOptions{
		Level: level,
	}

	handler := slog.NewJSONHandler(output, opts).WithAttrs([]slog.Attr{
		slog.String("service", "qrkit"),
	})
	Log = slog.New(handler)
}

// InitSTDIOLogger sets up logging for STDIO mode
// Only logs fatal errors to stderr to avoid corrupting stdout protocol
func InitSTDIOLogger() {
	initLogger(os.Stderr, slog.LevelError)
}

// InitHTTPLogger sets up logging for HTTP mode
// Logs everything (debug and above) to stdout
func InitHTTPLogger() {
	initLogger(os.Stdout, slog.LevelDebug)
}
