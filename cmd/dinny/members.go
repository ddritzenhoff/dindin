package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"

	"github.com/ddritzenhoff/dinny/http/grpc/pb"
)

// MembersCommand is a command to list the current members of dinner rotation.
type MembersCommand struct {
	ConfigPath string
}

// Run executes the members command.
func (c *MembersCommand) Run(ctx context.Context, args []string) error {
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
	stream, err := slackClient.GetMembers(context.Background(), &pb.EmptyMessage{})
	if err != nil {
		return fmt.Errorf("Run client.GetMembers: %w", err)
	}
	for {
		memberInfo, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("client.GetMembers failed: %v", err)
		}
		log.Printf("Full Name: %s\nSlackUID: %s\n\n", memberInfo.GetFullName(), memberInfo.GetSlack_UID())
	}
	return nil
}

// usage prints usage information for members to STDOUT.
func (c *MembersCommand) usage() {
	fmt.Println(`
List the current members of dinner rotation.

Usage:

		dinny members
	`[1:])
}
