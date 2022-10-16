package member

import (
	"fmt"

	"gorm.io/gorm"
)

type store interface {
	create(p *Member) error
	get(slackUID string) (*Member, error)
	delete(slackUID string) error
	update(p *Member) error
	getAll() ([]Member, error)
}

type memberStore struct {
	db *gorm.DB
}

func newStore(db *gorm.DB) (*memberStore, error) {
	return &memberStore{db: db}, nil
}

func (ps *memberStore) create(p *Member) error {
	err := ps.db.Create(p).Error
	if err != nil {
		return fmt.Errorf("Failed creating member: %w", err)
	}
	return nil
}

func (ps *memberStore) get(slackUID string) (*Member, error) {
	var member Member
	err := ps.db.First(&member, "id = ?", slackUID).Error
	if err != nil {
		return nil, fmt.Errorf("Member not found: %w", err)
	}
	return &member, nil
}

func (ps *memberStore) getAll() ([]Member, error) {
	var members []Member
	err := ps.db.Find(&members).Error
	if err != nil {
		return nil, fmt.Errorf("Couldn't retrieve all of the members: %w", err)
	}
	return members, nil
}

func (ps *memberStore) update(p *Member) error {
	err := ps.db.Model(&Member{}).Where("id = ?", p.SlackUID).Updates(p).Error
	// pretty sure that ps.db.Model(p).Updates(p).Error also works here.
	if err != nil {
		return fmt.Errorf("Failed to update member: %w", err)
	}
	return nil
}

func (ps *memberStore) delete(slackUID string) error {
	err := ps.db.Delete(&Member{}, "id = ?", slackUID).Error
	if err != nil {
		return fmt.Errorf("Failed to member: %w", err)
	}
	return nil
}
