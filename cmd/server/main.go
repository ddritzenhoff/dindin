package main

import (
	"fmt"
	"log"
	"net"

	"github.com/ddritzenhoff/dindin/internal/configs"
	"github.com/ddritzenhoff/dindin/internal/cooking"
	"github.com/ddritzenhoff/dindin/internal/http/rest"
	"github.com/ddritzenhoff/dindin/internal/http/rpc"
	"github.com/ddritzenhoff/dindin/internal/person"
	"github.com/slack-go/slack"
	"google.golang.org/grpc"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
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

	restCfg, err := cfg.REST()
	if err != nil {
		log.Fatal(err)
	}

	h, err := rest.NewRESTService(restCfg, ps)
	if err != nil {
		log.Fatal(err)
	}
	go h.Start()

	grpcCfg, err := cfg.GRPC()
	if err != nil {
		log.Fatal(err)
	}
	// create a listener on TCP port 7777
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", grpcCfg.Host, grpcCfg.Port))
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
