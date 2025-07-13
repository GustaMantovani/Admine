package config

import (
	"os"
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
Get the Singleton instance of the server configuration.

Checks whether it is possible to fetch data from a configuration file
or environment variables. If not, it closes the program.
*/
func GetInstance() *config {
	once.Do(func() {
		instance = &config{}
		configFile, err := GetConfigFileData()

		if err != nil {
			GetLogger().Warn("Could not fetch data from configuration file: " + err.Error())

			_, err := isEnvSetAndSetConfig(instance)
			if err != nil {
				GetLogger().Warn("Could not fetch data from env vars too: " + err.Error())
				GetLogger().Warn("Closing app because its not configured")
				os.Exit(1)
			}

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
	})

	return instance
}
