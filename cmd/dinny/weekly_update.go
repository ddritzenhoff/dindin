package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
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

	url := fmt.Sprintf("%s/cmd/weekly-update", config.URL)
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("Run http.Get: %w", err)
	}
	defer resp.Body.Close()
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
