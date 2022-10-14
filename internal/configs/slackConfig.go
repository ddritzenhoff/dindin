package configs

import (
	"fmt"

	"github.com/spf13/viper"
)

type slackConfig struct {
	BotUserToken          string
	BotTestChannel        string
	DinnerRotationChannel string
}

func newSlackConfig() (*slackConfig, error) {
	if !viper.IsSet("slack.botSigningKey") {
		return nil, fmt.Errorf("Couldn't find the config's slack botSigningKey")
	}
	if !viper.IsSet("slack.dev.channelID") {
		return nil, fmt.Errorf("Couldn't find the config's slack dev channelID")
	}
	if !viper.IsSet("slack.prod.channelID") {
		return nil, fmt.Errorf("Couldn't find the config's slack prod channelID")
	}
	return &slackConfig{
		BotUserToken:          viper.GetString("slack.botSigningKey"),
		BotTestChannel:        viper.GetString("slack.dev.channelID"),
		DinnerRotationChannel: viper.GetString("slack.prod.channelID"),
	}, nil
}
