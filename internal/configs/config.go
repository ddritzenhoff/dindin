package configs

import (
	"fmt"
	"log"

	"github.com/slack-go/slack"
	"github.com/spf13/viper"
)

// Configs handles all dependencies required for handling configurations
type Configs struct{}

func (cfg *Configs) DBName() (string, error) {
	if viper.IsSet("database.name") {
		return viper.GetString("database.name"), nil
	}
	return "", fmt.Errorf("couldn't find the config's database name attribute")
}

type HTTP struct {
	Host string
	Port string
}

func (cfg *Configs) REST() (*HTTP, error) {
	if !viper.IsSet("http.rest.port") {
		return nil, fmt.Errorf("couldn't find the config's rest port attribute")
	}
	return &HTTP{
		Host: viper.GetString("http.rest.host"),
		Port: viper.GetString("http.rest.port"),
	}, nil
}

func (cfg *Configs) GRPC() (*HTTP, error) {
	if !viper.IsSet("http.grpc.host") {
		return nil, fmt.Errorf("couldn't find the config's grpc host attribute")
	}
	if !viper.IsSet("http.grpc.port") {
		return nil, fmt.Errorf("couldn't find the config's grpc port attribute")
	}
	return &HTTP{
		Host: viper.GetString("http.grpc.host"),
		Port: viper.GetString("http.grpc.port"),
	}, nil
}

type SlackConfig struct {
	Channel string
	Client  *slack.Client
}

func (cfg *Configs) SlackConfig() (*SlackConfig, error) {
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
	client := slack.New(viper.GetString("slack.botSigningKey"))
	if client == nil {
		log.Fatal("slack client reference is nil")
	}
	return &SlackConfig{
		Channel: channel,
		Client:  client,
	}, nil
}

func NewConfigService() (*Configs, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	// TODO: (ddritzenhoff) find a better way of loading the config. Maybe
	// 	you can pass it in instead?
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/programming/dindin/internal/configs/.")
	viper.AddConfigPath("$HOME/config/.")
	viper.AddConfigPath("$HOME/.")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("fatal error config file: %w", err)
	}
	return &Configs{}, nil
}
