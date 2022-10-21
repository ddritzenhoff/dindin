package main

import (
	"context"
	"log"
	"time"

	"github.com/ddritzenhoff/dindin/internal/http/rpc/pb"
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
	assignCooks.Flags().StringVarP(&mondaySlackUID, "monday", "m", "", "set cook for monday <slackUID>")
	assignCooks.Flags().StringVarP(&tuesdaySlackUID, "tuesday", "t", "", "set cook for tuesday <slackUID>")
	assignCooks.Flags().StringVarP(&wednesdaySlackUID, "wednesday", "w", "", "set cook for wednesday <slackUID>")
	assignCooks.Flags().StringVarP(&thursdaySlackUID, "thursday", "h", "", "set cook for thursday <slackUID>")
	assignCooks.Flags().StringVarP(&fridaySlackUID, "friday", "f", "", "set cook for friday <slackUID>")
	assignCooks.Flags().StringVarP(&saturdaySlackUID, "saturday", "s", "", "set cook for saturday <slackUID>")
	assignCooks.Flags().StringVarP(&sundaySlackUID, "sunday", "u", "", "set cook for sunday <slackUID>")
}

func getDayDifference(now time.Weekday, then time.Weekday) int {
	return (int(then) - int(now) + 7) % 7
}

func buildCookingDay(now time.Time, then time.Weekday, slackUID string) *pb.CookingDay {
	year, month, day := now.AddDate(0, 0, getDayDifference(now.Weekday(), then)).Date()
	return &pb.CookingDay{
		Day:      int32(day),
		Month:    int32(month),
		Year:     int32(year),
		SlackUID: slackUID,
	}
}

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
		var cookingDays []*pb.CookingDay

		if sundaySlackUID != "" {
			cookingDays = append(cookingDays, buildCookingDay(now, time.Sunday, sundaySlackUID))
		}
		if mondaySlackUID != "" {
			cookingDays = append(cookingDays, buildCookingDay(now, time.Monday, mondaySlackUID))
		}
		if tuesdaySlackUID != "" {
			cookingDays = append(cookingDays, buildCookingDay(now, time.Tuesday, tuesdaySlackUID))
		}
		if wednesdaySlackUID != "" {
			cookingDays = append(cookingDays, buildCookingDay(now, time.Wednesday, wednesdaySlackUID))
		}
		if thursdaySlackUID != "" {
			cookingDays = append(cookingDays, buildCookingDay(now, time.Thursday, thursdaySlackUID))
		}
		if fridaySlackUID != "" {
			cookingDays = append(cookingDays, buildCookingDay(now, time.Friday, fridaySlackUID))
		}
		if saturdaySlackUID != "" {
			cookingDays = append(cookingDays, buildCookingDay(now, time.Saturday, saturdaySlackUID))
		}

		if len(cookingDays) == 0 {
			return
		}

		slackActionsClient := pb.NewSlackActionsClient(conn)
		_, err = slackActionsClient.AssignCooks(context.Background(), &pb.AssignCooksRequest{CookingDays: cookingDays})
		if err != nil {
			log.Fatal(err)
		}
	},
}
