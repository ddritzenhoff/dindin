package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	var testChannel bool
	var cmdSlackMessage = &cobra.Command{
		Use:   "slack",
		Short: "send a 'is eating' slack message",
		Run:   cmdSlackMessage,
	}

	cmdSlackMessage.Flags().BoolVarP(&testChannel, "test", "t", false, "use flag to execute command in the test slack channel")

	var rootCmd = &cobra.Command{
		Use: "app",
	}

	rootCmd.AddCommand(cmdSlackMessage)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
