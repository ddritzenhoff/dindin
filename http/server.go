package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/ddritzenhoff/dinny"
	"github.com/ddritzenhoff/dinny/slack"
	"github.com/go-chi/chi/v5"
	"github.com/slack-go/slack/slackevents"
)

// Server represents the dinny REST server.
type Server struct {
	ln     net.Listener
	server *http.Server
	router *chi.Mux

	Logger *log.Logger

	// Bind address the server's listener.
	Addr string

	// Servics used by the various HTTP routes.
	MemberService dinny.MemberService
	MealService   dinny.MealService
	SlackService  slack.Service
}

// NewServer creates a new dinny REST server instance.
func NewServer(logger *log.Logger, addr string, memberService dinny.MemberService, mealService dinny.MealService, slackService slack.Service) *Server {
	s := &Server{
		server:        &http.Server{},
		router:        chi.NewRouter(),
		Logger:        logger,
		Addr:          addr,
		MemberService: memberService,
		MealService:   mealService,
		SlackService:  slackService,
	}

	s.router.Put("/event", s.handleSlackEvent)
	s.router.Get("/ping", s.handlePing)
	s.router.Route("/cmd", func(r chi.Router) {
		r.Put("/assign-cooks", s.handleAssignCooks)
		r.Get("/eating-tomorrow", s.handleEatingTomorrow)
		r.Get("/members", s.handleMembers)
		r.Get("/upcoming-cooks", s.handleUpcomingCooks)
		r.Get("/weekly-update", s.handleWeeklyUpdate)
	})

	s.server.Handler = s.router

	return s
}

// Open creates a listener on a specific port and address.
func (s *Server) Open() error {
	var err error
	s.ln, err = net.Listen("tcp", s.Addr)
	if err != nil {
		return fmt.Errorf("Open net.Listen: %w", err)
	}
	go s.server.Serve(s.ln)
	s.Logger.Printf("Server listening on %s", s.ln.Addr())
	return nil
}

// handleSlackEvent handles Slack events like a Slack member liking a message.
func (s *Server) handleSlackEvent(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return
	}
	event, err := slackevents.ParseEvent(body, slackevents.OptionNoVerifyToken())
	if err != nil {
		log.Println("Unable to parse event.")
		return
	}

	switch event.Type {
	case slackevents.URLVerification:
		s.handleSlackURLVerification(w, r, body)
	case slackevents.CallbackEvent:
		s.handleCallbackEvent(w, r, event)
	default:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		s.Logger.Printf("handleSlackEvent: %s", err.Error())
		return
	}
}

// handleSlackURLVerification verifies the slack request.
func (s *Server) handleSlackURLVerification(w http.ResponseWriter, r *http.Request, body []byte) {
	var ch *slackevents.ChallengeResponse
	err := json.Unmarshal(body, &r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = w.Write([]byte(ch.Challenge))
	if err != nil {
		s.Logger.Println("handleSlackURLVerification: Unable to write response")
		return
	}
}

// handleCallbackEvent checks the slack reaction.
func (s *Server) handleCallbackEvent(w http.ResponseWriter, r *http.Request, event slackevents.EventsAPIEvent) {
	switch innerEvent := event.InnerEvent.Data.(type) {
	case *slackevents.ReactionAddedEvent:
		err := s.SlackService.ReactionAddedEvent(innerEvent)
		if err != nil {
			s.Logger.Printf("%s", err.Error())
		}
	case *slackevents.ReactionRemovedEvent:
		err := s.SlackService.ReactionRemovedEvent(innerEvent)
		if err != nil {
			s.Logger.Printf("%s", err.Error())
		}
	}
}

// CookAssignment represents the assignment of a cook on a specific date.
type CookAssignment struct {
	Date         dinny.Date `json:"date"`
	CookSlackUID string     `json:"cookSlackUID"`
}

// AssignCooksRequest represents the a collection of cook assignments.
type AssignCooksRequest struct {
	CookAssignments []CookAssignment `json:"cooks"`
}

// handleAssignCooks represents a handler for assigning multiple cooks.
func (s *Server) handleAssignCooks(w http.ResponseWriter, r *http.Request) {
	var req *AssignCooksRequest
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		s.Logger.Printf("handleAssignCooks: %s", err.Error())
		return
	}
	for _, assignment := range req.CookAssignments {
		m, err := s.MealService.FindMealByDate(assignment.Date)
		if errors.Is(err, dinny.ErrNotFound) {
			err := s.MealService.CreateMeal(&dinny.Meal{
				CookSlackUID: assignment.CookSlackUID,
				Date:         assignment.Date,
			})
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error()))
				s.Logger.Printf("handleAssignCooks MealService.CreateMeal: %s", err.Error())
				return
			}
		} else if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			s.Logger.Printf("handleAssignCooks MealService.FindMealByDate: %s", err.Error())
			return
		} else {
			err := s.MealService.UpdateMeal(m.ID, dinny.MealUpdate{
				ChefSlackUID: &assignment.CookSlackUID,
			})
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error()))
				s.Logger.Printf("handleAssignCooks MealService.UpdateMeal: %s", err.Error())
				return
			}

		}
	}

}

// handleEatingTomorrow is a handler for the eating_tomorrow command.
func (s *Server) handleEatingTomorrow(w http.ResponseWriter, r *http.Request) {
	err := s.SlackService.PostEatingTomorrow()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		s.Logger.Printf("EatingTomorrow: %s", err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
}

// handleMembers is a handler for the members command.
func (s *Server) handleMembers(w http.ResponseWriter, r *http.Request) {
	members, err := s.MemberService.ListMembers()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		s.Logger.Printf("Members: %s", err.Error())
		return
	}

	e := json.NewEncoder(w)
	for _, member := range members {
		e.Encode(member)
	}
	w.WriteHeader(http.StatusOK)
}

// handlePing returns 'pong' to every request to indicate the server being alive.
func (s *Server) handlePing(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}

// handleUpcomingCooks is a handler for the upcoming_cooks command.
func (s *Server) handleUpcomingCooks(w http.ResponseWriter, r *http.Request) {
	year, month, day := time.Now().Date()
	daysWanted, err := strconv.Atoi(r.URL.Query().Get("daysWanted"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		s.Logger.Printf("handleUpcomingCooks strconv.Atoi: %s", err.Error())
		return
	}
	e := json.NewEncoder(w)
	for ii := 0; ii < int(daysWanted); ii++ {
		date := dinny.Date{Year: year, Month: month, Day: day + ii}
		meal, err := s.MealService.FindMealByDate(date)
		if errors.Is(dinny.ErrNotFound, err) {
			e.Encode(dinny.Meal{CookSlackUID: "NOT SET", Date: date})
		} else if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			s.Logger.Printf("handleUpcomingCooks FindMealByDate: %s", err.Error())
			return
		} else {
			e.Encode(*meal)
		}
	}
}

// handleWeeklyUpdate is a handler for the weekly_update command.
func (s *Server) handleWeeklyUpdate(w http.ResponseWriter, r *http.Request) {
	err := s.SlackService.WeeklyUpdate()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		s.Logger.Printf("handleWeeklyUpdate SlackService.WeeklyUpdate: %s", err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
}
