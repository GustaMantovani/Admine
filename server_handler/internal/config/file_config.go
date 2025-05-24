package config

import (
	"log"
	"os"

	yaml "gopkg.in/yaml.v3"
)

type ConfigFile struct {
	ServerName       string   `yaml:"serverName"`
	ComposeDirectory string   `yaml:"composeDirectory"`
	ConsumerChannels []string `yaml:"consumerChannels"`
	SenderChannel    string   `yaml:"senderChannel"`
	Pubsub           string   `yaml:"pubsub"`
	Host             string   `yaml:"host"`
	Port             string   `yaml:"port"`
}

// GetConfigFileData returns a ConfigFile. It takes the data from config file in "~/.config/admine/server.yaml"
func GetConfigFileData() (ConfigFile, error) {
	var configFile ConfigFile

	configFilePath := getConfigFilePath()
	log.Printf("Reading config file from: %s", configFilePath)

	file, err := os.Open(configFilePath)

	if err != nil {
		log.Printf("Error opening config file: %v", err)
		return configFile, err
	}

	decoder := yaml.NewDecoder(file)

	err = decoder.Decode(&configFile)

	if err != nil {
		log.Printf("Error decoding config file: %v", err)
		return configFile, err
	}

	file.Close()
	log.Printf("Successfully loaded config file: %+v", configFile)

	return configFile, nil
}

// getConfigFilePath returns the yaml config file path in user home directory
func getConfigFilePath() string {
	home, err := os.UserHomeDir()

	if err != nil {
		log.Printf("Error finding user home directory: %v", err)
		return ""
	}

	configFilePath := home + "/.config/admine/server.yaml"
	log.Printf("Config file path: %s", configFilePath)

	return configFilePath
}
