package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/ddritzenhoff/dinny/http/grpc/pb"
)

// PingCommand is a command to ping the dinny service to check health.
type PingCommand struct {
	ConfigPath string
}

// Run executes the ping command.
func (c *PingCommand) Run(ctx context.Context, args []string) error {
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
	msg, err := slackClient.Ping(context.Background(), &pb.EmptyMessage{})
	if err != nil {
		return fmt.Errorf("Run slackClient.Ping: %w", err)
	} else {
		fmt.Printf("%s\n", msg.GetMessage())
	}
	return nil
}

// usage prints usage information for ping to STDOUT.
func (c *PingCommand) usage() {
	fmt.Println(`
Ping the dinny service to check health.

Usage:

		dinny ping
`[1:])
}
