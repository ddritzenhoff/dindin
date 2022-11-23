package cooking

import (
	"database/sql"
	"fmt"
	"time"
)

type repository interface {
	create(d Day) (*Day, error)
	get(id int64) (*Day, error)
	getBySlackMessageID(slackMessageID string) (*Day, error)
	getByDate(year int, month time.Month, day int) (*Day, error)
	deleteBySlackMessageID(slackMessageID string) error
	update(id int64, updated Day) (*Day, error)
}

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) (*Repository, error) {
	r := Repository{db}
	err := r.Migrate()
	if err != nil {
		return nil, fmt.Errorf("NewRepository: %w", err)
	}
	return &r, nil
}

func (r *Repository) Migrate() error {
	stmt := `
    CREATE TABLE IF NOT EXISTS Events(
        id INTEGER PRIMARY KEY AUTOINCREMENT,
		created_at INTEGER NOT NULL,
        updated_at INTEGER NOT NULL,
		cooking_time INTEGER NOT NULL UNIQUE,
        chef_slack_uid TEXT NOT NULL,
		meal_description TEXT,
		slack_message_id TEXT
    );
    `

	_, err := r.db.Exec(stmt)
	return err
}

func (r *Repository) create(d Day) (*Day, error) {
	now := time.Now().Unix()
	if d.CookingTime == 0 {
		return nil, fmt.Errorf("no cooking time specified")
	}
	if d.ChefSlackUID == "" {
		return nil, fmt.Errorf("no chef slackUID specified")
	}
	stmt, err := r.db.Prepare("insert or ignore into Events(created_at, updated_at, cooking_time, chef_slack_uid, meal_description, slack_message_id) values(?,?,?,?,?,?)")
	if err != nil {
		return nil, fmt.Errorf("create prepare: %w", err)
	}
	res, err := stmt.Exec(now, now, d.CookingTime, d.ChefSlackUID, d.MealDescription, d.SlackMessageID)
	if err != nil {
		return nil, fmt.Errorf("create Exec: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("create LastInsertId: %w", err)
	}
	d.ID = id
	return &d, nil
}

func (r *Repository) get(id int64) (*Day, error) {
	var d Day
	err := r.db.QueryRow("select * from Events where id = ?", id).Scan(&d.ID, &d.CreatedAt, &d.UpdatedAt, &d.CookingTime, &d.ChefSlackUID, &d.MealDescription, &d.SlackMessageID)
	if err != nil {
		return nil, fmt.Errorf("get: %w", err)
	}
	return &d, nil
}

func (r *Repository) getBySlackMessageID(slackMessageID string) (*Day, error) {
	var d Day
	err := r.db.QueryRow("select * from Events where slack_message_id = ?", slackMessageID).Scan(&d.ID, &d.CreatedAt, &d.UpdatedAt, &d.CookingTime, &d.ChefSlackUID, &d.MealDescription, &d.SlackMessageID)
	if err != nil {
		return nil, fmt.Errorf("getBySlackMessageID: %w", err)
	}
	return &d, nil
}

func TimeFromDate(year int, month time.Month, day int) time.Time {
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}

func UnixTimeFromDate(year int, month time.Month, day int) int64 {
	return TimeFromDate(year, month, day).Unix()
}

func (r *Repository) getByDate(year int, month time.Month, day int) (*Day, error) {
	unixTime := UnixTimeFromDate(year, month, day)
	var d Day
	err := r.db.QueryRow("select * from Events where cooking_time = ?", unixTime).Scan(&d.ID, &d.CreatedAt, &d.UpdatedAt, &d.CookingTime, &d.ChefSlackUID, &d.MealDescription, &d.SlackMessageID)
	if err != nil {
		return nil, fmt.Errorf("getByDate: %w", err)
	}
	return &d, nil
}

func (r *Repository) deleteBySlackMessageID(slackMessageID string) error {
	res, err := r.db.Exec("delete * from Events where slack_message_id = ?", slackMessageID)
	if err != nil {
		return fmt.Errorf("deleteBySlackMessageID: %w", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("deleteBySlackMessageID RowsAffected: %w", err)
	}
	if rowsAffected != 1 {
		return fmt.Errorf("deleteBySlackMessageID %d rows were affected. expecting 1", rowsAffected)
	}
	return nil
}

func (r *Repository) update(id int64, updated Day) (*Day, error) {
	if id == 0 {
		return nil, fmt.Errorf("invalid updated ID")
	}
	res, err := r.db.Exec("update Events set updated_at = ?, cooking_time = ?, chef_slack_uid = ?, meal_description = ?, slack_message_id = ? where id = ?", time.Now().Unix(), updated.CookingTime, updated.ChefSlackUID, updated.MealDescription, updated.SlackMessageID, id)
	if err != nil {
		return nil, fmt.Errorf("update: %w", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("update RowsAffected: %w", err)
	}
	if rowsAffected == 0 {
		return nil, fmt.Errorf("updated failed")
	}
	return &updated, nil
}
