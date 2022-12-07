package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ddritzenhoff/dinny/http/rpc/pb"
	"github.com/spf13/cobra"
)

var daysWanted *int64

func init() {
	rootCmd.AddCommand(upcomingCooks)
	daysWanted = upcomingCooks.Flags().Int64P("days", "d", 1, "get the number cooking events wanted")
}

// upcomingCooks gets the cooks who are in charge of the meals over the next few days.
var upcomingCooks = &cobra.Command{
	Use:   "upcoming_cooks",
	Short: "get the upcoming cooks",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		conn, err := generateGRPCClientConnection()
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()
		slackClient := pb.NewSlackActionsClient(conn)
		r, err := slackClient.UpcomingCooks(context.Background(), &pb.UpcomingCooksRequest{DaysWanted: *daysWanted})
		if err != nil {
			log.Fatal(err)
		}
		for _, m := range r.Meals {
			fmt.Printf("Name: %s\n\tSlackUID: %s,\n\tCooking Time: %s,\n\tDesc: %s,\n\tMessageID: %s,", m.FullName, m.CookSlack_UID, fmt.Sprintf("%d/%d/%d\n\t", m.Date.Month, m.Date.Day, m.Date.Year), m.Description, m.SlackMessage_ID)
		}
	},
}
