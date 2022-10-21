package day

import (
	"fmt"
	"time"

	"github.com/ddritzenhoff/dindin/internal/http/rpc/pb"
	"github.com/slack-go/slack"
	"gorm.io/gorm"
)

const DAY_LENGTH_SECONDS = time.Hour * 24

type Event struct {
	ID              uint
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt `gorm:"index"`
	CookingDay      int
	CookingMonth    int
	CookingYear     int
	ChefSlackUID    string
	MealDescription string
	timeCreated     time.Time
	slackMessageID  string
}

func (e *Event) IsEatingMessageExpired() bool {
	return time.Since(e.timeCreated) > DAY_LENGTH_SECONDS
}

// EventService struct holds all the dependencies required for the CookingEvent struct and exposes all services
// provided by this package as its methods
type EventService struct {
	store        store
	slackChannel string
	slackClient  *slack.Client
}

func NewEventService(db *gorm.DB, slackChannel string, slackClient *slack.Client) (*EventService, error) {
	pStore, err := newStore(db)
	if err != nil {
		return nil, err
	}
	return &EventService{store: pStore, slackChannel: slackChannel, slackClient: slackClient}, nil
}

func (s *EventService) GetEatingEvent(slackMessageID string) (eatingEvent *Event, exists bool) {
	event, err := s.store.get(slackMessageID)
	if err != nil {
		return nil, false
	}
	return event, true
}

func isEatingTomorrowBlock() slack.MsgOption {
	// Header Section
	headerText := slack.NewTextBlockObject("mrkdwn", "hey <!channel>, please react to this message (:thumbsup:) if you are eating tomorrow", false, false)
	headerSection := slack.NewSectionBlock(headerText, nil, nil)

	return slack.MsgOptionBlocks(
		headerSection,
	)
}

func (s *EventService) PostEatingTomorrow(slackUID string) error {
	_, respTimestamp, err := s.slackClient.PostMessage(s.slackChannel, isEatingTomorrowBlock())
	if err != nil {
		return err
	}
	now := time.Now()
	e := Event{
		CookingDay:      int(now.Weekday()),
		ChefSlackUID:    slackUID,
		MealDescription: "",
		timeCreated:     now,
		slackMessageID:  respTimestamp,
	}

	err = s.store.create(&e)
	if err != nil {
		return err
	}

	return nil
}

func (s *EventService) AssignCooks(cookingDays []*pb.CookingDay) error {
	for ii, cookingDay := range cookingDays {
		event, err := s.store.getByDateOrCreate(int(cookingDay.Year), int(cookingDay.Month), int(cookingDay.Day))
		if err != nil {
			return fmt.Errorf("getByDate %d: %w", ii, err)
		}
		event.ChefSlackUID = cookingDay.SlackUID
		s.store.update(event)
	}
	return nil
}
