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

// StartServerCompose starts the minecraft server using docker compose
func StartServerCompose() error {
	compose := configureCompose()

	log.Println("Starting minecraft server using docker compose...")
	upCh, err := compose.Up(&client.UpOptions{Detach: true}, nil)

	if err != nil {
		log.Printf("Failed to start docker compose: %v", err)
		return fmt.Errorf("error starting compose: %w", err)
	}

	err = <-upCh

	if err != nil {
		log.Printf("Docker compose operation failed: %v", err)
		return fmt.Errorf("error in compose: %w", err)
	}

	log.Println("Docker compose started successfully")
	return nil
}

// StopServerCompose stops the minecraft server using docker compose
func StopServerCompose() {
	compose := configureCompose()

	log.Println("Stopping minecraft server using docker compose...")
	downCh, err := compose.Down(&client.DownOptions{}, nil)

	if err != nil {
		log.Printf("Failed to stop docker compose: %v", err)
		log.Fatal("error stopping compose: ", err.Error())
	}

	err = <-downCh

	if err != nil {
		log.Printf("Docker compose down operation failed: %v", err)
		log.Fatal("error in compose: ", err.Error())
	}

	log.Println("Docker compose stopped successfully")
}
