package config

import (
	"fmt"
	"os"
	"strings"
)

// Set a Config variable and if any value is an empty string returns false, because the config is not full set
func isEnvSetAndSetConfig(config *config) (bool, error) {
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

	var envVarsNotSet []string
	if config.ComposeContainerName == "" {
		envVarsNotSet = append(envVarsNotSet, "SERVER_NAME")
	}

	if config.ComposeAbsPath == "" {
		envVarsNotSet = append(envVarsNotSet, "COMPOSE_DIRECTORY")
	}

	if config.ConsumerChannel[len(config.ConsumerChannel)-1] == "" {
		envVarsNotSet = append(envVarsNotSet, "CONSUMER_CHANNEL")
	}

	if config.SenderChannel == "" {
		envVarsNotSet = append(envVarsNotSet, "SENDER_CHANNEL")
	}

	if config.SenderChannel == "" {
		envVarsNotSet = append(envVarsNotSet, "PUBSUB")
	}

	if config.SenderChannel == "" {
		envVarsNotSet = append(envVarsNotSet, "HOST")
	}

	if config.SenderChannel == "" {
		envVarsNotSet = append(envVarsNotSet, "PORT")
	}

	if len(envVarsNotSet) != 0 {
		return false, fmt.Errorf("env vars not set: [%s]", strings.Join(envVarsNotSet, " "))
	}

	return true, nil
}
