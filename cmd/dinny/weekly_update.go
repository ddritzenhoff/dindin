package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ddritzenhoff/dinny/http/rpc/pb"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(weeklyUpdate)
}

// weeklyUpdate publishes a weekly update message into slack.
var weeklyUpdate = &cobra.Command{
	Use:   "weekly_update",
	Short: "publish a weekly update message into slack",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		conn, err := generateGRPCClientConnection()
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()
		slackClient := pb.NewSlackActionsClient(conn)
		_, err = slackClient.WeeklyUpdate(context.Background(), &pb.EmptyMessage{})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("success")
	},
}
