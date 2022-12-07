package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ddritzenhoff/dinny/http/rpc/pb"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(cmdEatingTomorrow)
}

// cmdEatingTomorrow sends a 'like to eat tomorrow' slack message if a user is set to cook tomorrow.
var cmdEatingTomorrow = &cobra.Command{
	Use:   "eating_tomorrow",
	Short: "send a 'like to eat tomorrow' slack message if a user is set to cook tomorrow",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		conn, err := generateGRPCClientConnection()
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()
		slackClient := pb.NewSlackActionsClient(conn)
		_, err = slackClient.EatingTomorrow(context.Background(), &pb.EmptyMessage{})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("success")
	},
}
