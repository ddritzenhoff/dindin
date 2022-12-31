package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"
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

	url := fmt.Sprintf("%s/cmd/ping", config.URL)
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("Run http.Get: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Run io.ReadAll: %w", err)
	}
	fmt.Println(string(body))
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
