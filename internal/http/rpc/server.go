package rpc

import (
	"context"
	"fmt"

	"github.com/ddritzenhoff/dindin/internal/configs"
	"github.com/ddritzenhoff/dindin/internal/cooking"
	"github.com/ddritzenhoff/dindin/internal/http/rpc/pb"
	"github.com/ddritzenhoff/dindin/internal/member"
)

// Server represents the gRPC http
type Server struct {
	cookingService *cooking.Service
	memberService  *member.Service
	slackCfg       *configs.SlackConfig
	pb.UnimplementedSlackActionsServer
}

func NewServer(es *cooking.Service, ms *member.Service, slackCfg *configs.SlackConfig) Server {
	return Server{cookingService: es, memberService: ms, slackCfg: slackCfg}
}

// Ping generates a response to indicate http is working
func (s *Server) Ping(_ context.Context, in *pb.EmptyMessage) (*pb.PingResponse, error) {
	return &pb.PingResponse{Message: "pong"}, nil
}

// EatingTomorrow generates triggers a message in the dinner-rotation Slack channel
func (s *Server) EatingTomorrow(_ context.Context, in *pb.EmptyMessage) (*pb.EmptyMessage, error) {
	err := s.cookingService.PostEatingTomorrow()
	if err != nil {
		return nil, fmt.Errorf("PostEatingTomorrow: %w", err)
	}
	return &pb.EmptyMessage{}, nil
}

func (s *Server) GetMembers(_ *pb.EmptyMessage, stream pb.SlackActions_GetMembersServer) error {
	members, err := s.memberService.GetAllMembers()
	if err != nil {
		return fmt.Errorf("GetAllMembers: %w", err)
	}

	var membersSlackUIDs []string
	for _, member := range members {
		membersSlackUIDs = append(membersSlackUIDs, member.SlackUID)
	}

	slackMembers, err := s.slackCfg.Client.GetUsersInfo(membersSlackUIDs...)
	if err != nil {
		return fmt.Errorf("GetUsersInfo: %w", err)
	}

	for _, slackMember := range *slackMembers {
		stream.Send(&pb.GetMembersResponse{
			FirstName:   slackMember.Profile.FirstName,
			LastName:    slackMember.Profile.LastName,
			RealName:    slackMember.Profile.RealNameNormalized,
			DisplayName: slackMember.Profile.DisplayNameNormalized,
			SlackUID:    slackMember.ID,
		})
	}

	return nil
}

func (s *Server) WeeklyUpdate(_ context.Context, in *pb.EmptyMessage) (*pb.EmptyMessage, error) {
	return &pb.EmptyMessage{}, s.memberService.WeeklyUpdate()
}

func (s *Server) AssignCooks(_ context.Context, acr *pb.AssignCooksRequest) (*pb.EmptyMessage, error) {
	return &pb.EmptyMessage{}, s.cookingService.AssignCooks(acr.GetCookings())
}
