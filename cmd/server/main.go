package main

import (
	"fmt"
	"github.com/ddritzenhoff/dindin/internal/configs"
	"github.com/ddritzenhoff/dindin/internal/cooking"
	"github.com/ddritzenhoff/dindin/internal/http/rest"
	"github.com/ddritzenhoff/dindin/internal/http/rpc"
	"github.com/ddritzenhoff/dindin/internal/person"
	"github.com/slack-go/slack"
	"google.golang.org/grpc"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"net"
)

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

	h, err := rest.NewHTTPService(httpCfg, ps)
	if err != nil {
		log.Fatal(err)
	}
	go h.Start()

	// create a listener on TCP port 7777
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 7777))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	// create a http instance
	s := rpc.NewServer(ces)
	// create a gRPC http object
	grpcServer := grpc.NewServer()
	// attach the Ping service to the http
	rpc.RegisterSlackActionsServer(grpcServer, &s)
	// start the http
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
