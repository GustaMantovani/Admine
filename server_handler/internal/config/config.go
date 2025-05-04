package config

import (
	"log"
	"sync"
)

type Config struct {
	ComposeAbsPath  string
	ConsumerChannel string
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

		instance = &Config{
			ComposeAbsPath:  composeAbsPath,
			ConsumerChannel: configFile.ConsumerChannel,
		}
	})

	return instance
}
