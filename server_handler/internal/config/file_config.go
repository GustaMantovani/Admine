package config

import (
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

// Return a ConfigFile. It takes the data from config file in "~/.config/admine/server.yaml"
func GetConfigFileData() (ConfigFile, error) {
	var configFile ConfigFile

	configFilePath, err := getConfigFilePath()

	if err != nil {
		return configFile, err
	}

	file, err := os.Open(configFilePath)

	if err != nil {
		return configFile, err
	}

	decoder := yaml.NewDecoder(file)

	err = decoder.Decode(&configFile)

	if err != nil {
		return configFile, err
	}

	file.Close()

	return configFile, nil
}

// Return the yaml config file in user home directory
func getConfigFilePath() (string, error) {
	home, err := os.UserHomeDir()

	if err != nil {
		return "", err
	}

	configFilePath := home + "/.config/admine/server.yaml"

	return configFilePath, nil
}
