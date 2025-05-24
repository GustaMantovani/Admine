package config

import (
	"log"
	"path"
	"sync"
)

type config struct {
	ComposeAbsPath       string
	ComposeContainerName string
	ConsumerChannel      []string
	SenderChannel        string
	Pubsub               string
	Host                 string
	Port                 string
}

var instance *config
var once sync.Once

/*
GetInstance returns the Singleton instance of the server configuration.

Checks whether it is possible to fetch data from a configuration file
or environment variables. If not, it closes the program.
*/
func GetInstance() *config {
	once.Do(func() {
		log.Println("Initializing configuration...")
		instance = &config{}
		configFile, err := GetConfigFileData()

		if err != nil {
			log.Printf("Could not fetch data from configuration file. Error: %v", err)

			if !isEnvSetAndSetConfig(instance) {
				log.Fatalln("Could not fetch data from environment variables either. Closing program.")
			}

			log.Println("Using environment variables for configuration")
			return
		}

		composeAbsPath := configFile.ComposeDirectory + "/" + "docker-compose.yaml"
		containerName := path.Base(configFile.ComposeDirectory) + "-" + configFile.ServerName + "-1"

		instance = &config{
			ComposeAbsPath:       composeAbsPath,
			ConsumerChannel:      configFile.ConsumerChannels,
			SenderChannel:        configFile.SenderChannel,
			ComposeContainerName: containerName,
			Pubsub:               configFile.Pubsub,
			Host:                 configFile.Host,
			Port:                 configFile.Port,
		}

		log.Printf("Configuration initialized from file: %+v", instance)
	})

	return instance
}
