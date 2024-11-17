package cmd

import (
	"fmt"
	"log"
	"os"
	"server/handler/internal"
	"time"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/spf13/cobra"
)

var minecraftServer = internal.MinecraftServerContainerByCompose{}

var runLongDescription = `Up the server container by compose and continuosly monitors its status to ensure it stays up and running.
If the docker container is down, the program will up him again.

The program use env and a yaml config file to get the server info. 

The env vars are MINECRAFT_SERVER_SERVICE which refers to the service service name 
in the compose file and MINECRAFT_SERVER_DIRECTORY which refers to compose directory.

The config file is a yaml in ~/.config/admine/adhandler.yaml with the fields 'serviceName' and 'composeDirectory'.
The directory must be full name.

The program will first look for the env, if not defined will then look for the config file.

The program can receive the comand using args, which refers to the service name in the compose. 
The directory is the working directory in the shell`

var env bool
var file bool

var rootCmd = &cobra.Command{
	Use:   "Run",
	Short: "Up the server container and continuosly monitors its status to ensure it stays up and running",
	Long:  runLongDescription,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			configureWithArgs(args)
		} else if verifyEnvVars() {
			configureWithEnv()
		} else if verifyConfigFile() {
			configWithFile()
		} else {
			fmt.Println("Não foi possível obter as configurações do servidor")
			os.Exit(0)
		}

		if env && file {
			fmt.Println("Flag excludentes foram chamadas.")
			os.Exit(1)
		}

		if env {
			configureWithEnv()
		}

		if file {
			configureWithEnv()
		}

		fmt.Println("a: ", minecraftServer)

		for {
			minecraftServer.VerifyContainerAndUpIfDown()
			time.Sleep(1 * time.Second)
		}

	},
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

func configureWithArgs(args []string) {
	minecraftServer.SetContainerNameByServiceAndDirectory(args[0], getLocalDirectory())
}

func configureWithEnv() {
	serverName := os.Getenv("MINECRAFT_SERVER_SERVICE")
	directory := os.Getenv("MINECRAFT_SERVER_DIRECTORY")

	minecraftServer.SetContainerNameByServiceAndDirectory(serverName, directory)
	minecraftServer.UpMinecraftServerContainerByCompose()
}

func configureWithCfgFile() {}

func verifyEnvVars() bool {
	_, serverNameExists := os.LookupEnv("MINECRAFT_SERVER_SERVICE")
	_, directoryExists := os.LookupEnv("MINECRAFT_SERVER_DIRECTORY")

	return serverNameExists && directoryExists
}

func verifyConfigFile() bool {
	return internal.VerifyIfConfigFileExists()
}

func configWithFile() {
	configFileData, err := internal.GetConfigFileData()
	if err != nil {
		fmt.Println("Erro: ", err)
		os.Exit(1)
	}

	serverName := configFileData.ServerName
	directory := configFileData.ComposeDirectory

	minecraftServer.SetContainerNameByServiceAndDirectory(serverName, directory)
	minecraftServer.UpMinecraftServerContainerByCompose()
}

// FUNÇÕES A SEREM MOVIDAS PRA SEUS DEVIDOS DIRETÓRIOS
func getLocalDirectory() string {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Erro ao obter o diretório de trabalho: %v\n", err)
		return ""
	}

	return wd
}
