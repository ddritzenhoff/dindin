package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/ddritzenhoff/dinny/http/grpc/pb"
)

var daysWanted int64

// UpcomingCooksCommand is a command to list the upcoming cooks for the next week.
type UpcomingCooksCommand struct {
	ConfigPath string
}

// Run executes the upcoming_cooks command.
func (c *UpcomingCooksCommand) Run(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.Int64Var(&daysWanted, "days", 7, "list the next days' cooks")
	attachConfigFlags(fs, &c.ConfigPath)
	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("Run fs.Parse: %w", err)
	}

	if daysWanted < 1 {
		return fmt.Errorf("Run: -days flag must a value above 0")
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
	slackClient := pb.NewSlackActionsClient(conn)
	r, err := slackClient.UpcomingCooks(context.Background(), &pb.UpcomingCooksRequest{DaysWanted: daysWanted})
	if err != nil {
		return fmt.Errorf("Run slackClient.UpcomingCooks: %w", err)
	}
	for _, m := range r.Meals {
		fmt.Printf("Name: %s\n\tSlackUID: %s,\n\tCooking Time: %s,\n\tDesc: %s,\n\tMessageID: %s,", m.FullName, m.CookSlack_UID, fmt.Sprintf("%d/%d/%d\n\t", m.Date.Month, m.Date.Day, m.Date.Year), m.Description, m.SlackMessage_ID)
	}
	return nil
}

// usage prints usage information for upcoming_cooks to STDOUT.
func (c *UpcomingCooksCommand) usage() {
	fmt.Println(`
List the upcoming cooks.

Usage:

		dinny upcoming_cooks -days

Arguments:

		-days <int64>
`[1:])
}
