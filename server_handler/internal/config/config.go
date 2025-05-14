package config

import (
	"log"
	"path"
	"sync"
)

type Config struct {
	ComposeAbsPath       string
	ComposeContainerName string
	ConsumerChannel      []string
	SenderChannel        string
	Pubsub               string
}

var instance *Config
var once sync.Once

/*
Get the Singleton instance of the server configuration.

Checks whether it is possible to fetch ddata from a configuration file
or environment variables. If not, it closes the program.
*/
func GetInstance() *Config {
	once.Do(func() {
		instance = &Config{}
		configFile, err := GetConfigFileData()

		if err != nil {
			log.Println("Could not fetch data from configuration file. Error: ", err.Error())

			if !isEnvSetAndSetConfig(instance) {
				log.Fatalln("Coult not fetch data from env vars too. Closing program.")
			}

			return
		}

		composeAbsPath := configFile.ComposeDirectory + "/" + "docker-compose.yaml"
		containerName := path.Base(configFile.ComposeDirectory) + "-" + configFile.ServerName + "-1"

		instance = &Config{
			ComposeAbsPath:       composeAbsPath,
			ConsumerChannel:      configFile.ConsumerChannels,
			SenderChannel:        configFile.SenderChannel,
			ComposeContainerName: containerName,
			Pubsub:               configFile.Pubsub,
		}
	})

	return instance
}
