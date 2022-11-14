package cooking

import (
	"database/sql"
	"errors"
	"log"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func TestRepository_create(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	r, err := NewRepository(db)
	if err != nil {
		log.Fatal(err)
	}
	d1 := Day{
		CookingTime:     time.Now().Unix(),
		ChefSlackUID:    "dom",
		MealDescription: "meal",
		SlackMessageID:  "slack",
	}
	d2, err := r.create(d1)
	if err != nil {
		log.Fatal(err)
	}
	if d1.ChefSlackUID != d2.ChefSlackUID || d1.CookingTime != d2.CookingTime || d1.MealDescription != d2.MealDescription || d1.SlackMessageID != d2.SlackMessageID || d2.ID != 1 {
		t.Errorf("Repository.create() = %v, want %v", d2, d1)
	}
}

func TestRepository_get(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	r, err := NewRepository(db)
	if err != nil {
		log.Fatal(err)
	}
	d1, err := r.create(Day{
		CookingTime:     time.Now().Unix(),
		ChefSlackUID:    "dom",
		MealDescription: "meal",
		SlackMessageID:  "slack",
	})
	if err != nil {
		log.Fatal(err)
	}
	d2, err := r.get(d1.ID)
	if err != nil {
		log.Fatal(err)
	}
	if d1.ChefSlackUID != d2.ChefSlackUID || d1.CookingTime != d2.CookingTime || d1.MealDescription != d2.MealDescription || d1.SlackMessageID != d2.SlackMessageID || d2.ID != 1 {
		t.Errorf("Repository.create() = %v, want %v", d2, d1)
	}
	idNotInDB := 341
	_, err = r.get(int64(idNotInDB))
	if !errors.Is(err, sql.ErrNoRows) {
		t.Errorf("Repository.create () = %v, want %v", err, sql.ErrNoRows)
	}
}
