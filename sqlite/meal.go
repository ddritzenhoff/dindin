package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/ddritzenhoff/dinny"
	"github.com/ddritzenhoff/dinny/sqlite/gen"
)

// Ensure service implements interface.
var _ dinny.MealService = (*MealService)(nil)

// MemberService represents a service for managing members.
type MealService struct {
	query *gen.Queries
	db    *sql.DB
}

// NewMemberService returns a new instance of MemberService.
func NewMealService(query *gen.Queries, db *sql.DB) *MealService {
	return &MealService{query, db}
}

// toDindinMeal converts a gen.Meal to a dinny.Meal
func toDindinMeal(m gen.Meal) *dinny.Meal {
	var desc string
	var smid string
	if m.Description.Valid {
		desc = m.Description.String
	} else {
		desc = ""
	}
	if m.SlackMessageID.Valid {
		smid = m.SlackMessageID.String
	} else {
		smid = ""
	}

	return &dinny.Meal{
		ID:           m.ID,
		CookSlackUID: m.CookSlackUid,
		Date: dinny.Date{
			Year:  int(m.Year),
			Month: time.Month(m.Month),
			Day:   int(m.Day),
		},
		Description:    desc,
		SlackMessageID: smid,
	}
}

// FindMealByID retrieves a meal by ID.
// Returns ErrNotFound if meal does not exist.
func (ms *MealService) FindMealByID(id int64) (*dinny.Meal, error) {
	m, err := ms.query.FindMealByID(context.Background(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, dinny.ErrNotFound
		} else {
			return nil, fmt.Errorf("FindMealByID: %w", err)
		}
	}
	return toDindinMeal(m), nil
}

// FindMealByDate retrieves a meal by Date.
// Returns ErrNotFound if meal does not exist.
func (ms *MealService) FindMealByDate(date dinny.Date) (*dinny.Meal, error) {
	params := gen.FindMealByDateParams{
		Year:  int64(date.Year),
		Month: int64(date.Month),
		Day:   int64(date.Day),
	}
	m, err := ms.query.FindMealByDate(context.Background(), params)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, dinny.ErrNotFound
		} else {
			return nil, fmt.Errorf("FindMealByDate: %w", err)
		}
	}
	return toDindinMeal(m), nil
}

// FindMealBySlackMessageID retrieves a meal by SlackMessageID.
// Returns ErrNotFound if meal does not exist.
func (ms *MealService) FindMealBySlackMessageID(slackMessageID string) (*dinny.Meal, error) {
	param := sql.NullString{
		String: slackMessageID,
		Valid:  true,
	}
	m, err := ms.query.FindMealBySlackMessageID(context.Background(), param)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, dinny.ErrNotFound
		} else {
			return nil, fmt.Errorf("FindMealBySlackMessageID: %w", err)
		}
	}
	return toDindinMeal(m), nil
}

// CreateMeal creates a new meal.
func (ms *MealService) CreateMeal(m *dinny.Meal) error {
	arg := gen.CreateMealParams{
		CookSlackUid: m.CookSlackUID,
		Year:         int64(m.Date.Year),
		Month:        int64(m.Date.Month),
		Day:          int64(m.Date.Day),
	}
	_, err := ms.query.CreateMeal(context.Background(), arg)
	if err != nil {
		return fmt.Errorf("CreateMeal: %w", err)
	}
	return nil
}

// UpdateMeal updates a meal object.
func (ms *MealService) UpdateMeal(id int64, upd dinny.MealUpdate) error {
	tx, err := ms.db.Begin()
	if err != nil {
		return fmt.Errorf("UpdateMember db.Begin: %w", err)
	}
	defer tx.Rollback()
	qtx := ms.query.WithTx(tx)
	if upd.Description != nil {
		s := sql.NullString{
			String: *upd.Description,
			Valid:  true,
		}
		params := gen.UpdateMealDescriptionParams{ID: id, Description: s}
		err := qtx.UpdateMealDescription(context.Background(), params)
		if err != nil {
			return fmt.Errorf("UpdateMeal UpdateMealDescription: %w", err)
		}
	}

	if upd.ChefSlackUID != nil {
		params := gen.UpdateMealSlackUIDParams{
			ID:           id,
			CookSlackUid: *upd.ChefSlackUID,
		}
		err := qtx.UpdateMealSlackUID(context.Background(), params)
		if err != nil {
			return fmt.Errorf("UpdateMeal UpdateMealSlackUID: %w", err)
		}
	}

	if upd.SlackMessageID != nil {
		s := sql.NullString{
			String: *upd.SlackMessageID,
			Valid:  true,
		}
		params := gen.UpdateMealSlackMessageIDParams{ID: id, SlackMessageID: s}
		err := qtx.UpdateMealSlackMessageID(context.Background(), params)
		if err != nil {
			return fmt.Errorf("UpdateMeal UpdateMealSlackMessageID: %w", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("UpdateMeal tx.Commit: %w", err)
	}
	return nil
}

// DeleteMeal permanently deletes a meal.
func (ms *MealService) DeleteMeal(id int64) error {
	err := ms.query.DeleteMeal(context.Background(), id)
	if err != nil {
		return fmt.Errorf("DeleteMeal: %w", err)
	}
	return nil
}
