package config

import (
	"log"
	"sync"
)

type Config struct {
	ComposeAbsPath string
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
		log.Println("Compose Abs Path: ", composeAbsPath)
		instance = &Config{ComposeAbsPath: composeAbsPath}
	})

	return instance
}
