package pkg

import (
	"io"
	"log"
	"os"
)

// Logger is the global logging instance
var Logger *CustomLogger

// CustomLogger wraps the standard logger with level methods
type CustomLogger struct {
	internal *log.Logger
}

// Setup initializes the global logger with console and file output
func Setup(logFile string) (*CustomLogger, error) {
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	multiWriter := io.MultiWriter(os.Stdout, file)

	l := &CustomLogger{
		internal: log.New(multiWriter, "", log.Ldate|log.Ltime|log.Lshortfile),
	}

	Logger = l
	return l, nil
}

// Info logs informational messages
func (l *CustomLogger) Info(format string, v ...interface{}) {
	l.internal.Printf("[INFO] "+format, v...)
}

// Error logs error messages
func (l *CustomLogger) Error(format string, v ...interface{}) {
	l.internal.Printf("[ERROR] "+format, v...)
}

// Debug logs debug messages
func (l *CustomLogger) Debug(format string, v ...interface{}) {
	l.internal.Printf("[DEBUG] "+format, v...)
}
