package main

import (
	"context"
	_ "embed"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/pelletier/go-toml/v2"

	rest "github.com/ddritzenhoff/dinny/http"
	"github.com/ddritzenhoff/dinny/slack"
	"github.com/ddritzenhoff/dinny/sqlite"
	"github.com/ddritzenhoff/dinny/sqlite/gen"
	_ "github.com/mattn/go-sqlite3"
)

func main() {

	// Setup signal handlers.
	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() { <-c; cancel() }()

	// Parse command line flag and load configuration.
	config, err := ParseFlag(context.Background(), os.Args[1:])
	if err != nil {
		if err == flag.ErrHelp {
			os.Exit(1)
		} else {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}

	// Execute program.
	if err := Run(config); err != nil {
		log.Fatal(err)
	}

	// Wait for CTRL-C.
	<-ctx.Done()

	// Clean up program.
	// TODO (ddritzenhoff)
}

// ParseFlag parses the config flag and loads the config.
func ParseFlag(context context.Context, args []string) (*Config, error) {
	fs := flag.NewFlagSet("dinnyd", flag.ContinueOnError)
	var configPath string
	fs.StringVar(&configPath, "config", DefaultConfigPath, "config path")
	if err := fs.Parse(args); err != nil {
		return nil, fmt.Errorf("ParseFlag fs.Parse: %w", err)
	}

	// The expand() function is here to automatically expand "~" to the user's
	// home directory.
	configPath, err := expand(configPath)
	if err != nil {
		return nil, fmt.Errorf("ParseFlag expand: %w", err)
	}

	// Read the TOML formatted configuration file.
	config, err := ReadConfigFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("ParseFlag ReadConfigFile: %w", err)
	}
	return &config, nil
}

// run initializes the member, meal, and Slack services and starts the REST and GRPC servers.
func Run(config *Config) error {
	logger := log.New(os.Stdout, "DEBUG: ", log.LstdFlags)

	DSNPath, err := expandDSN(config.DB.DSN)
	if err != nil {
		return fmt.Errorf("Run expandDSN: %w", err)
	}

	db, err := sqlite.Open(DSNPath)
	if err != nil {
		return fmt.Errorf("Run sqlite.Open: %w", err)
	}

	queries := gen.New(db)

	mealService := sqlite.NewMealService(queries, db)

	memberService := sqlite.NewMemberService(queries, db)

	slackConfig := slack.Config{
		Channel:       config.Slack.ChannelID,
		BotSigningKey: config.Slack.BotSigningKey,
	}

	slackService, err := slack.NewService(&slackConfig, mealService, memberService)
	if err != nil {
		return fmt.Errorf("Run slack.NewService: %w", err)
	}

	restServer := rest.NewServer(logger, config.HTTP.URL, memberService, mealService, slackService)
	restServer.Open()

	return nil
}

const (
	// DefaultConfigPath is the the default path to the application configuration.
	DefaultConfigPath = "~/dinnyd.toml"

	// DefaultDSN is the default datasource name.
	DefaultDSN = "~/.dinnyd/db"
)

// Config represents the CLI configuration file.
type Config struct {
	DB struct {
		DSN string `toml:"dsn"`
	} `toml:"db"`

	HTTP struct {
		URL string `toml:"url"`
	} `toml:"http"`

	Slack struct {
		BotSigningKey string `toml:"botSigningKey"`
		AppID         string `toml:"appID"`
		ClientID      string `toml:"clientID"`
		ClientSecret  string `toml:"clientSecret"`
		SigningSecret string `toml:"signingSecret"`
		ChannelID     string `toml:"channelID"`
	} `toml:"slack"`
}

// DefaultConfig returns a new instance of Config with defaults set.
func DefaultConfig() Config {
	var config Config
	config.DB.DSN = DefaultDSN
	return config
}

// ReadConfigFile unmarshals config from
func ReadConfigFile(filename string) (Config, error) {
	var config Config
	if buf, err := os.ReadFile(filename); os.IsNotExist(err) {
		return config, fmt.Errorf("config file with path %s not found: %w", filename, err)
	} else if err != nil {
		return config, fmt.Errorf("NewConfigService os.ReadFile: %w", err)
	} else if toml.Unmarshal(buf, &config); err != nil {
		return config, fmt.Errorf("NewConfigService toml.Unmarshal: %w", err)
	}
	return config, nil
}

// expand returns path using tilde expansion. This means that a file path that
// begins with the "~" will be expanded to prefix the user's home directory.
func expand(path string) (string, error) {
	// Ignore if path has no leading tilde.
	if path != "~" && !strings.HasPrefix(path, "~"+string(os.PathSeparator)) {
		return path, nil
	}

	// Fetch the current user to determine the home path.
	u, err := user.Current()
	if err != nil {
		return path, fmt.Errorf("expand user.Current: %w", err)
	} else if u.HomeDir == "" {
		return path, fmt.Errorf("expand u.HomeDir: home directory unset")
	}

	if path == "~" {
		return u.HomeDir, nil
	}
	return filepath.Join(u.HomeDir, strings.TrimPrefix(path, "~"+string(os.PathSeparator))), nil
}

// expandDSN expands a datasource name. Ignores in-memory databases.
func expandDSN(dsn string) (string, error) {
	if dsn == ":memory:" {
		return dsn, nil
	}
	return expand(dsn)
}
