package pkg

import (
	"io"
	"log/slog"
	"os"
	"strings"
)

// Setup initializes the global slog with console and file output
func Setup(logFile string, logLevel string) error {
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	multiWriter := io.MultiWriter(os.Stdout, file)

	opts := &slog.HandlerOptions{
		Level:     defineLogLevel(logLevel),
		AddSource: true,
	}

	handler := slog.NewTextHandler(multiWriter, opts)
	l := slog.New(handler)

	// Set global slog
	slog.SetDefault(l)

	return nil
}

func defineLogLevel(ll string) slog.Leveler {
	switch strings.ToUpper(ll) {
	case "DEBUG":
		return slog.LevelDebug.Level()
	case "INFO":
		return slog.LevelInfo.Level()
	case "WARN":
		return slog.LevelWarn.Level()
	case "ERROR":
		return slog.LevelError.Level()
	default:
		return slog.LevelInfo.Level()
	}
}
