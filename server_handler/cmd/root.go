package cmd

import (
	"fmt"
	"log"
	"os"
	"server_handler/cmd/queue"

	"github.com/spf13/cobra"
)

var shortDescription = "Start minecraft server from a docker compose file."

var longDescription = `
Start minecraft server from a docker compose file.
The compose file must be specified in a YAML file in ~/.config/admine/server.yaml or in environment variables.
If the environment variables are not fully set, then the file is used to configure the handler.

server.yaml content:
serverName: "name-of-the-service-in-the-compose-file"
composeDirectory: "/compose/absolute/path.yaml"
host: "pubsub-host-address"
port: "pubsub-port"
senderChannel: "channel"
consumerChannel:
- "channel1"
- "channel2"

Environment variables:
SERVER_NAME "channel"
COMPOSE_DIRECTORY "/path"
CONSUMER_CHANNEL "channel1:channel2"
SENDER_CHANNEL "channel"
PUBSUB "pubsub-type"
HOST "pubsub-host-address"
PORT "pubsub-port"
`

var rootCmd = &cobra.Command{
	Use:   "",
	Short: shortDescription,
	Long:  longDescription,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Starting Admine server handler...")
		go queue.RunListenQueue()
		// Keep the main goroutine alive
		select {}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
