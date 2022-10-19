package day

import (
	"fmt"

	"gorm.io/gorm"
)

type store interface {
	create(p *Event) error
	get(slackMessageID string) (*Event, error)
	delete(slackMessageID string) error
	update(p *Event) error
}

type eatingStore struct {
	db *gorm.DB
}

func newStore(db *gorm.DB) (*eatingStore, error) {
	return &eatingStore{db: db}, nil
}

func (es *eatingStore) create(e *Event) error {
	err := es.db.Create(e).Error
	if err != nil {
		return fmt.Errorf("create event: %w", err)
	}
	return nil
}

func (es *eatingStore) get(slackMessageID string) (*Event, error) {
	var person Event
	err := es.db.First(&person, slackMessageID).Error
	if err != nil {
		return nil, fmt.Errorf("get event: %w", err)
	}
	return &person, nil
}

func (es *eatingStore) update(e *Event) error {
	err := es.db.Model(&Event{}).Updates(e).Error
	if err != nil {
		return fmt.Errorf("update event: %w", err)
	}
	return nil
}

func (es *eatingStore) delete(slackMessageID string) error {
	err := es.db.Model(&Event{}).Delete(slackMessageID).Error
	if err != nil {
		return fmt.Errorf("delete event: %w", err)
	}
	return nil
}
