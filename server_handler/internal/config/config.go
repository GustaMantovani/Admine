package config

import (
	"log"
	"path"
	"sync"
)

type Config struct {
	ComposeAbsPath       string
	ComposeContainerName string
	ConsumerChannel      string
	SenderChannel        string
}

var instance *Config
var once sync.Once

// Obter instancia Singleton da configuração do servidor
func GetInstance() *Config {
	once.Do(func() {
		configFile, err := GetConfigFileData()
		if err != nil {
			log.Println("erro get config file data: ", err.Error())
		}

		composeAbsPath := configFile.ComposeDirectory + "/" + "docker-compose.yaml"
		containerName := path.Base(configFile.ComposeDirectory) + "-" + configFile.ServerName + "-1"

		instance = &Config{
			ComposeAbsPath:       composeAbsPath,
			ConsumerChannel:      configFile.ConsumerChannel,
			SenderChannel:        configFile.SenderChannel,
			ComposeContainerName: containerName,
		}
	})

	return instance
}
