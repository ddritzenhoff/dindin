package rest

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ddritzenhoff/dindin/internal/configs"
	"github.com/ddritzenhoff/dindin/internal/member"
)

type Server struct {
	server *http.Server
	config *configs.HTTP
}

func (h *Server) Start() {
	log.Printf("HTTP server listening on host %s and port %s\n", h.config.Host, h.config.Port)
	err := h.server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func NewRESTService(logger *log.Logger, cfg *configs.HTTP, memberService *member.Service) (*Server, error) {
	h := NewHandlers(logger, memberService)
	httpServer := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Handler: h.routes(),
	}

	return &Server{
		server: httpServer,
		config: cfg,
	}, nil
}
