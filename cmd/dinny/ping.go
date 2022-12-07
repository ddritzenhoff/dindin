package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ddritzenhoff/dinny/http/rpc/pb"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(cmdPing)
}

// cmdPing pings the dinner rotation bot server to see if it's active.
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
		msg, err := slackClient.Ping(context.Background(), &pb.EmptyMessage{})
		if err != nil {
			fmt.Printf("error from server: %s", err.Error())
		} else {
			fmt.Printf("%s\n", msg.GetMessage())
		}
	},
}
