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

func (ms *memberStore) create(p *Member) error {
	err := ms.db.Create(p).Error
	if err != nil {
		return fmt.Errorf("create member: %w", err)
	}
	return nil
}

func (ms *memberStore) get(slackUID string) (*Member, error) {
	var member Member
	err := ms.db.First(&member, "id = ?", slackUID).Error
	if err != nil {
		return nil, fmt.Errorf("get member: %w", err)
	}
	return &member, nil
}

func (ms *memberStore) getAll() ([]Member, error) {
	var members []Member
	err := ms.db.Find(&members).Error
	if err != nil {
		return nil, fmt.Errorf("getAll members: %w", err)
	}
	return members, nil
}

func (ms *memberStore) update(p *Member) error {
	err := ms.db.Model(&Member{}).Where("id = ?", p.SlackUID).Updates(p).Error
	// pretty sure that ps.db.Model(p).Updates(p).Error also works here.
	if err != nil {
		return fmt.Errorf("update member: %w", err)
	}
	return nil
}

func (ms *memberStore) delete(slackUID string) error {
	err := ms.db.Delete(&Member{}, "id = ?", slackUID).Error
	if err != nil {
		return fmt.Errorf("delete member: %w", err)
	}
	return nil
}
