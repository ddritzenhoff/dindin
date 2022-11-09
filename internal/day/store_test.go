package day

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Test_eatingStore_get(t *testing.T) {
	t.Run("correct values are retrieved", func(t *testing.T) {
		connection, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		if err != nil {
			t.Errorf("gorm.Open() error = %v", err)
			return
		}
		sqlDB, err := connection.DB()
		if err != nil {
			t.Errorf("connection.DB() error = %v", err)
			return
		}
		defer sqlDB.Close()
		es, err := newStore(connection)
		if err != nil {
			t.Errorf("newStore() error = %v", err)
			return
		}
		e := &Event{
			CookingDay:      1,
			CookingMonth:    2,
			CookingYear:     3,
			ChefSlackUID:    "dom123",
			MealDescription: "it will be yummy",
			SlackMessageID:  "123.123",
		}
		err = es.create(e)
		if err != nil {
			t.Errorf("eatingStore.create() error = %v", err)
			return
		}
		event, err := es.get(e.SlackMessageID)
		if err != nil {
			wantErr := false
			t.Errorf("eatingStore.get() error = %v, wantErr %v", err, wantErr)
			return
		}
		if e.CookingDay != event.CookingDay || e.CookingMonth != event.CookingMonth || e.CookingYear != event.CookingYear || e.ChefSlackUID != event.ChefSlackUID || e.MealDescription != event.MealDescription || e.SlackMessageID != event.SlackMessageID {
			t.Errorf("eatingStore.get() = %v, want %v", event, e)
		}
	})
}
