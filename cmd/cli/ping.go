package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ddritzenhoff/dindin/internal/http/rpc/pb"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(cmdPing)
}

var cmdPing = &cobra.Command{
	Use:   "ping",
	Short: "ping the server to check to see if it's working",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		conn, err := generateGRPCClientConnection()
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()
		slackClient := pb.NewSlackActionsClient(conn)
		msg, err := slackClient.Ping(context.Background(), &pb.PingMessage{Message: "sent from the client"})
		fmt.Printf("Got a response from the server: %s", msg.GetMessage())
	},
}
