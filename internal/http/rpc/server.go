package rpc

import (
	"context"
	"fmt"

	"github.com/ddritzenhoff/dindin/internal/day"
	"github.com/ddritzenhoff/dindin/internal/http/rpc/pb"
	"github.com/ddritzenhoff/dindin/internal/member"
	"github.com/slack-go/slack"
)

// Server represents the gRPC http
type Server struct {
	eventService  *day.EventService
	memberService *member.Service
	slackClient   *slack.Client
	pb.UnimplementedSlackActionsServer
}

func NewServer(es *day.EventService, ms *member.Service, slackClient *slack.Client) Server {
	return Server{eventService: es, memberService: ms}
}

// Ping generates a response to indicate http is working
func (s *Server) Ping(_ context.Context, in *pb.PingMessage) (*pb.PingMessage, error) {
	return &pb.PingMessage{Message: in.Message + " received"}, nil
}

// EatingTomorrow generates triggers a message in the dinner-rotation Slack channel
func (s *Server) EatingTomorrow(_ context.Context, in *pb.EatingTomorrowRequest) (*pb.EatingTomorrowResponse, error) {
	err := s.eventService.PostEatingTomorrow(in.SlackUID)
	if err != nil {
		return nil, fmt.Errorf("PostEatingTomorrow: %w", err)
	}
	return &pb.EatingTomorrowResponse{}, nil
}

func (s *Server) GetMembers(_ *pb.GetMembersRequest, stream pb.SlackActions_GetMembersServer) error {
	members, err := s.memberService.GetAllMembers()
	if err != nil {
		return fmt.Errorf("GetAllMembers: %w", err)
	}

	var membersSlackUIDs []string
	for _, member := range members {
		membersSlackUIDs = append(membersSlackUIDs, member.SlackUID)
	}

	slackMembers, err := s.slackClient.GetUsersInfo(membersSlackUIDs...)
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

func (s *Server) WeeklyUpdate(_ context.Context, in *pb.WeeklyUpdateRequest) (*pb.WeeklyUpdateResponse, error) {
	return &pb.WeeklyUpdateResponse{}, s.memberService.WeeklyUpdate()
}

func (s *Server) AssignCooks(_ context.Context, acr *pb.AssignCooksRequest) (*pb.AssignCooksResponse, error) {
	return &pb.AssignCooksResponse{}, s.eventService.AssignCooks(acr.GetCookingDays())
}
