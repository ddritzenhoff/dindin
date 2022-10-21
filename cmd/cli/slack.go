package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ddritzenhoff/dindin/internal/http/rpc/pb"
	"github.com/spf13/cobra"
)

var devChannel bool
var slackUID string

func init() {
	cmdSlackMessage.Flags().BoolVarP(&devChannel, "dev", "d", false, "use flag to execute command in the dev slack channel")
	cmdSlackMessage.Flags().StringVarP(&slackUID, "slackUID", "u", "", "The person responsible for cooking (required)")
	cmdSlackMessage.MarkFlagRequired("slackUID")
	rootCmd.AddCommand(cmdSlackMessage)
}

var cmdSlackMessage = &cobra.Command{
	Use:   "slack",
	Short: "send a 'is eating' slack message",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		conn, err := generateGRPCClientConnection()
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()
		slackClient := pb.NewSlackActionsClient(conn)
		msg, err := slackClient.EatingTomorrow(context.Background(), &pb.EatingTomorrowRequest{SlackUID: slackUID})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Got a response from the server: %s", msg.String())
	},
}
