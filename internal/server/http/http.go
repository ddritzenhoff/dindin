package http

import (
	"fmt"
	"github.com/ddritzenhoff/dindin/internal/person"
	"log"
	"net/http"
)

type HTTP struct {
	server *http.Server
	config *Config
}

func (h *HTTP) Start() {
	log.Printf("Listening on host %s and port %s\n", h.config.Host, h.config.Port)
	err := h.server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

type Config struct {
	Host string
	Port string
}

func NewHTTPService(cfg *Config, personService *person.Service) (*HTTP, error) {
	h := &Handlers{personService: personService}
	httpServer := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Handler: h.routes(),
	}

	return &HTTP{
		server: httpServer,
		config: cfg,
	}, nil
}
