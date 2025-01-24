package logger

import (
	"log"
	"os"
	"path/filepath"

	"github.com/suryamp/receipt-processor/internal/testing"
)

var (
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
)

func Init() error {
	// For tests, just log to stdout
	if testing.Testing() {
		InfoLogger = log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile)
		ErrorLogger = log.New(os.Stdout, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile)
		return nil
	}

	// Create logs directory if it doesn't exist
	if err := os.MkdirAll("logs", 0755); err != nil {
		return err
	}

	// Open log file with append mode
	logFile, err := os.OpenFile(
		filepath.Join("logs", "receipt-processor.log"),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)
	if err != nil {
		return err
	}

	// Create multiple loggers for different levels
	InfoLogger = log.New(logFile, "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(logFile, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile)

	return nil
}
