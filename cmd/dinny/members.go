package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"
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

	url := fmt.Sprintf("%s/cmd/members", config.URL)
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("Run http.Get: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Run io.ReadAll: %w", err)
	}
	b, err := prettyPrint(body)
	if err != nil {
		return fmt.Errorf("Run prettyPrint: %w", err)
	}
	fmt.Println(string(b))
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
