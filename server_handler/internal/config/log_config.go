package config

import (
	"log"
	"log/slog"
	"os"
	"sync"
)

var logger *slog.Logger
var onceOpenFile sync.Once
var onceCreateLogger sync.Once
var logFile *os.File

func GetLogger() *slog.Logger {
	return logger
}

func CreateLogger() {
	onceCreateLogger.Do(func() {
		handler := slog.NewTextHandler(logFile, nil)
		logger = slog.New(handler)
	})
}

func OpenLogFile(name string) {
	onceOpenFile.Do(func() {
		file, err := os.OpenFile(name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			log.Println("Erro ao criar arquivo de log")
		}

		logFile = file
		file.WriteString("\n\n\n")
	})
}

func CloseLogFile() {
	logFile.Close()
}
