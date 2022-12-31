package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
)

// EatingTomorrowCommand is a command to send a 'who's eating tomorrow' message into Slack.
type EatingTomorrowCommand struct {
	ConfigPath string
}

// Run executes the eating_tomorrow command.
func (c *EatingTomorrowCommand) Run(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	attachConfigFlags(fs, &c.ConfigPath)
	fs.Usage = c.usage
	err := fs.Parse(args)
	if err != nil {
		if err == flag.ErrHelp {
			os.Exit(1)
		} else {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}

	// Load the configuration.
	config, err := ReadConfigFile(c.ConfigPath)
	if err != nil {
		return fmt.Errorf("Run ReadConfigFile: %w", err)
	}

	url := fmt.Sprintf("%s/cmd/eating-tomorrow", config.URL)
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("Run http.Get: %w", err)
	}
	defer resp.Body.Close()
	fmt.Println("success")
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
