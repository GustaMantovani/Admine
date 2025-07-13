package minecraftserver

import (
	"fmt"
	"server_handler/internal/config"

	"github.com/harrim91/docker-compose-go/client"
)

func configureCompose() *client.ComposeClient {
	config := config.GetInstance()
	return client.New(
		&client.GlobalOptions{
			Files: []string{
				config.ComposeAbsPath,
			},
		},
	)
}

// Inicia o servidor
func StartServerCompose() error {
	compose := configureCompose()

	upCh, err := compose.Up(&client.UpOptions{Detach: false}, config.GetLogFile())

	if err != nil {
		return fmt.Errorf("error starting compose: %w", err)
	}

	err = <-upCh

	if err != nil {
		return fmt.Errorf("error in compose: %w", err)
	}

	return nil
}

// Parar o servidor
func StopServerCompose() error {
	compose := configureCompose()

	downCh, err := compose.Down(&client.DownOptions{}, nil)

	if err != nil {
		return err
	}

	err = <-downCh

	if err != nil {
		return err
	}

	return nil
}
