package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/pelletier/go-toml/v2"
)

const (
	DefaultURL        = "localhost:7777"
	DefaultConfigPath = "~/dinny.toml"
)

func main() {
	// Execute program.
	if err := Run(context.Background(), os.Args[1:]); err == flag.ErrHelp {
		os.Exit(1)
	} else if err != nil {
		fmt.Printf("main: %s", err.Error())
	}
}

// Run executes the main program.
func Run(ctx context.Context, args []string) error {
	var cmd string
	if len(args) > 0 {
		cmd, args = args[0], args[1:]
	}

	switch cmd {
	case "", "help", "-h", "--help":
		usage()
		return flag.ErrHelp
	case "assign_cooks":
		return (&AssignCooksCommand{}).Run(ctx, args)
	case "eating_tomorrow":
		return (&EatingTomorrowCommand{}).Run(ctx, args)
	case "members":
		return (&MembersCommand{}).Run(ctx, args)
	case "ping":
		return (&PingCommand{}).Run(ctx, args)
	case "upcoming_cooks":
		return (&UpcomingCooksCommand{}).Run(ctx, args)
	case "weekly_update":
		return (&WeeklyUpdateCommand{}).Run(ctx, args)
	default:
		return fmt.Errorf("dinny %s: unknown command", cmd)
	}
}

// usage prints the top-level CLI usage message.
func usage() {
	fmt.Println(`
Command line utility for interacting with the dinny service.

Usage:
		dinny <command> [arguments]

The commands are:

		assign_cooks		assign cooks for the next week
		eating_tomorrow		send a 'who's eating tomorrow' message within slack
		members			list the current members of dinner rotation
		ping			ping the dinny service to check health
		upcoming_cooks		list the upcoming cooks for the next week
		weekly_update		send a message into slack with each member's meals eaten to meals cooked ratio
`[1:])
}

// Config represents a configuration file common to all subcommands.
type Config struct {
	// URL represents the base url of the server.
	URL string `toml:"url"`
}

func DefaultConfig() Config {
	return Config{
		URL: DefaultURL,
	}
}

// ReadConfigFile unmarshals config from filename. Expands path if needed.
func ReadConfigFile(filename string) (Config, error) {
	config := DefaultConfig()

	// Expand filename, if necessary. This means substituting a "~" prefix
	// with the user's home directory, if available.
	if prefix := "~" + string(os.PathSeparator); strings.HasPrefix(filename, prefix) {
		u, err := user.Current()
		if err != nil {
			return config, err
		} else if u.HomeDir == "" {
			return config, fmt.Errorf("ReadConfigFile: home directory unset")
		}
		filename = filepath.Join(u.HomeDir, strings.TrimPrefix(filename, prefix))
	}

	// Read & deserialize configuration.
	if buf, err := os.ReadFile(filename); os.IsNotExist(err) {
		return config, fmt.Errorf("ReadConfigFile: config file not found: %s", filename)
	} else if err != nil {
		return config, err
	} else if err := toml.Unmarshal(buf, &config); err != nil {
		return config, err
	}
	return config, nil
}

// attachConfigFlags adds a common "-config" flag to a flag set.
func attachConfigFlags(fs *flag.FlagSet, p *string) {
	fs.StringVar(p, "config", DefaultConfigPath, "config path")
}

// prettyPrint takes raw JSON bytes and converts into a more legible form.
func prettyPrint(b []byte) ([]byte, error) {
	var out bytes.Buffer
	err := json.Indent(&out, b, "", "  ")
	return out.Bytes(), err
}
