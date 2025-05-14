package config

import (
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

	if config.ComposeContainerName == "" || config.ComposeAbsPath == "" || len(config.ConsumerChannel) == 0 || config.SenderChannel == "" {
		return false
	}

	return true
}
