package person

import (
	"gorm.io/gorm"
	"log"
)

type store interface {
	create(p *Person) error
	get(slackUID string) (*Person, error)
	delete(slackUID string) error
	update(p *Person) error
}

type personStore struct {
	db *gorm.DB
}

func newStore(db *gorm.DB) (*personStore, error) {
	return &personStore{db: db}, nil
}

func (ps *personStore) create(p *Person) error {
	err := ps.db.Create(p).Error
	if err != nil {
		log.Println("failed creating person")
		return err
	}
	return nil
}

func (ps *personStore) get(slackUID string) (*Person, error) {
	var person Person
	err := ps.db.First(&person, "id = ?", slackUID).Error
	if err != nil {
		log.Println("person not found")
		return nil, err
	}
	return &person, nil
}

func (ps *personStore) update(p *Person) error {
	err := ps.db.Model(&Person{}).Where("id = ?", p.SlackUID).Updates(p).Error
	// pretty sure that ps.db.Model(p).Updates(p).Error also works here.
	if err != nil {
		log.Println("failed to update person")
		return err
	}
	return nil
}

func (ps *personStore) delete(slackUID string) error {
	err := ps.db.Delete(&Person{}, "id = ?", slackUID).Error
	if err != nil {
		log.Println("failed to delete person")
		return err
	}
	return nil
}
