package main

import (
	"github.com/ddritzenhoff/dindin/internal/configs"
	"github.com/ddritzenhoff/dindin/internal/cooking"
	"github.com/ddritzenhoff/dindin/internal/person"
	"github.com/ddritzenhoff/dindin/internal/server/http"
	"github.com/slack-go/slack"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
)

// TODO (ddritzenhoff) add logging
// TODO (ddritzenhoff) add better error handling

func main() {
	cfg, err := configs.NewConfigService()
	if err != nil {
		log.Fatal(err)
	}

	dbName, err := cfg.DBName()
	if err != nil {
		log.Fatal(err)
	}

	db, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	slackConfig, err := cfg.SlackConfig()
	if err != nil {
		log.Fatal(err)
	}
	slackClient := slack.New(slackConfig.BotUserToken)

	ces, err := cooking.NewEventService(db, slackConfig.BotTestChannel, slackClient)
	if err != nil {
		log.Fatal(err)
	}

	ps, err := person.NewService(db, ces)
	if err != nil {
		log.Fatal(err)
	}

	httpCfg, err := cfg.HTTP()
	if err != nil {
		log.Fatal(err)
	}

	h, err := http.NewHTTPService(httpCfg, ps)
	if err != nil {
		log.Fatal(err)
	}
	h.Start()
}
