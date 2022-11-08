package main

import (
	"context"
	"io"
	"log"

	"github.com/ddritzenhoff/dindin/internal/http/rpc/pb"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(cmdMembers)
}

var cmdMembers = &cobra.Command{
	Use:   "members",
	Short: "get the members of dinner rotation in [first name, last name | realName | displayName | SlackUID] order",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		conn, err := generateGRPCClientConnection()
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()
		slackClient := pb.NewSlackActionsClient(conn)
		stream, err := slackClient.GetMembers(context.Background(), &pb.EmptyMessage{})
		if err != nil {
			log.Fatalf("client.GetMembers failed: %v", err)
		}
		for {
			memberInfo, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("client.GetMembers failed: %v", err)
			}
			log.Printf("Real Name: %s\nDisplay Name: %s\nSlackUID: %s\n\n", memberInfo.GetRealName(), memberInfo.GetDisplayName(), memberInfo.GetSlackUID())
		}
	},
}
