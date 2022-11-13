package cooking

import (
	"database/sql"
	"fmt"
	"time"
)

type repository interface {
	create(e Day) (*Day, error)
	get(id int64) (*Day, error)
	getBySlackMessageID(slackMessageID string) (*Day, error)
	getByDate(year int, month int, day int) (*Day, error)
	deleteBySlackMessageID(slackMessageID string) error
	update(id int64, updated Day) (*Day, error)
}

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) (*Repository, error) {
	dr := Repository{db}
	err := dr.Migrate()
	if err != nil {
		return nil, fmt.Errorf("NewRepository: %w", err)
	}
	return &dr, nil
}

func (dr *Repository) Migrate() error {
	stmt := `
    CREATE TABLE IF NOT EXISTS Events(
        id INTEGER PRIMARY KEY AUTOINCREMENT,
		created_at INTEGER NOT NULL
        updated_at INTEGER NOT NULL,
		cooking_time INTEGER NOT NULL UNIQUE,
        chef_slack_uid TEXT NOT NULL,
		meal_description TEXT,
		slack_message_id TEXT
    );
    `

	_, err := dr.db.Exec(stmt)
	return err
}

func (dr *Repository) create(e Day) (*Day, error) {
	now := time.Now().Unix()
	if e.CookingTime == 0 {
		return nil, fmt.Errorf("no cooking time specified")
	}
	if e.ChefSlackUID == "" {
		return nil, fmt.Errorf("no chef slackUID specified")
	}
	stmt, err := dr.db.Prepare("insert or ignore into Events(created_at, updated_at, cooking_time, chef_slack_uid, meal_description, slack_message_id) values(?,?,?,?,?,?)")
	if err != nil {
		return nil, fmt.Errorf("create prepare: %w", err)
	}
	res, err := stmt.Exec(now, now, e.CookingTime, e.ChefSlackUID, e.MealDescription, e.SlackMessageID)
	if err != nil {
		return nil, fmt.Errorf("create Exec: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("create LastInsertId: %w", err)
	}
	e.ID = id
	return &e, nil
}

func (dr *Repository) get(id int64) (*Day, error) {
	var e Day
	err := dr.db.QueryRow("select * from Events where id = ?", id).Scan(&e)
	if err != nil {
		return nil, fmt.Errorf("get: %w", err)
	}
	return &e, nil
}

func (dr *Repository) getBySlackMessageID(slackMessageID string) (*Day, error) {
	var e Day
	err := dr.db.QueryRow("select * from Events where slack_message_id = ?", slackMessageID).Scan(&e)
	if err != nil {
		return nil, fmt.Errorf("getBySlackMessageID: %w", err)
	}
	return &e, nil
}

func TimeFromDate(year int, month int, day int) time.Time {
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}

func UnixTimeFromDate(year int, month int, day int) int64 {
	return TimeFromDate(year, month, day).Unix()
}

func (dr *Repository) getByDate(year int, month int, day int) (*Day, error) {
	unixTime := UnixTimeFromDate(year, month, day)
	var e Day
	err := dr.db.QueryRow("select * from Events where cooking_time = ?", unixTime).Scan(&e)
	if err != nil {
		return nil, fmt.Errorf("getByDate: %w", err)
	}
	return &e, nil
}

func (dr *Repository) deleteBySlackMessageID(slackMessageID string) error {
	res, err := dr.db.Exec("delete * from Events where slack_message_id = ?", slackMessageID)
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

func (dr *Repository) update(id int64, updated Day) (*Day, error) {
	if id == 0 {
		return nil, fmt.Errorf("invalid updated ID")
	}
	res, err := dr.db.Exec("update Events set updated_at = ?, cooking_time = ?, chef_slack_uid = ?, meal_description = ?, slack_message_id = ? where id = ?", time.Now().Unix(), updated.CookingTime, updated.ChefSlackUID, updated.MealDescription, updated.SlackMessageID, id)
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
