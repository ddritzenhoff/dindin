package slack

import (
	"errors"
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/ddritzenhoff/dinny"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

// Service represents the service to communicate with the Slack API.
type Service interface {
	PostEatingTomorrow() error
	WeeklyUpdate() error
	ReactionAddedEvent(e *slackevents.ReactionAddedEvent) error
	ReactionRemovedEvent(e *slackevents.ReactionRemovedEvent) error
}

// Config represents the configuration values to communicate with the slack API.
type Config struct {
	Channel       string
	BotSigningKey string
}

// service represents the implementation of the Service interface.
type service struct {
	client        *slack.Client
	config        *Config
	mealService   dinny.MealService
	memberService dinny.MemberService
}

// NewService returns a new instance of slack.Service.
func NewService(config *Config, mealService dinny.MealService, memberService dinny.MemberService) (*service, error) {
	client := slack.New(config.BotSigningKey)
	if client == nil {
		return nil, fmt.Errorf("NewService: couldn't generate slack client")
	}
	return &service{
		client,
		config,
		mealService,
		memberService,
	}, nil
}

// isEatingTomorrowBlock creates a 'who's eating' message to be sent into the slack channel.
func isEatingTomorrowBlock() slack.MsgOption {
	// Header Section
	headerText := slack.NewTextBlockObject("mrkdwn", "hey <!channel>, please react to this message (:thumbsup:) if you are eating tomorrow", false, false)
	headerSection := slack.NewSectionBlock(headerText, nil, nil)

	return slack.MsgOptionBlocks(
		headerSection,
	)
}

// PostEatingTomorrow sends the 'who's eating' messages into the slack channel.
func (s *service) PostEatingTomorrow() error {
	year, month, day := time.Now().AddDate(0, 0, 1).Date()
	meal, err := s.mealService.FindMealByDate(dinny.Date{Year: year, Month: month, Day: day})
	if err != nil {
		return fmt.Errorf("PostEatingTomorrow FindMealByDate: %w", err)
	}
	if meal.SlackMessageID != "" {
		return fmt.Errorf("the slack message has already been posted for tomorrow")
	}
	_, respTimestamp, err := s.client.PostMessage(s.config.Channel, isEatingTomorrowBlock())
	if err != nil {
		return fmt.Errorf("PostEatingTomorrow PostMessage: %w", err)
	}
	err = s.mealService.UpdateMeal(meal.ID, dinny.MealUpdate{SlackMessageID: &respTimestamp})
	if err != nil {
		return fmt.Errorf("PostEatingTomorrow UpdateMeal: %w", err)
	}
	return nil
}

// mealsEatenToMealsCooked calculates the meals eaten to meals cooked ratio. Returns infinity for 0 meals cooked and >0 meals eaten.
func mealsEatenToMealsCooked(mealsEaten int64, mealsCooked int64) float32 {
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

// ratioStatus calculates the the meals eaten to meals cooked ratio. Returns a string instead of a float.
func ratioStatus(mealsEaten int64, mealsCooked int64) string {
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

// weeklyUpdateBlock represents a slack message to give meals eaten to meals cooked ratio statuses.
func weeklyUpdateBlock(members []*dinny.Member) slack.MsgOption {

	var sectionBlocks []slack.Block
	// Header Section
	headerText := slack.NewTextBlockObject("mrkdwn", "dinner rotation members with the *worst* meals eaten to meals cooked ratios:", false, false)
	headerSection := slack.NewSectionBlock(headerText, nil, nil)
	sectionBlocks = append(sectionBlocks, *headerSection)

	for ii, member := range members {
		if ii > 10 {
			break
		}

		realNameField := slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*Real Name:*\n%s", member.FullName), false, false)
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

// WeeklyUpdate sends the weeklyUpdateBlock into Slack.
func (s *service) WeeklyUpdate() error {
	members, err := s.memberService.ListMembers()
	if err != nil {
		return fmt.Errorf("WeeklyUpdate ListMembers: %w", err)
	}

	sort.Slice(members, func(ii, jj int) bool {
		return mealsEatenToMealsCooked(members[ii].MealsEaten, members[ii].MealsCooked) > mealsEatenToMealsCooked(members[jj].MealsEaten, members[jj].MealsCooked)
	})

	_, _, err = s.client.PostMessage(s.config.Channel, weeklyUpdateBlock(members))
	if err != nil {
		return fmt.Errorf("WeeklyUpdate PostMessage: %w", err)
	}
	return nil
}

// ReactionAddedEvent adds 1 to the Slack member's meals_eaten if a valid 'is eating' message were liked.
func (s *service) ReactionAddedEvent(e *slackevents.ReactionAddedEvent) error {
	// Don't bother if the reaction isn't a like
	if e.Reaction != "+1" {
		return fmt.Errorf("ReactionAddedEvent +1: got %s in channel %s from user %s", e.Reaction, e.Item.Channel, e.User)
	}

	// check to see whether the reaction was on a 'who's eating tomorrow' post
	slackMessageID := e.Item.Timestamp
	meal, err := s.mealService.FindMealBySlackMessageID(slackMessageID)
	if err != nil {
		return fmt.Errorf("ReactionAddedEvent GetEatingEvent: got %s in channel %s from user %s with timestamp %s", e.Reaction, e.Item.Channel, e.User, e.Item.Timestamp)
	}

	// if the reaction was on an expired 'who's eating tomorrow' post, don't do anything
	if meal.Expired() {
		return fmt.Errorf("ReactionAddedEvent IsEatingMessageExpired: slackMessageID: %s", slackMessageID)
	}

	// add the member if he doesn't exist yet
	slackUID := e.User
	member, err := s.memberService.FindMemberBySlackUID(slackUID)
	if errors.Is(err, dinny.ErrNotFound) {
		userInfo, err := s.client.GetUserInfo(slackUID)
		if err != nil {
			return fmt.Errorf("ReactionAddedEvent GetUserInfo: %w", err)
		}
		err = s.memberService.CreateMember(&dinny.Member{
			SlackUID: slackUID,
			FullName: userInfo.RealName,
		})
		if err != nil {
			return fmt.Errorf("ReactionAddedEvent CreateMember: %w", err)
		}
	} else if err != nil {
		return fmt.Errorf("ReactionAddedEvent FindMemberBySlackUID: %w", err)
	}

	// update the member's meals eaten
	newMealsEaten := member.MealsEaten + 1
	err = s.memberService.UpdateMember(member.ID, dinny.MemberUpdate{
		MealsEaten: &newMealsEaten,
	})
	if err != nil {
		return fmt.Errorf("ReactionAddedEvent UpdateMember: %w", err)
	}
	return nil
}

// ReactionAddedEvent adds 1 to the Slack member's meals_eaten if a valid 'is eating' message were un-liked.
func (s *service) ReactionRemovedEvent(e *slackevents.ReactionRemovedEvent) error {
	// Don't bother if the reaction isn't a like
	if e.Reaction != "+1" {
		return fmt.Errorf("ReactionRemovedEvent +1: got %s in channel %s from user %s", e.Reaction, e.Item.Channel, e.User)
	}

	// check to see whether the reaction was on a 'who's eating tomorrow' post
	slackMessageID := e.Item.Timestamp
	meal, err := s.mealService.FindMealBySlackMessageID(slackMessageID)
	if err != nil {
		return fmt.Errorf("ReactionRemovedEvent GetEatingEvent: got %s in channel %s from user %s with timestamp %s", e.Reaction, e.Item.Channel, e.User, e.Item.Timestamp)
	}

	// if the reaction was on an expired 'who's eating tomorrow' post, don't do anything
	if meal.Expired() {
		return fmt.Errorf("ReactionRemovedEvent IsEatingMessageExpired: slackMessageID: %s", slackMessageID)
	}

	// add the member if he doesn't exist yet
	slackUID := e.User
	member, err := s.memberService.FindMemberBySlackUID(slackUID)
	if errors.Is(err, dinny.ErrNotFound) {
		userInfo, err := s.client.GetUserInfo(slackUID)
		if err != nil {
			return fmt.Errorf("ReactionRemovedEvent GetUserInfo: %w", err)
		}
		err = s.memberService.CreateMember(&dinny.Member{
			SlackUID: slackUID,
			FullName: userInfo.RealName,
		})
		if err != nil {
			return fmt.Errorf("ReactionRemovedEvent CreateMember: %w", err)
		}
	} else if err != nil {
		return fmt.Errorf("ReactionRemovedEvent FindMemberBySlackUID: %w", err)
	}

	// update the member's meals eaten

	newMealsEaten := member.MealsEaten + 1
	if newMealsEaten < 0 {
		newMealsEaten = 0
	}
	err = s.memberService.UpdateMember(member.ID, dinny.MemberUpdate{
		MealsEaten: &newMealsEaten,
	})
	if err != nil {
		return fmt.Errorf("ReactionRemovedEvent UpdateMember: %w", err)
	}
	return nil
}
