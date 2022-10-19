package member

import (
	"fmt"
	"log"
	"math"
	"sort"
	"time"

	"github.com/ddritzenhoff/dindin/internal/day"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"gorm.io/gorm"
)

type Member struct {
	ID          uint
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	SlackUID    string         `gorm:"primaryKey"`
	FirstName   string
	LastName    string
	MealsEaten  uint
	MealsCooked uint
}

type Service struct {
	store         store
	eatingService *day.EventService
	slackClient   *slack.Client
	slackChannel  string
}

func NewService(db *gorm.DB, eatingService *day.EventService) (*Service, error) {
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

func (s *Service) GetMember(slackUID string) (*Member, error) {
	m, err := s.store.get(slackUID)
	if err != nil {
		return s.createMember(slackUID)
	}
	return m, nil
}

func (s *Service) createMember(slackUID string) (*Member, error) {
	slackUser, err := s.slackClient.GetUserInfo(slackUID)
	if err != nil {
		return nil, fmt.Errorf("GetUserInfo: %w", err)
	}
	m := &Member{SlackUID: slackUID, FirstName: slackUser.Profile.FirstName, LastName: slackUser.Profile.LastName, MealsEaten: 0, MealsCooked: 0}
	err = s.store.create(m)
	if err != nil {
		return nil, fmt.Errorf("db create: %w", err)
	}
	return m, nil
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
	m, err := s.GetMember(slackUID)
	if err != nil {
		return
	}

	m.MealsEaten += 1
	err = s.store.update(m)
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
	m, err := s.GetMember(slackUID)
	if err != nil {
		return
	}
	if m.MealsEaten > 0 {
		m.MealsEaten -= 1
	}
	err = s.store.update(m)
	if err != nil {
		return
	}
}

func (s *Service) GetAllMembers() ([]Member, error) {
	members, err := s.store.getAll()
	if err != nil {
		return nil, err
	}
	return members, nil
}

func weeklyUpdateBlock(members []Member) slack.MsgOption {

	var sectionBlocks []slack.Block
	// Header Section
	headerText := slack.NewTextBlockObject("mrkdwn", "dinner rotation members with the *worst* meals eaten to meals cooked ratios:", false, false)
	headerSection := slack.NewSectionBlock(headerText, nil, nil)
	sectionBlocks = append(sectionBlocks, *headerSection)

	for ii, member := range members {
		if ii > 10 {
			break
		}

		realNameField := slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*Real Name:*\n%s %s", member.FirstName, member.LastName), false, false)
		slackNameField := slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*Slack Name:*\n<@%s>", member.SlackUID), false, false)
		ratioField := slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*Ratio Status:*\n%s", ratioStatus(member.MealsEaten, member.MealsCooked)), false, false)

		fieldSlice := make([]*slack.TextBlockObject, 0)
		fieldSlice = append(fieldSlice, realNameField)
		fieldSlice = append(fieldSlice, slackNameField)
		fieldSlice = append(fieldSlice, ratioField)
		fieldsSection := slack.NewSectionBlock(nil, fieldSlice, nil)
		sectionBlocks = append(sectionBlocks, *fieldsSection)
	}

	return slack.MsgOptionBlocks(sectionBlocks...)
}

func ratioStatus(mealsEaten uint, mealsCooked uint) string {
	if mealsCooked == 0 {
		if mealsEaten > 0 {
			return "Infinity! You've eaten but never cooked"
		} else {
			return "Neither cooked nor eaten"
		}
	} else {
		return fmt.Sprintf("%.3f", float32(mealsEaten)/float32(mealsCooked))
	}
}

func mealsEatenToMealsCooked(mealsEaten uint, mealsCooked uint) float32 {
	if mealsCooked == 0 {
		if mealsEaten > 0 {
			return math.MaxFloat32
		} else {
			return 0
		}
	} else {
		return float32(mealsEaten) / float32(mealsCooked)
	}
}

func (s *Service) WeeklyUpdate() error {
	members, err := s.GetAllMembers()
	if err != nil {
		return fmt.Errorf("GetAllMembers: %w", err)
	}

	sort.Slice(members, func(ii, jj int) bool {
		return mealsEatenToMealsCooked(members[ii].MealsEaten, members[ii].MealsCooked) > mealsEatenToMealsCooked(members[jj].MealsEaten, members[jj].MealsCooked)
	})

	_, _, err = s.slackClient.PostMessage(s.slackChannel, weeklyUpdateBlock(members))
	if err != nil {
		return fmt.Errorf("slack PostMessage: %w", err)
	}
	return nil
}
