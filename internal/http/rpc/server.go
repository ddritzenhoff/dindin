package rpc

import (
	"context"

	"github.com/ddritzenhoff/dindin/internal/day"
	"github.com/ddritzenhoff/dindin/internal/member"
	"github.com/slack-go/slack"
)

// Server represents the gRPC http
type Server struct {
	eatingService *day.EventService
	memberService *member.Service
	slackClient   *slack.Client
	UnimplementedSlackActionsServer
}

func NewServer(es *day.EventService, ms *member.Service, slackClient *slack.Client) Server {
	return Server{eatingService: es, memberService: ms}
}

// Ping generates a response to indicate http is working
func (s *Server) Ping(_ context.Context, in *PingMessage) (*PingMessage, error) {
	return &PingMessage{Message: in.Message + " received"}, nil
}

// EatingTomorrow generates triggers a message in the dinner-rotation Slack channel
func (s *Server) EatingTomorrow(_ context.Context, in *EatingTomorrowRequest) (*EatingTomorrowResponse, error) {
	err := s.eatingService.PostEatingTomorrow(in.SlackUID)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (s *Server) GetMembers(_ *GetMembersRequest, stream SlackActions_GetMembersServer) error {
	members, err := s.memberService.GetAllMembers()
	if err != nil {
		return err
	}

	var membersSlackUIDs []string
	for _, member := range members {
		membersSlackUIDs = append(membersSlackUIDs, member.SlackUID)
	}

	slackMembers, err := s.slackClient.GetUsersInfo(membersSlackUIDs...)
	if err != nil {
		return err
	}

	for _, slackMember := range *slackMembers {
		stream.Send(&GetMembersResponse{
			FirstName:   slackMember.Profile.FirstName,
			LastName:    slackMember.Profile.LastName,
			RealName:    slackMember.Profile.RealNameNormalized,
			DisplayName: slackMember.Profile.DisplayNameNormalized,
			SlackUID:    slackMember.ID,
		})
	}

	return nil
}
