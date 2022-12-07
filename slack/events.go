package slack

import (
	"errors"
	"fmt"

	"github.com/ddritzenhoff/dindin"
	"github.com/slack-go/slack/slackevents"
)

// ReactionAddedEvent adds 1 to the Slack member's meals_eaten if a valid 'is eating' message were liked.
func (s *Service) ReactionAddedEvent(e *slackevents.ReactionAddedEvent) error {
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
	if errors.Is(err, dindin.ErrNotFound) {
		userInfo, err := s.client.GetUserInfo(slackUID)
		if err != nil {
			return fmt.Errorf("ReactionAddedEvent GetUserInfo: %w", err)
		}
		err = s.memberService.CreateMember(&dindin.Member{
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
	err = s.memberService.UpdateMember(member.ID, dindin.MemberUpdate{
		MealsEaten: &newMealsEaten,
	})
	if err != nil {
		return fmt.Errorf("ReactionAddedEvent UpdateMember: %w", err)
	}
	return nil
}

// ReactionAddedEvent adds 1 to the Slack member's meals_eaten if a valid 'is eating' message were un-liked.
func (s *Service) ReactionRemovedEvent(e *slackevents.ReactionRemovedEvent) error {
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
	if errors.Is(err, dindin.ErrNotFound) {
		userInfo, err := s.client.GetUserInfo(slackUID)
		if err != nil {
			return fmt.Errorf("ReactionRemovedEvent GetUserInfo: %w", err)
		}
		err = s.memberService.CreateMember(&dindin.Member{
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
	err = s.memberService.UpdateMember(member.ID, dindin.MemberUpdate{
		MealsEaten: &newMealsEaten,
	})
	if err != nil {
		return fmt.Errorf("ReactionRemovedEvent UpdateMember: %w", err)
	}
	return nil
}
