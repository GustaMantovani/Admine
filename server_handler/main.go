package main

import (
	"server_handler/cmd"
	"server_handler/internal/config"
)

func main() {
	config.OpenLogFile()
	config.CreateLogger()
	config.GetLogger().Info("========= STARTING APP ===========")
	cmd.Execute()
}
