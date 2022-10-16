package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ddritzenhoff/dindin/internal/http/rpc"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(cmdPing)
}

var cmdPing = &cobra.Command{
	Use:   "ping",
	Short: "ping the server to check to see if it's working",
	Run: func(cmd *cobra.Command, args []string) {
		conn, err := generateGRPCClientConnection()
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()
		slackClient := rpc.NewSlackActionsClient(conn)
		msg, err := slackClient.Ping(context.Background(), &rpc.PingMessage{Message: "sent from the client"})
		fmt.Printf("Got a response from the server: %s", msg.GetMessage())
	},
}
