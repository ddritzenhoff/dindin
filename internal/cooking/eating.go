package cooking

import (
	"github.com/slack-go/slack"
	"gorm.io/gorm"
	"time"
)

const day = time.Hour * 24

type Event struct {
	ID              uint
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt `gorm:"index"`
	CookingDay      int
	ChefSlackUID    string
	MealDescription string
	timeCreated     time.Time
	slackMessageID  string `gorm:"primaryKey"`
}

func (e *Event) IsEatingMessageExpired() bool {
	return time.Since(e.timeCreated) > day
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

func (s *EventService) GetEatingEvent(slackMessageID string) (_eatingEvent *Event, _exists bool) {
	event, err := s.store.get(slackMessageID)
	if err != nil {
		return nil, false
	}
	return event, true
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

func isEatingTomorrowBlock() slack.MsgOption {
	// Header Section
	headerText := slack.NewTextBlockObject("mrkdwn", "hey <!channel>, please react to this message (:thumbsup:) if you are eating tomorrow", false, false)
	headerSection := slack.NewSectionBlock(headerText, nil, nil)

	return slack.MsgOptionBlocks(
		headerSection,
	)
}
