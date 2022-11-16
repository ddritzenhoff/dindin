package cooking

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/ddritzenhoff/dindin/internal/configs"
	"github.com/ddritzenhoff/dindin/internal/http/rpc/pb"
	"github.com/slack-go/slack"
)

const DAY_LENGTH_SECONDS = time.Hour * 24

type Day struct {
	ID              int64
	CreatedAt       int64
	UpdatedAt       int64
	CookingTime     int64
	ChefSlackUID    string
	MealDescription string
	SlackMessageID  string
}

func (d *Day) ToTime(unixTime int64) time.Time {
	return time.Unix(unixTime, 0)
}

// IsEatingMessageExpired returns true if it's day(s) after the meal should have been cooked
func (d *Day) IsEatingMessageExpired() bool {
	expiredTime := d.ToTime(d.CookingTime).Add(DAY_LENGTH_SECONDS)
	return time.Now().After(expiredTime)
}

// Service struct holds all the dependencies required for the CookingEvent struct and exposes all services
// provided by this package as its methods
type Service struct {
	repository repository
	slackCfg   *configs.SlackConfig
}

func NewService(db *sql.DB, slackCfg *configs.SlackConfig) (*Service, error) {
	pStore, err := NewRepository(db)
	if err != nil {
		return nil, err
	}
	return &Service{repository: pStore, slackCfg: slackCfg}, nil
}

func (s *Service) GetEatingEvent(slackMessageID string) (eatingEvent *Day, exists bool) {
	event, err := s.repository.getBySlackMessageID(slackMessageID)
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

func (s *Service) PostEatingTomorrow() error {
	year, month, day := time.Now().AddDate(0, 0, 1).Date()
	e, err := s.repository.getByDate(year, int(month), day)
	if err != nil {
		return err
	}
	if e.SlackMessageID != "" {
		return fmt.Errorf("the slack message has already been posted for tomorrow")
	}
	_, respTimestamp, err := s.slackCfg.Client.PostMessage(s.slackCfg.Channel, isEatingTomorrowBlock())
	if err != nil {
		return err
	}
	e.SlackMessageID = respTimestamp
	_, err = s.repository.update(e.ID, *e)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) AssignCooks(cookings []*pb.Cooking) error {
	for _, cooking := range cookings {
		d, err := s.repository.getByDate(int(cooking.Year), int(cooking.Month), int(cooking.Day))
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				_, err = s.repository.create(Day{
					ChefSlackUID: cooking.SlackUID,
					CookingTime:  UnixTimeFromDate(int(cooking.Year), int(cooking.Month), int(cooking.Day)),
				})
				if err != nil {
					return fmt.Errorf("AssignCooks create: %w", err)
				}
				return err
			}
		}
		d.ChefSlackUID = cooking.SlackUID
		s.repository.update(d.ID, *d)
	}
	return nil
}
