package pb

import (
	"context"
	"github.com/ddritzenhoff/dindin/internal/cooking"
)

// Server represents the gRPC server
type Server struct {
	eatingService *cooking.EventService
	UnimplementedSlackActionsServer
}

func NewServer(es *cooking.EventService) Server {
	return Server{eatingService: es}
}

// Ping generates a response to indicate server is working
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
