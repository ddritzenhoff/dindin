package main

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/ddritzenhoff/dinny/configs"
	"github.com/ddritzenhoff/dinny/http/rest"
	"github.com/ddritzenhoff/dinny/http/rpc"
	"github.com/ddritzenhoff/dinny/http/rpc/pb"
	"github.com/ddritzenhoff/dinny/slack"
	"github.com/ddritzenhoff/dinny/sqlite"
	"github.com/ddritzenhoff/dinny/sqlite/gen"
	_ "github.com/mattn/go-sqlite3"
	"google.golang.org/grpc"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

// run initializes the member, meal, and Slack services and starts the REST and GRPC servers.
func run() error {
	logger := log.New(os.Stdout, "DEBUG: ", log.LstdFlags)

	cfg, err := configs.NewConfigService()
	if err != nil {
		log.Fatal(err)
	}

	slackConfig, err := cfg.Slack()
	if err != nil {
		log.Fatal(err)
	}

	dbName, err := cfg.DBName()
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// create tables
	_, err = db.ExecContext(context.Background(), sqlite.Schema)
	if err != nil {
		log.Fatal(err)
	}

	queries := gen.New(db)

	mealService := sqlite.NewMealService(queries, db)

	memberService := sqlite.NewMemberService(queries, db)

	slackService := slack.NewService(slackConfig, mealService, memberService)

	restCfg, err := cfg.REST()
	if err != nil {
		log.Fatal(err)
	}

	h, err := rest.NewServer(logger, restCfg, memberService, slackService)
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
	s := rpc.NewServer(mealService, memberService, slackService)
	// create a gRPC http object
	grpcServer := grpc.NewServer()
	// attach the Ping service to the http
	pb.RegisterSlackActionsServer(grpcServer, &s)
	// start the http
	log.Printf("gRPC server listening on host %s and port %s\n", grpcCfg.Host, grpcCfg.Port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
	return nil
}
