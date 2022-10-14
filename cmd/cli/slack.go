package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ddritzenhoff/dindin/internal/configs"
	"github.com/ddritzenhoff/dindin/internal/http/rpc"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func cmdSlackMessage(cmd *cobra.Command, args []string) {
	cfg, err := configs.NewConfigService()
	if err != nil {
		log.Fatal(err)
	}
	grpcCfg, err := cfg.GRPC()
	if err != nil {
		log.Fatal(err)
	}

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", grpcCfg.Host, grpcCfg.Port), opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	slackClient := rpc.NewSlackActionsClient(conn)
	msg, err := slackClient.Ping(context.Background(), &rpc.PingMessage{Message: "sent from the client"})
	fmt.Printf("Got a response from the server: %s", msg.Message)
}
