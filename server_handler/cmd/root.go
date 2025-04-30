package cmd

import (
	"fmt"
	"os"
	"server_handler/cmd/queue"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "hugo",
	Short: "short hugo",
	Long:  "long hugo",
	Run: func(cmd *cobra.Command, args []string) {
		print("rodando queue")

		queue.RunListenQueue()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
