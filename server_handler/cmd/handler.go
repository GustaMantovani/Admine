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

var rootCmd = &cobra.Command{
	Use:   "Run",
	Short: "Up the server container and continuosly monitors its status to ensure it stays up and running",
	Long: `Up the server container by compose and continuosly monitors its status to ensure it stays up and running.
  If the docker container is down, the program will up him again`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := docker.NewClientFromEnv()
		if err != nil {
			log.Fatal(err)
		}

		minecraftServerContainerByCompose := internal.NewMinecraftServerContainerByCompose(client, "mine_server", "/home/andre/pgm/pessoal/minecraft-server-on-docker/")
		minecraftServerContainerByCompose.UpMinecraftServerContainerByCompose()

		for {
			minecraftServerContainerByCompose.VerifyContainerAndUpIfDown()
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
