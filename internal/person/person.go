package person

import (
	"github.com/ddritzenhoff/dindin/internal/cooking"
	"github.com/slack-go/slack/slackevents"
	"gorm.io/gorm"
	"log"
	"time"
)

type Person struct {
	ID          uint
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	SlackUID    string         `gorm:"primaryKey"`
	MealsEaten  uint
	MealsCooked uint
}

type Service struct {
	store         store
	eatingService *cooking.EventService
}

func NewService(db *gorm.DB, eatingService *cooking.EventService) (*Service, error) {
	pStore, err := newStore(db)
	if err != nil {
		return nil, err
	}

	return &Service{store: pStore, eatingService: eatingService}, nil
}

func (s *Service) LikedMessage(slackUID string) error {
	person, err := s.store.get(slackUID)
	if err != nil {
		return err
	}
	person.MealsEaten += 1
	err = s.store.update(person)
	if err != nil {
		log.Println("wasn't able to update meals eaten")
		return err
	}
	return nil
}

func (s *Service) GetPerson(slackUID string) (p *Person, _exists bool) {
	p, err := s.store.get(slackUID)
	if err != nil {
		return nil, false
	}
	return p, true
}

// ReactionAddedEvent determines whether to add +1 to the meals eaten of this specific user
func (s *Service) ReactionAddedEvent(e *slackevents.ReactionAddedEvent) {
	// Don't bother if the reaction isn't a like
	if e.Reaction != "+1" {
		return
	}

	// check to see whether the reaction was on a cooking event message
	eventMessageID := e.Item.Timestamp
	eatingEvent, exists := s.eatingService.GetEatingEvent(eventMessageID)
	if !exists {
		return
	}

	// if the reaction was on an expired cooking event message, don't do anything
	if eatingEvent.IsEatingMessageExpired() {
		return
	}

	// add the user if he doesn't exist yet
	slackUID := e.User
	p, exists := s.GetPerson(slackUID)
	if !exists {
		p = &Person{SlackUID: slackUID}
		err := s.store.create(p)
		if err != nil {
			return
		}
	}

	p.MealsEaten += 1
	err := s.store.update(p)
	if err != nil {
		return
	}
}

// ReactionRemovedEvent determines whether to add -1 to the meals eaten of this specific user
func (s *Service) ReactionRemovedEvent(e *slackevents.ReactionRemovedEvent) {
	// Don't bother if the reaction isn't a like
	if e.Reaction != "+1" {
		return
	}

	// check to see whether the reaction was on a cooking event message
	eventMessageID := e.Item.Timestamp
	eatingEvent, exists := s.eatingService.GetEatingEvent(eventMessageID)
	if !exists {
		return
	}

	// if the reaction was on an expired cooking event message, don't do anything
	if eatingEvent.IsEatingMessageExpired() {
		return
	}

	// add the user if he doesn't exist yet
	slackUID := e.User
	p, exists := s.GetPerson(slackUID)
	if !exists {
		p = &Person{SlackUID: slackUID}
		err := s.store.create(p)
		if err != nil {
			return
		}
	}

	if p.MealsEaten > 0 {
		p.MealsEaten -= 1
	}
	err := s.store.update(p)
	if err != nil {
		return
	}
}
