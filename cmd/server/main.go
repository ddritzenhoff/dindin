package main

import (
	"fmt"
	"log"
	"net"

	"database/sql"

	"github.com/ddritzenhoff/dindin/internal/configs"
	"github.com/ddritzenhoff/dindin/internal/cooking"
	"github.com/ddritzenhoff/dindin/internal/http/rest"
	"github.com/ddritzenhoff/dindin/internal/http/rpc"
	"github.com/ddritzenhoff/dindin/internal/http/rpc/pb"
	"github.com/ddritzenhoff/dindin/internal/member"
	_ "github.com/mattn/go-sqlite3"
	"google.golang.org/grpc"
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

	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	slackConfig, err := cfg.SlackConfig()
	if err != nil {
		log.Fatal(err)
	}
	ces, err := cooking.NewService(db, slackConfig)
	if err != nil {
		log.Fatal(err)
	}

	ms, err := member.NewService(db, ces)
	if err != nil {
		log.Fatal(err)
	}

	restCfg, err := cfg.REST()
	if err != nil {
		log.Fatal(err)
	}

	h, err := rest.NewRESTService(restCfg, ms)
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
	s := rpc.NewServer(ces, ms, slackConfig)
	// create a gRPC http object
	grpcServer := grpc.NewServer()
	// attach the Ping service to the http
	pb.RegisterSlackActionsServer(grpcServer, &s)
	// start the http
	log.Printf("gRPC server listening on host %s and port %s\n", grpcCfg.Host, grpcCfg.Port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
