package rpc

import (
	"context"
	"fmt"
	"time"

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
	members, err := s.memberService.AllMembers()
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
			Slack_UID:   slackMember.ID,
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

func (s *Server) UpcomingCooks(_ context.Context, in *pb.UpcomingCooksRequest) (*pb.UpcomingCooksResponse, error) {
	year, month, day := time.Now().Date()
	var cooks []*pb.Cook
	for ii := 0; ii < int(in.DaysWanted); ii++ {
		d, _ := s.cookingService.CookingByDate(year, month, day+ii)
		if d == nil {
			continue
		}
		m, err := s.memberService.GetMember(d.ChefSlackUID)
		if err != nil {
			return nil, fmt.Errorf("UpcomingCooks: %w", err)
		}
		c := pb.Cook{
			CookingTime:     d.CookingTime,
			ChefSlack_UID:   d.ChefSlackUID,
			MealDescription: d.MealDescription,
			SlackMessage_ID: d.SlackMessageID,
			FirstName:       m.FirstName,
			LastName:        m.LastName,
		}
		cooks = append(cooks, &c)
	}
	return &pb.UpcomingCooksResponse{Cooks: cooks}, nil
}
