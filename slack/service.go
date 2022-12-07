package slack

import (
	"fmt"
	"log"
	"math"
	"sort"
	"time"

	"github.com/ddritzenhoff/dindin"
	"github.com/slack-go/slack"
)

// Config represents the configuration values to communicate with the slack API.
type Config struct {
	Channel       string
	BotSigningKey string
}

// Service represents the service to communciate with the slack API.
type Service struct {
	client        *slack.Client
	config        *Config
	mealService   dindin.MealService
	memberService dindin.MemberService
}

// NewService returns a new instance of slack.Service.
func NewService(config *Config, mealService dindin.MealService, memberService dindin.MemberService) *Service {
	client := slack.New(config.BotSigningKey)
	if client == nil {
		log.Fatal("NewService slack.New")
	}
	return &Service{
		client,
		config,
		mealService,
		memberService,
	}
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
func (s *Service) PostEatingTomorrow() error {
	year, month, day := time.Now().AddDate(0, 0, 1).Date()
	meal, err := s.mealService.FindMealByDate(dindin.Date{Year: year, Month: month, Day: day})
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
	err = s.mealService.UpdateMeal(meal.ID, dindin.MealUpdate{SlackMessageID: &respTimestamp})
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
func weeklyUpdateBlock(members []*dindin.Member) slack.MsgOption {

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
func (s *Service) WeeklyUpdate() error {
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
