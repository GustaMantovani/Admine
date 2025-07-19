package config

import (
	"log"
	"log/slog"
	"os"
	"path/filepath"
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

func OpenLogFile() {
	logDir := verifyLogDir()

	onceOpenFile.Do(func() {
		file, err := os.OpenFile(logDir+"/server_handler.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
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

func GetLogFile() *os.File {
	return logFile
}

func verifyLogDir() string {
	homedir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Error getting home dir. Verify if $HOME env var is set.", err.Error())
	}

	admineStateDir := filepath.Join(homedir, ".local", "state", "admine")

	if _, err := os.Stat(admineStateDir); os.IsNotExist(err) {
		err := os.MkdirAll(admineStateDir, 0755)
		if err != nil {
			log.Fatal("Error creating admine state dir: ", err.Error())
		}
	}

	return admineStateDir
}
