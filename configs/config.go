package configs

import (
	"fmt"
	"log"

	"github.com/ddritzenhoff/dinny/http/rest"
	"github.com/ddritzenhoff/dinny/http/rpc"
	"github.com/ddritzenhoff/dinny/slack"
	"github.com/spf13/viper"
)

// Configs handles all dependencies required for handling configurations
type Configs struct{}

// DBName returns the name of the database.
func (cfg *Configs) DBName() (string, error) {
	if viper.IsSet("database.name") {
		return viper.GetString("database.name"), nil
	}
	return "", fmt.Errorf("couldn't find the config's database name attribute")
}

// REST represents the config values for the REST server.
func (cfg *Configs) REST() (*rest.Config, error) {
	if !viper.IsSet("http.rest.host") {
		return nil, fmt.Errorf("couldn't find the config's rest host attribute")
	}
	if !viper.IsSet("http.rest.port") {
		return nil, fmt.Errorf("couldn't find the config's rest port attribute")
	}
	return &rest.Config{
		Host: viper.GetString("http.rest.host"),
		Port: viper.GetString("http.rest.port"),
	}, nil
}

// GRPC represents the config values for the GRPC server.
func (cfg *Configs) GRPC() (*rpc.Config, error) {
	if !viper.IsSet("http.grpc.host") {
		return nil, fmt.Errorf("couldn't find the config's grpc host attribute")
	}
	if !viper.IsSet("http.grpc.port") {
		return nil, fmt.Errorf("couldn't find the config's grpc port attribute")
	}
	return &rpc.Config{
		Host: viper.GetString("http.grpc.host"),
		Port: viper.GetString("http.grpc.port"),
	}, nil
}

// Slack represents the configs for slack.
func (cfg *Configs) Slack() (*slack.Config, error) {
	if !viper.IsSet("slack.botSigningKey") {
		return nil, fmt.Errorf("couldn't find the config's slack botSigningKey")
	}
	if !viper.IsSet("slack.dev.channelID") {
		return nil, fmt.Errorf("couldn't find the config's slack dev channelID")
	}
	if !viper.IsSet("slack.prod.channelID") {
		return nil, fmt.Errorf("couldn't find the config's slack prod channelID")
	}
	var channel string
	if viper.GetBool("slack.isProd") {
		channel = viper.GetString("slack.prod.channelID")
	} else {
		channel = viper.GetString("slack.dev.channelID")
	}
	return &slack.Config{
		Channel:       channel,
		BotSigningKey: viper.GetString("slack.botSigningKey"),
	}, nil
}

// NewConfigService reads in the config file.
func NewConfigService() (*Configs, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	// TODO: (ddritzenhoff) find a better way of loading the config. Maybe
	// 	you can pass it in instead?
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/programming/dinny/configs/.")
	viper.AddConfigPath("$HOME/config/.")
	viper.AddConfigPath("$HOME/.")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("fatal error config file: %w", err)
	}
	return &Configs{}, nil
}
