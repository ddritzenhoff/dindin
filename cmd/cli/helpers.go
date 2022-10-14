package main

import (
	"fmt"
	"log"

	"github.com/ddritzenhoff/dindin/internal/configs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func generateGRPCClientConnection() (*grpc.ClientConn, error) {
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
	return conn, err
}
