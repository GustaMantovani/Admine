package handler

import (
	"errors"
	"server_handler/internal/server"
)

func ManageCommand(command string) error {
	if command == "start" {
		server.StartServerCompose()
	} else if command == "stop" {
		server.StopServerCompose()
	} else {
		return errors.New("invalid command")
	}

	return nil
}
