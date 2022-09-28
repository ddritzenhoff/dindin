package cooking

import (
	"gorm.io/gorm"
	"log"
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
		log.Println("failed creating eating entry")
		return err
	}
	return nil
}

func (es *eatingStore) get(slackMessageID string) (*Event, error) {
	var person Event
	err := es.db.First(&person, "id = ?", slackMessageID).Error
	if err != nil {
		log.Println("eating entry not found")
		return nil, err
	}
	return &person, nil
}

func (es *eatingStore) update(e *Event) error {
	err := es.db.Model(&Event{}).Where("id = ?", e.slackMessageID).Updates(e).Error
	// pretty sure that ps.db.Model(p).Updates(p).Error also works here.
	if err != nil {
		log.Println("failed to update eating entry")
		return err
	}
	return nil
}

func (es *eatingStore) delete(slackMessageID string) error {
	err := es.db.Delete(&Event{}, "id = ?", slackMessageID).Error
	if err != nil {
		log.Println("failed to delete eating entry")
		return err
	}
	return nil
}
