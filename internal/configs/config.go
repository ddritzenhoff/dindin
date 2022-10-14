package configs

import (
	"fmt"

	"github.com/ddritzenhoff/dindin/internal/http/rest"
	"github.com/spf13/viper"
)

// Configs handles all dependencies required for handling configurations
type Configs struct{}

func (cfg *Configs) DBName() (string, error) {
	if viper.IsSet("database.name") {
		return viper.GetString("database.name"), nil
	}
	return "", fmt.Errorf("Couldn't find the config's database name attribute")
}

func (cfg *Configs) REST() (*rest.Config, error) {
	if !viper.IsSet("http.rest.port") {
		return nil, fmt.Errorf("Couldn't find the config's rest port attribute")
	}
	return &rest.Config{
		Host: viper.GetString("http.rest.host"),
		Port: viper.GetString("http.rest.port"),
	}, nil
}

func (cfg *Configs) GRPC() (*rest.Config, error) {
	if !viper.IsSet("http.grpc.host") {
		return nil, fmt.Errorf("Couldn't find the config's grpc host attribute")
	}
	if !viper.IsSet("http.grpc.port") {
		return nil, fmt.Errorf("Couldn't find the config's grpc port attribute")
	}
	return &rest.Config{
		Host: viper.GetString("http.grpc.host"),
		Port: viper.GetString("http.grpc.port"),
	}, nil
}

func (cfg *Configs) SlackConfig() (*slackConfig, error) {
	return newSlackConfig()
}

func NewConfigService() (*Configs, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	// TODO: (ddritzenhoff) find a better way of loading the config. Maybe
	// 	you can pass it in instead?
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/programming/dindin/internal/configs/.")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	return &Configs{}, nil
}
