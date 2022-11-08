package rest

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/ddritzenhoff/dindin/internal/member"
	"github.com/slack-go/slack/slackevents"
)

type Handlers struct {
	personService *member.Service
}

func (h *Handlers) routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/event", h.handleSlackEvent)
	mux.HandleFunc("/ping", h.handlePing)
	return mux
}

func (h *Handlers) handlePing(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}

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
			h.personService.ReactionAddedEvent(innerEvent)
		case *slackevents.ReactionRemovedEvent:
			h.personService.ReactionRemovedEvent(innerEvent)
		}
	}
}
