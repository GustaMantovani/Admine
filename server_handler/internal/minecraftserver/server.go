package minecraftserver

import (
	"fmt"
	"log"
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

	upCh, err := compose.Up(&client.UpOptions{Detach: true}, nil)

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
func StopServerCompose() {
	compose := configureCompose()

	downCh, err := compose.Down(&client.DownOptions{}, nil)

	if err != nil {
		log.Fatal("error stopping compose: ", err.Error())
	}

	err = <-downCh

	if err != nil {
		log.Fatal("error in compose: ", err.Error())
	}
}
