package rest

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ddritzenhoff/dinny"
	"github.com/ddritzenhoff/dinny/slack"
)

// Config represents the values needed to start an HTTP server.
type Config struct {
	Host string `toml:"host"`
	Port string `toml:"port"`
}

// Server represents an HTTP server.
type Server struct {
	logger *log.Logger
	server *http.Server
	config *Config
}

// Start starts the HTTP server on a specific host and port.
func (s *Server) Start() {
	s.logger.Printf("HTTP server listening on host %s and port %s\n", s.config.Host, s.config.Port)
	err := s.server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

// NewServer creates a new HTTP server.
func NewServer(logger *log.Logger, cfg *Config, memberService dinny.MemberService, slackService *slack.Service) *Server {
	h := NewHandlers(logger, memberService, slackService)
	httpServer := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Handler: h.routes(),
	}

	return &Server{
		logger: logger,
		server: httpServer,
		config: cfg,
	}
}
