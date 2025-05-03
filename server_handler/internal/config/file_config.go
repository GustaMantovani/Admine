package config

import (
	"log"
	"os"

	yaml "gopkg.in/yaml.v3"
)

type ConfigFile struct {
	ServerName       string `yaml:"serverName"`
	ComposeDirectory string `yaml:"composeDirectory"`
}

// Return a ConfigFile. It takes the data from config file in "~/.config/admine/server.yaml"
func GetConfigFileData() (ConfigFile, error) {
	var configFile ConfigFile

	configFilePath := getConfigFilePath()

	file, err := os.Open(configFilePath)

	if err != nil {
		log.Println("error opening config file: ", err.Error())
		return configFile, err
	}

	decoder := yaml.NewDecoder(file)

	err = decoder.Decode(&configFile)

	if err != nil {
		log.Println("error decoding config file: ", err.Error())
		return configFile, err
	}

	file.Close()

	return configFile, nil
}

// Return the yaml config file in user home directory
func getConfigFilePath() string {
	home, err := os.UserHomeDir()

	if err != nil {
		log.Println("Error finding home user directory: ", err.Error())
		return ""
	}

	configFilePath := home + "/.config/admine/server.yaml"

	return configFilePath
}
