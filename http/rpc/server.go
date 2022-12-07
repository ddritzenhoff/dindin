package rpc

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ddritzenhoff/dinny"
	"github.com/ddritzenhoff/dinny/http/rpc/pb"
	"github.com/ddritzenhoff/dinny/slack"
)

// Config represents the values needed for a gRPC server.
type Config struct {
	Host string
	Port string
}

// Server represents the gRPC server.
type Server struct {
	mealService   dinny.MealService
	memberService dinny.MemberService
	slackService  *slack.Service
	pb.UnimplementedSlackActionsServer
}

// NewServer creates a new gRPC server instance.
func NewServer(es dinny.MealService, ms dinny.MemberService, slackService *slack.Service) Server {
	return Server{mealService: es, memberService: ms, slackService: slackService}
}

// Ping generates a response to indicate http is working.
func (s *Server) Ping(_ context.Context, in *pb.EmptyMessage) (*pb.PingResponse, error) {
	return &pb.PingResponse{Message: "pong"}, nil
}

// EatingTomorrow generates triggers a message in the dinner-rotation Slack channel.
func (s *Server) EatingTomorrow(_ context.Context, in *pb.EmptyMessage) (*pb.EmptyMessage, error) {
	err := s.slackService.PostEatingTomorrow()
	if err != nil {
		return nil, fmt.Errorf("EatingTomorrow: %w", err)
	}
	return &pb.EmptyMessage{}, nil
}

// GetMembers fetches the recorded members from the database.
func (s *Server) GetMembers(_ *pb.EmptyMessage, stream pb.SlackActions_GetMembersServer) error {
	members, err := s.memberService.ListMembers()
	if err != nil {
		return fmt.Errorf("GetAllMembers: %w", err)
	}

	for _, member := range members {
		err := stream.Send(&pb.GetMembersResponse{
			FullName:  member.FullName,
			Slack_UID: member.SlackUID,
		})
		if err != nil {
			return fmt.Errorf("GetMembers: %w", err)
		}
	}

	return nil
}

// WeeklyUpdate triggers a weekly status update of the members in dinner rotation.
func (s *Server) WeeklyUpdate(_ context.Context, in *pb.EmptyMessage) (*pb.EmptyMessage, error) {
	return &pb.EmptyMessage{}, s.slackService.WeeklyUpdate()
}

// AssignCooks assigns cooks for up a maximum of the next week.
func (s *Server) AssignCooks(_ context.Context, acr *pb.AssignCooksRequest) (*pb.EmptyMessage, error) {
	assignments := acr.GetAssignments()
	for ii, assignment := range assignments {
		d := dinny.Date{
			Year:  int(assignment.Date.Year),
			Month: time.Month(assignment.Date.Month),
			Day:   int(assignment.Date.Day),
		}
		m, err := s.mealService.FindMealByDate(d)
		if errors.Is(err, dinny.ErrNotFound) {
			s.mealService.CreateMeal(&dinny.Meal{
				CookSlackUID: assignment.Slack_UID,
				Date:         d,
			})
		} else if err != nil {
			return &pb.EmptyMessage{}, fmt.Errorf("AssignCooks FindMealByDate %d: %w", ii, err)
		} else {
			err := s.mealService.UpdateMeal(m.ID, dinny.MealUpdate{
				ChefSlackUID: &assignment.Slack_UID,
			})
			if err != nil {
				return &pb.EmptyMessage{}, fmt.Errorf("AssignCooks UpdateMeal: %w", err)
			}
		}
	}
	return &pb.EmptyMessage{}, nil
}

// UpcomingCooks fetches the cooks set to cook over the next few days.
func (s *Server) UpcomingCooks(_ context.Context, in *pb.UpcomingCooksRequest) (*pb.UpcomingCooksResponse, error) {
	year, month, day := time.Now().Date()
	var meals []*pb.Meal
	for ii := 0; ii < int(in.DaysWanted); ii++ {
		d, _ := s.mealService.FindMealByDate(dinny.Date{Year: year, Month: month, Day: day + ii})
		if d == nil {
			continue
		}
		m, err := s.memberService.FindMemberBySlackUID(d.CookSlackUID)
		if err != nil {
			return nil, fmt.Errorf("UpcomingCooks: %w", err)
		}
		date := &pb.Date{
			Year:  int64(year),
			Month: int64(month),
			Day:   int64(day),
		}
		meal := pb.Meal{
			Date:            date,
			CookSlack_UID:   d.CookSlackUID,
			Description:     d.Description,
			SlackMessage_ID: d.SlackMessageID,
			FullName:        m.FullName,
		}
		meals = append(meals, &meal)
	}
	return &pb.UpcomingCooksResponse{Meals: meals}, nil
}
