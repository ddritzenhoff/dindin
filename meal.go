package dinny

import "time"

// Date represents the year, month, and day of the meal.
type Date struct {
	Year  int        `json:"year"`
	Month time.Month `json:"month"`
	Day   int        `json:"day"`
}

// Meal represents a meal in dinner rotation.
type Meal struct {
	ID             int64  `json:"id"`
	CookSlackUID   string `json:"cookSlackUID"`
	Date           Date   `json:"date"`
	Description    string `json:"description"`
	SlackMessageID string `json:"slackMessageID"`
}

// Expired determines if a meal is expired if it's after the date the meal was supposed to occur.
func (m *Meal) Expired() bool {
	year, month, day := time.Now().Date()
	return year > m.Date.Year || month > m.Date.Month || day > m.Date.Day
}

// MealService represents a service for managing meals.
type MealService interface {
	// FindMealByID retrieves a meal by ID.
	// Returns ErrNotFound if meal does not exist.
	FindMealByID(id int64) (*Meal, error)

	// FindMealByDate retrieves a meal by Date.
	// Returns ErrNotFound if meal does not exist.
	FindMealByDate(date Date) (*Meal, error)

	// FindMealBySlackMessageID retrieves a meal by SlackMessageID.
	// Returns ErrNotFound if meal does not exist.
	FindMealBySlackMessageID(slackMessageID string) (*Meal, error)

	// CreateMeal creates a new meal.
	CreateMeal(m *Meal) error

	// UpdateMeal updates a meal object.
	UpdateMeal(id int64, upd MealUpdate) error

	// DeleteMeal permanently deletes a meal.
	DeleteMeal(id int64) error
}

// MealUpdate represents a set of fields to be updated via UpdateMeal().
type MealUpdate struct {
	ChefSlackUID   *string
	Description    *string
	SlackMessageID *string
}
