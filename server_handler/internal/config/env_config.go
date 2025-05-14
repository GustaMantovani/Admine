package config

import (
	"log"
	"os"
	"strings"
)

// Set a Config variable and if any value is an empty string returns false, because the config is not full set
func isEnvSetAndSetConfig(config *Config) bool {
	var channels []string

	envChannels := os.Getenv("CONSUMER_CHANNEL")
	parts := strings.Split(envChannels, ":")
	channels = append(channels, parts...)

	config.ComposeContainerName = os.Getenv("SERVER_NAME")
	config.ComposeAbsPath = os.Getenv("COMPOSE_DIRECTORY")
	config.ConsumerChannel = channels
	config.SenderChannel = os.Getenv("SENDER_CHANNEL")
	config.Pubsub = os.Getenv("PUBSUB")

	log.Println(config)

	if config.ComposeContainerName == "" || config.ComposeAbsPath == "" || len(config.ConsumerChannel) == 0 || config.SenderChannel == "" || config.Pubsub == "" {
		return false
	}

	return true
}
