package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/ddritzenhoff/dinny/http/grpc/pb"
)

// WeeklyUpdateCommand is a command to send a message into slack with each member's meals eaten to meals cooked ratio.
type WeeklyUpdateCommand struct {
	ConfigPath string
}

// Run executes the weekly_update command.
func (c *WeeklyUpdateCommand) Run(ctx context.Context, args []string) error {
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
		return fmt.Errorf("Run generateGRPCClientConnection: %w", err)
	}
	defer conn.Close()
	slackClient := pb.NewSlackActionsClient(conn)
	_, err = slackClient.WeeklyUpdate(context.Background(), &pb.EmptyMessage{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("success")
	return nil
}

// usage prints usage information for weekly_update to STDOUT.
func (c *WeeklyUpdateCommand) usage() {
	fmt.Println(`
Send a message into slack with each member's meals eaten to meals cooked ratio.

Usage:

		dinny weekly_update
`[1:])
}
