package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ddritzenhoff/dindin/http/rpc/pb"
	"github.com/spf13/cobra"
)

var mondaySlackUID string
var tuesdaySlackUID string
var wednesdaySlackUID string
var thursdaySlackUID string
var fridaySlackUID string
var saturdaySlackUID string
var sundaySlackUID string

func init() {
	rootCmd.AddCommand(assignCooks)
	assignCooks.Flags().StringVar(&mondaySlackUID, "monday", "", "set cook for monday <slackUID>")
	assignCooks.Flags().StringVar(&tuesdaySlackUID, "tuesday", "", "set cook for tuesday <slackUID>")
	assignCooks.Flags().StringVar(&wednesdaySlackUID, "wednesday", "", "set cook for wednesday <slackUID>")
	assignCooks.Flags().StringVar(&thursdaySlackUID, "thursday", "", "set cook for thursday <slackUID>")
	assignCooks.Flags().StringVar(&fridaySlackUID, "friday", "", "set cook for friday <slackUID>")
	assignCooks.Flags().StringVar(&saturdaySlackUID, "saturday", "", "set cook for saturday <slackUID>")
	assignCooks.Flags().StringVar(&sundaySlackUID, "sunday", "", "set cook for sunday <slackUID>")
}

// getDayDifference returns the number of days between now (today) and then (some other weekday).
// This function will return some value between 0 and 6. 0 --> now=Monday, then=Monday 6 --> now=Wednesday, then=Tuesday
func getDayDifference(now time.Weekday, then time.Weekday) int {
	return (int(then) - int(now) + 7) % 7
}

// buildCookAssignment creates a pb.CookAssignment.
func buildCookAssignment(now time.Time, then time.Weekday, slackUID string) *pb.CookAssignment {
	year, month, day := now.AddDate(0, 0, getDayDifference(now.Weekday(), then)).Date()
	return &pb.CookAssignment{
		Date: &pb.Date{
			Year:  int64(year),
			Month: int64(month),
			Day:   int64(day),
		},
		Slack_UID: slackUID,
	}
}

// assignCooks assigns cooks for the next week.
var assignCooks = &cobra.Command{
	Use:   "assign_cooks",
	Short: "assign the cooks for the next week of dinner rotation",
	Args:  cobra.NoArgs,
	Long:  "Will automatically set the upcoming day with the cook. That is, if it's on Tuesday that I issue this command, and I set a value for the Wednesday flag, the next day would be assigned a cook.",
	Run: func(cmd *cobra.Command, args []string) {
		conn, err := generateGRPCClientConnection()
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()

		now := time.Now()
		var assignments []*pb.CookAssignment

		if sundaySlackUID != "" {
			assignments = append(assignments, buildCookAssignment(now, time.Sunday, sundaySlackUID))
		}
		if mondaySlackUID != "" {
			assignments = append(assignments, buildCookAssignment(now, time.Monday, mondaySlackUID))
		}
		if tuesdaySlackUID != "" {
			assignments = append(assignments, buildCookAssignment(now, time.Tuesday, tuesdaySlackUID))
		}
		if wednesdaySlackUID != "" {
			assignments = append(assignments, buildCookAssignment(now, time.Wednesday, wednesdaySlackUID))
		}
		if thursdaySlackUID != "" {
			assignments = append(assignments, buildCookAssignment(now, time.Thursday, thursdaySlackUID))
		}
		if fridaySlackUID != "" {
			assignments = append(assignments, buildCookAssignment(now, time.Friday, fridaySlackUID))
		}
		if saturdaySlackUID != "" {
			assignments = append(assignments, buildCookAssignment(now, time.Saturday, saturdaySlackUID))
		}

		if len(assignments) == 0 {
			fmt.Println("didn't specify any days, so nothing happened..")
			return
		}

		slackActionsClient := pb.NewSlackActionsClient(conn)
		_, err = slackActionsClient.AssignCooks(context.Background(), &pb.AssignCooksRequest{Assignments: assignments})
		if err != nil {
			log.Fatal(err)
		}
	},
}
