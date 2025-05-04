package cmd

import (
	"fmt"
	"os"
	"server_handler/cmd/queue"

	"github.com/spf13/cobra"
)

var shortDescription = "Up minecraft server from a docker compose file."

var longDescription = `Up minecraft server from a docker compose file.
The compose file must be specified in a YAML file in ~/.config/admine/server.yaml.

server.yaml content
composeDirectory: "/compose/absolute/path.yaml"
consumerChannel: "channel-who-receives-commands"
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
