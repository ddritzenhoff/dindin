package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/ddritzenhoff/dinny/http/grpc/pb"
)

// EatingTomorrowCommand is a command to send a 'who's eating tomorrow' message into Slack.
type EatingTomorrowCommand struct {
	ConfigPath string
}

// Run executes the eating_tomorrow command.
func (c *EatingTomorrowCommand) Run(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet("", flag.ContinueOnError)
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
		return fmt.Errorf("Run generateGRPCClient: %w", err)
	}
	defer conn.Close()
	slackClient := pb.NewSlackActionsClient(conn)
	_, err = slackClient.EatingTomorrow(context.Background(), &pb.EmptyMessage{})
	if err != nil {
		return fmt.Errorf("Run slackClient.EatingTomorrow: %w", err)
	}
	return nil
}

// usage prints usage information for eating_tomorrow to STDOUT.
func (c *EatingTomorrowCommand) usage() {
	fmt.Println(`
Send a 'like to eat tomorrow' slack message.

Usage:
        dinny eating_tomorrow
`[1:])
}
