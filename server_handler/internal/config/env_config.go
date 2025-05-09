package config

import "os"

// Set a Config variable and if any value is an empty string returns false, because the config is not full set
func isEnvSetAndSetConfig(config *Config) bool {
	config.ComposeContainerName = os.Getenv("SERVER_NAME")
	config.ComposeAbsPath = os.Getenv("COMPOSE_DIRECTORY")
	config.ConsumerChannel = os.Getenv("CONSUMER_CHANNEL")
	config.SenderChannel = os.Getenv("SENDER_CHANNEL")

	if config.ComposeContainerName == "" || config.ComposeAbsPath == "" || config.ConsumerChannel == "" || config.SenderChannel == "" {
		return false
	}

	return true
}
