package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
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

	if daysWanted < 1 {
		return fmt.Errorf("Run: -days flag must a value above 0")
	}

	// Load the configuration.
	config, err := ReadConfigFile(c.ConfigPath)
	if err != nil {
		return fmt.Errorf("Run ReadConfigFile: %w", err)
	}

	url := fmt.Sprintf("%s/cmd/upcoming-cooks", config.URL)
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
