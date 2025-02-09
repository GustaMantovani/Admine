package cmd

import (
	"fmt"
	"log"
	"os"
	"time"

	commandhandler "server/handler/internal/command-handler"
	"server/handler/internal/health-checker"
	minecraftserver "server/handler/internal/minecraft-server"
	serverhandler "server/handler/internal/server-handler"

	"github.com/spf13/cobra"
)

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

func runRootCmd(cmd *cobra.Command, args []string) {
	ms := minecraftserver.NewMinecraftServerContainerByCompose()
	ms.ConfigureMinecraftServer(env, file, args)
	log.Println(ms.ContainerName)
	time.Sleep(time.Second * 2)
	go serverhandler.RunServerHandler(ms)
	go healthchecker.RunHealthChecker(ms)
	go commandhandler.RunCommandHandler(ms.ContainerName)

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

// func verifyConfigFile() bool {
// 	return internal.VerifyIfConfigFileExists()
// }
