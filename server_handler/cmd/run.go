package cmd

import (
	"fmt"
	"log"
	"os"
	"server/handler/internal"
	"server/handler/pubsub"
	"time"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/spf13/cobra"
)

var minecraftServer = internal.MinecraftServerContainerByCompose{}
var subscriber = pubsub.RedisPubSubSubscriber{}

var runLongDescription = `Up the server container by compose and continuosly monitors its status to ensure it stays up and running.
If the docker container is down, the program will up him again.

The program use env and a yaml config file to get the server info. 

The env vars are MINECRAFT_SERVER_SERVICE which refers to the service service name 
in the compose file and MINECRAFT_SERVER_DIRECTORY which refers to compose directory.

The config file is a yaml in ~/.config/admine/server.yaml with the fields 'serviceName' and 'composeDirectory'.
The directory must be full name.

The program will first look for the env, if not defined will then look for the config file.

The program can receive the comand using args, which refers to the service name in the compose. 
The directory is the working directory in the shell`

var env bool
var file bool

// Roda o server handler
func runServerHandler() {
	iniciado := false
	config := pubsub.GetConfigServerChannelFromDotEnv("REDIS_SERVER_CHANNEL")
	subscriber := pubsub.CreateSubscriber(config.Addr, config.Channel)

	var isUp bool
	for {
		_, isUp = minecraftServer.VerifyContainerAndUpIfDown()
		if isUp == true && iniciado == false {
			subscriber.SendMessage(internal.ConvertMessageToJson("server_up", minecraftServer.ContainerName))
		} else if isUp == false && iniciado == true {
			subscriber.SendMessage(internal.ConvertMessageToJson("server_down", minecraftServer.ContainerName))
		}

		iniciado = isUp

		time.Sleep(1 * time.Second)
	}
}

func configureMinecraftServer(args []string) {
	if len(args) > 0 {
		minecraftServer.ConfigureWithArgs(args)
	} else if verifyEnvVars() {
		minecraftServer.ConfigureWithEnv()
	} else if verifyConfigFile() {
		minecraftServer.ConfigureWithFile()
	} else {
		fmt.Println("Não foi possível obter as configurações do servidor. Elas não estão definidas.")
		os.Exit(0)
	}

	if env && file {
		fmt.Println("Flag excludentes foram chamadas.")
		os.Exit(1)
	}

	if env {
		minecraftServer.ConfigureWithEnv()
	}

	if file {
		minecraftServer.ConfigureWithFile()
	}

}

func runRootCmd(cmd *cobra.Command, args []string) {
	configureMinecraftServer(args)
	go runServerHandler()
	go internal.RunCommandHandler(minecraftServer.ContainerName)

	for {
	}
}

var rootCmd = &cobra.Command{
	Use:   "Run",
	Short: "Up the server container and continuosly monitors its status to ensure it stays up and running",
	Long:  runLongDescription,
	Run:   runRootCmd,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func initMineServerDockerClient() {
	client, err := docker.NewClientFromEnv()
	if err != nil {
		log.Fatal(err)
	}

	minecraftServer.DockerClient = client
}

func init() {
	cobra.OnInitialize(initMineServerDockerClient)
	rootCmd.Flags().BoolVarP(&env, "env", "e", false, "use environment variables to read server info")
	rootCmd.Flags().BoolVarP(&file, "file", "f", false, "use file to read server info")
}

func verifyEnvVars() bool {
	_, serverNameExists := os.LookupEnv("MINECRAFT_SERVER_SERVICE")
	_, directoryExists := os.LookupEnv("MINECRAFT_SERVER_DIRECTORY")

	return serverNameExists && directoryExists
}

func verifyConfigFile() bool {
	return internal.VerifyIfConfigFileExists()
}
