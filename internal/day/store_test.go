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

func Test_eatingStore_getByDate(t *testing.T) {
	t.Run("event is retrieved by date if exists", func(t *testing.T) {
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
			CookingYear:     3,
			CookingMonth:    2,
			CookingDay:      1,
			ChefSlackUID:    "a",
			MealDescription: "b",
			SlackMessageID:  "c",
		}
		err = es.create(e)
		if err != nil {
			t.Errorf("eatingStore.create() error = %v", err)
			return
		}
		event, err := es.getByDate(3, 2, 1)
		if err != nil {
			wantErr := false
			t.Errorf("eatingStore.getByDate() error = %v, wantErr %v", err, wantErr)
			return
		}
		if e.CookingDay != event.CookingDay || e.CookingMonth != event.CookingMonth || e.CookingYear != event.CookingYear || e.ChefSlackUID != event.ChefSlackUID || e.MealDescription != event.MealDescription || e.SlackMessageID != event.SlackMessageID {
			t.Errorf("eatingStore.getByDate() = %v, want %v", event, e)
		}
	})
}

func Test_eatingStore_getByDateOrCreate(t *testing.T) {
	t.Run("event is created if not there", func(t *testing.T) {
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
		e, err := es.getByDateOrCreate(1, 2, 3)
		if err != nil {
			t.Errorf("eatingStore.getByDateOrCreate() error = %v", err)
			return
		}
		event, err := es.getByDate(1, 2, 3)
		if err != nil {
			wantErr := false
			t.Errorf("eatingStore.getByDate() error = %v, wantErr %v", err, wantErr)
			return
		}
		if e.CookingDay != event.CookingDay || e.CookingMonth != event.CookingMonth || e.CookingYear != event.CookingYear || e.ChefSlackUID != event.ChefSlackUID || e.MealDescription != event.MealDescription || e.SlackMessageID != event.SlackMessageID {
			t.Errorf("eatingStore.get() = %v, want %v", event, e)
		}
	})
}

func Test_eatingStore_update(t *testing.T) {
	t.Run("event is properly updated", func(t *testing.T) {
		connection, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
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
			CookingYear:     3,
			CookingMonth:    2,
			CookingDay:      1,
			ChefSlackUID:    "a",
			MealDescription: "b",
			SlackMessageID:  "c",
		}
		err = es.create(e)
		if err != nil {
			t.Errorf("eatingStore.create() error = %v", err)
			return
		}
		e2, err := es.get(e.SlackMessageID)
		if err != nil {
			t.Errorf("eatingStore.get() error = %v", err)
			return
		}
		e2.CookingYear = 6
		e2.CookingMonth = 5
		e2.CookingDay = 4
		e2.ChefSlackUID = "dom"
		e2.MealDescription = "something else"

		err = es.update(e2)
		if err != nil {
			wantErr := false
			t.Errorf("eatingStore.update() error = %v, wantErr %v", err, wantErr)
			return
		}
		e3, err := es.get(e.SlackMessageID)
		if err != nil {
			t.Errorf("eatingStore.get() error = %v", err)
			return
		}
		if e3.CookingDay != e2.CookingDay || e3.CookingMonth != e2.CookingMonth || e3.CookingYear != e2.CookingYear || e3.ChefSlackUID != e2.ChefSlackUID || e3.MealDescription != e2.MealDescription || e3.SlackMessageID != e2.SlackMessageID {
			t.Errorf("eatingStore.get() = %v, want %v", e3, e2)
		}
	})
}
