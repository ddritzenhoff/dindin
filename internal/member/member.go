package member

import (
	"database/sql"
	"errors"
	"fmt"
	"math"
	"sort"

	"github.com/ddritzenhoff/dindin/internal/configs"
	"github.com/ddritzenhoff/dindin/internal/cooking"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

type Member struct {
	ID          int64
	CreatedAt   int64
	UpdatedAt   int64
	SlackUID    string
	FirstName   string
	LastName    string
	MealsEaten  int64
	MealsCooked int64
}

type Service struct {
	repository     repository
	cookingService *cooking.Service
	slackCfg       *configs.SlackConfig
}

func NewService(db *sql.DB, cs *cooking.Service) (*Service, error) {
	r, err := newRepository(db)
	if err != nil {
		return nil, fmt.Errorf("NewService: %w", err)
	}

	return &Service{repository: r, cookingService: cs}, nil
}

func (s *Service) LikedMessage(slackUID string) error {
	person, err := s.repository.getBySlackUID(slackUID)
	if err != nil {
		return err
	}
	person.MealsEaten += 1
	_, err = s.repository.updateMealsEaten(person.ID, person.MealsEaten)
	if err != nil {
		return fmt.Errorf("LikedMessage: %w", err)
	}
	return nil
}

func (s *Service) GetMemberOrCreate(slackUID string) (*Member, error) {
	m, err := s.repository.getBySlackUID(slackUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			m, err := s.createMember(slackUID)
			if err != nil {
				return nil, fmt.Errorf("GetMember createMember: %w", err)
			}
			return m, nil
		}
		return nil, fmt.Errorf("GetMember: %w", err)
	}
	return m, nil
}

func (s *Service) createMember(slackUID string) (*Member, error) {
	slackUser, err := s.slackCfg.Client.GetUserInfo(slackUID)
	if err != nil {
		return nil, fmt.Errorf("createMember GetUserInfo: %w", err)
	}
	m, err := s.repository.create(Member{SlackUID: slackUID, FirstName: slackUser.Profile.FirstName, LastName: slackUser.Profile.LastName, MealsEaten: 0, MealsCooked: 0})
	if err != nil {
		return nil, fmt.Errorf("createMember create: %w", err)
	}
	return m, nil
}

// TODO (ddritzenhoff) either add logging or figure out a way to display the error to the user
// ReactionAddedEvent determines whether to add +1 to the meals eaten of this specific user
func (s *Service) ReactionAddedEvent(e *slackevents.ReactionAddedEvent) {
	// Don't bother if the reaction isn't a like
	if e.Reaction != "+1" {
		return
	}

	// check to see whether the reaction was on a cooking event message
	eventMessageID := e.Item.Timestamp
	eatingEvent, exists := s.cookingService.GetEatingEvent(eventMessageID)
	if !exists {
		return
	}

	// if the reaction was on an expired cooking event message, don't do anything
	if eatingEvent.IsEatingMessageExpired() {
		return
	}

	// add the user if he doesn't exist yet
	slackUID := e.User
	m, err := s.GetMemberOrCreate(slackUID)
	if err != nil {
		return
	}

	m.MealsEaten += 1
	_, err = s.repository.updateMealsEaten(m.ID, m.MealsEaten)
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
	eatingEvent, exists := s.cookingService.GetEatingEvent(eventMessageID)
	if !exists {
		return
	}

	// if the reaction was on an expired cooking event message, don't do anything
	if eatingEvent.IsEatingMessageExpired() {
		return
	}

	// add the user if he doesn't exist yet
	slackUID := e.User
	m, err := s.GetMemberOrCreate(slackUID)
	if err != nil {
		return
	}
	if m.MealsEaten > 0 {
		m.MealsEaten -= 1
	}
	_, err = s.repository.updateMealsEaten(m.ID, m.MealsEaten)
	if err != nil {
		return
	}
}

func (s *Service) GetAllMembers() ([]Member, error) {
	members, err := s.repository.getAll()
	if err != nil {
		return nil, fmt.Errorf("GetAllMembers: %w", err)
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
		ratioField := slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*Ratio Status:*\n%s", ratioStatus(uint(member.MealsEaten), uint(member.MealsCooked))), false, false)

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
		return fmt.Errorf("WeeklyUpdate GetAllMembers: %w", err)
	}

	sort.Slice(members, func(ii, jj int) bool {
		return mealsEatenToMealsCooked(uint(members[ii].MealsEaten), uint(members[ii].MealsCooked)) > mealsEatenToMealsCooked(uint(members[jj].MealsEaten), uint(members[jj].MealsCooked))
	})

	_, _, err = s.slackCfg.Client.PostMessage(s.slackCfg.Channel, weeklyUpdateBlock(members))
	if err != nil {
		return fmt.Errorf("WeeklyUpdate PostMessage: %w", err)
	}
	return nil
}
