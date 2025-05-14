package cmd

import (
	"fmt"
	"os"
	"server_handler/cmd/queue"

	"github.com/spf13/cobra"
)

var shortDescription = "Up minecraft server from a docker compose file."

var longDescription = `Up minecraft server from a docker compose file.
The compose file must be specified in a YAML file in ~/.config/admine/server.yaml or in environment variables.
If the env vars is not fully set, then the file is used to configure the handler.

server.yaml content
serverName: "name-of-the-service-in-the-compose-file"
composeDirectory: "/compose/absolute/path.yaml"
consumerChannel: "channel-who-receives-commands"
senderChannel: "channel-who-send-responses"

env vars:
SERVER_NAME
COMPOSE_DIRECTORY
CONSUMER_CHANNEL
SENDER_CHANNEL
`

var rootCmd = &cobra.Command{
	Use:   "",
	Short: shortDescription,
	Long:  longDescription,
	Run: func(cmd *cobra.Command, args []string) {
		go queue.RunListenQueue()
		for {
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
