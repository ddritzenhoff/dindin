package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ddritzenhoff/dindin/internal/http/rpc"
	"github.com/spf13/cobra"
)

var devChannel bool
var slackUID string

func init() {
	cmdSlackMessage.Flags().BoolVarP(&devChannel, "dev", "d", false, "use flag to execute command in the dev slack channel")
	cmdSlackMessage.Flags().StringVarP(&slackUID, "slackUID", "uid", "", "The person responsible for cooking (required)")
	cmdSlackMessage.MarkFlagRequired("slackUID")
	rootCmd.AddCommand(cmdSlackMessage)
}

var cmdSlackMessage = &cobra.Command{
	Use:   "slack",
	Short: "send a 'is eating' slack message",
	Run: func(cmd *cobra.Command, args []string) {
		conn, err := generateGRPCClientConnection()
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()
		slackClient := rpc.NewSlackActionsClient(conn)
		msg, err := slackClient.EatingTomorrow(context.Background(), &rpc.EatingTomorrowRequest{SlackUID: slackUID})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Got a response from the server: %s", msg.String())
	},
}
