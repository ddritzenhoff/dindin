package configs

import (
	"log"
	"os"
	"strings"
)

const botTestChannel = "C028HTSA42K"
const dinnerRotationChannel = "CTEKWPTD1"

type SlackConfig struct {
	BotUserToken          string
	BotTestChannel        string
	DinnerRotationChannel string
}

func NewSlackConfig() (*SlackConfig, error) {
	botUserToken, ok := os.LookupEnv("BOT_SIGNING_KEY")
	if !ok {
		log.Fatal("BOT_SIGNING_KEY is not set")
	}

	s := &SlackConfig{
		BotUserToken:          strings.TrimSpace(botUserToken),
		BotTestChannel:        botTestChannel,
		DinnerRotationChannel: dinnerRotationChannel,
	}

	return s, nil
}
