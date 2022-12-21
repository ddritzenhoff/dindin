package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	"github.com/ddritzenhoff/dinny/http/grpc/pb"
)

var mondaySlackUID string
var tuesdaySlackUID string
var wednesdaySlackUID string
var thursdaySlackUID string
var fridaySlackUID string
var saturdaySlackUID string
var sundaySlackUID string

// AssignCooksCommand is a command to assign cooks.
type AssignCooksCommand struct {
	ConfigPath string
}

// Run executes the assign_cooks command.
func (c *AssignCooksCommand) Run(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.StringVar(&mondaySlackUID, "monday", "", "set cook for monday <slackUID>")
	fs.StringVar(&tuesdaySlackUID, "tuesday", "", "set cook for tuesday <slackUID>")
	fs.StringVar(&wednesdaySlackUID, "wednesday", "", "set cook for wednesday <slackUID>")
	fs.StringVar(&thursdaySlackUID, "thursday", "", "set cook for thursday <slackUID>")
	fs.StringVar(&fridaySlackUID, "friday", "", "set cook for friday <slackUID>")
	fs.StringVar(&saturdaySlackUID, "saturday", "", "set cook for saturday <slackUID>")
	fs.StringVar(&sundaySlackUID, "sunday", "", "set cook for sunday <slackUID>")
	attachConfigFlags(fs, &c.ConfigPath)
	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("Run fs.Parse: %w", err)
	}

	// Load the configuration.
	config, err := ReadConfigFile(c.ConfigPath)
	if err != nil {
		return fmt.Errorf("Run ReadConfigFile: %w", err)
	}

	conn, err := generateGRPCClientConnectionWithAddress(config.URL)
	if err != nil {
		return fmt.Errorf("Run generateGRPCClientConnection: %w", err)
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
		return fmt.Errorf("Run: didn't specify any days, so nothing happened..")
	}

	slackActionsClient := pb.NewSlackActionsClient(conn)
	_, err = slackActionsClient.AssignCooks(context.Background(), &pb.AssignCooksRequest{Assignments: assignments})
	if err != nil {
		return fmt.Errorf("Run slackActionsClient.AssignCooks: %w", err)
	}

	return nil
}

// getDayDifference returns the number of days between now (today) and then (some other weekday).
// This function will return some value between 0 and 6. 0 --> now=Monday, then=Monday 6 --> now=Wednesday, then=Tuesday.
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

// usage prints usage information for assign_cooks to STDOUT.
func (c *AssignCooksCommand) usage() {
	fmt.Println(`
Assign cooks for the next week starting with the current day.
If today were Monday, and a cook were to be assigned with the -monday flag, that cook would be set to cook today.

Usage:

		dinny assign_cooks -monday <slackUID> -tuesday <slackUID> -wednesday <slackUID> -thursday <slackUID> -friday <slackUID> -saturday <slackUID> -sunday <slackUID>

Arguments:

		-monday <slackUID>
			Set the cook for the upcoming Monday
		-tuesday <slackUID>
			Set the cook for the upcoming Tuesday
		-wednesday <slackUID>
			Set the cook for the upcoming Wednesday
		-thursday <slackUID>
			Set the cook for the upcoming Thursday
		-friday <slackUID>
			Set the cook for the upcoming Friday
		-saturday <slackUID>
			Set the cook for the upcoming Saturday
		-sunday <slackUID>
			Set the cook for the upcoming Sunday
	`[1:])
}
