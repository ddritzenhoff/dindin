package rest

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/ddritzenhoff/dinny"
	"github.com/ddritzenhoff/dinny/slack"
	"github.com/slack-go/slack/slackevents"
)

// Handlers represents the handler to an HTTP server.
type Handlers struct {
	logger        *log.Logger
	memberService dinny.MemberService
	slackService  *slack.Service
}

// NewHandlers creates a new Handlers instance.
func NewHandlers(logger *log.Logger, memberService dinny.MemberService, slackService *slack.Service) *Handlers {
	return &Handlers{
		logger:        logger,
		memberService: memberService,
		slackService:  slackService,
	}
}

// routes returns the endpoints with dedicated handlers.
func (h *Handlers) routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/event", h.handleSlackEvent)
	mux.HandleFunc("/ping", h.handlePing)
	return mux
}

// handlePing returns 'pong' to every request to indicate the server being alive.
func (h *Handlers) handlePing(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}

// handleSlackEvent handles Slack events like a Slack member liking a message.
func (h *Handlers) handleSlackEvent(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return
	}
	event, err := slackevents.ParseEvent(body, slackevents.OptionNoVerifyToken())
	if err != nil {
		log.Println("Unable to parse event.")
		return
	}

	if event.Type == slackevents.URLVerification {
		var r *slackevents.ChallengeResponse
		err := json.Unmarshal(body, &r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_, err = w.Write([]byte(r.Challenge))
		if err != nil {
			log.Println("Unable to write response")
			return
		}
	}

	if event.Type == slackevents.CallbackEvent {
		switch innerEvent := event.InnerEvent.Data.(type) {
		case *slackevents.ReactionAddedEvent:
			err := h.slackService.ReactionAddedEvent(innerEvent)
			if err != nil {
				h.logger.Printf("%s", err.Error())
			}
		case *slackevents.ReactionRemovedEvent:
			err := h.slackService.ReactionRemovedEvent(innerEvent)
			if err != nil {
				h.logger.Printf("%s", err.Error())
			}
		}
	}
}
