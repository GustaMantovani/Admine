package config

import (
	"log"
	"os"
	"strings"
)

// isEnvSetAndSetConfig sets a Config variable and returns false if any value is an empty string,
// because the config is not fully set
func isEnvSetAndSetConfig(config *config) bool {
	var channels []string

	envChannels := os.Getenv("CONSUMER_CHANNEL")
	parts := strings.Split(envChannels, ":")
	channels = append(channels, parts...)

	config.ComposeContainerName = os.Getenv("SERVER_NAME")
	config.ComposeAbsPath = os.Getenv("COMPOSE_DIRECTORY")
	config.ConsumerChannel = channels
	config.SenderChannel = os.Getenv("SENDER_CHANNEL")
	config.Pubsub = os.Getenv("PUBSUB")
	config.Host = os.Getenv("HOST")
	config.Port = os.Getenv("PORT")

	log.Printf("Environment configuration loaded: %+v", config)

	if config.ComposeContainerName == "" || config.ComposeAbsPath == "" || len(config.ConsumerChannel) == 0 || config.SenderChannel == "" || config.Pubsub == "" || config.Host == "" || config.Port == "" {
		log.Println("Environment configuration is incomplete")
		return false
	}

	log.Println("Environment configuration is complete")
	return true
}
