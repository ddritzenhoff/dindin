package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ddritzenhoff/dindin/internal/http/rpc/pb"
	"github.com/spf13/cobra"
)

var daysWanted *int64

func init() {
	rootCmd.AddCommand(upcomingCooks)
	daysWanted = upcomingCooks.Flags().Int64P("days", "d", 1, "get the number cooking events wanted")
}

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
		for _, c := range r.Cooks {
			fmt.Printf("Name: %s\n\tSlackUID: %s,\n\tCooking Time: %d,\n\tDesc: %s,\n\tMessageID: %s,", c.FirstName+" "+c.LastName, c.ChefSlack_UID, c.CookingTime, c.MealDescription, c.SlackMessage_ID)
			fmt.Printf("%#v\n", c)
		}
	},
}
