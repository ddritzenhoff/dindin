package member

import (
	"database/sql"
	"fmt"
	"time"
)

type repository interface {
	create(m Member) (*Member, error)
	get(id int64) (*Member, error)
	getBySlackUID(slackUID string) (*Member, error)
	getAll() ([]Member, error)
	deleteBySlackUID(slackUID string) error
	updateMealsEaten(id int64, mealsCooked int64) (*Member, error)
	updateMealsCooked(id int64, mealsCooked int64) (*Member, error)
}

type Repository struct {
	db *sql.DB
}

func newRepository(db *sql.DB) (*Repository, error) {
	r := Repository{db}
	err := r.Migrate()
	if err != nil {
		return nil, fmt.Errorf("NewRepository: %w", err)
	}
	return &r, nil
}

func (r *Repository) Migrate() error {
	stmt := `
    CREATE TABLE IF NOT EXISTS Members(
        id INTEGER PRIMARY KEY AUTOINCREMENT,
		created_at INTEGER NOT NULL
        updated_at INTEGER NOT NULL,
		slack_uid TEXT NOT NULL UNIQUE,
        first_name TEXT,
		last_name TEXT,
		meals_eaten INTEGER NOT NULL DEFAULT 0,
		meals_cooked INTEGER NOT NULL DEFAULT 0
    );
    `
	_, err := r.db.Exec(stmt)
	return err
}

func (r *Repository) create(m Member) (*Member, error) {
	now := time.Now().Unix()
	if m.SlackUID == "" {
		return nil, fmt.Errorf("no slackUID specified")
	}
	stmt, err := r.db.Prepare("insert or ignore into Members(created_at, updated_at, slack_uid, first_name, last_name, meals_eaten, meals_cooked) values(?,?,?,?,?,?,?)")
	if err != nil {
		return nil, fmt.Errorf("create prepare: %w", err)
	}
	res, err := stmt.Exec(now, now, m.SlackUID, m.FirstName, m.LastName, m.MealsEaten, m.MealsCooked)
	if err != nil {
		return nil, fmt.Errorf("create Exec: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("create LastInsertId: %w", err)
	}
	m.ID = id
	return &m, nil

}

func (r *Repository) get(id int64) (*Member, error) {
	var m Member
	err := r.db.QueryRow("select * from Members where id = ?", id).Scan(&m)
	if err != nil {
		return nil, fmt.Errorf("get: %w", err)
	}
	return &m, nil
}

func (r *Repository) getBySlackUID(slackUID string) (*Member, error) {
	var m Member
	err := r.db.QueryRow("select * from Members where slack_uid = ?", slackUID).Scan(&m)
	if err != nil {
		return nil, fmt.Errorf("getBySlackUID: %w", err)
	}
	return &m, nil
}

func (r *Repository) getAll() ([]Member, error) {
	rows, err := r.db.Query("select * from Members")
	if err != nil {
		return nil, fmt.Errorf("getAll: %w", err)
	}
	defer rows.Close()
	var all []Member
	for rows.Next() {
		var m Member
		if err := rows.Scan(&m.ID, &m.CreatedAt, &m.UpdatedAt, &m.SlackUID, &m.FirstName, &m.LastName, &m.MealsEaten, &m.MealsCooked); err != nil {
			return nil, fmt.Errorf("getAll: %w", err)
		}
		all = append(all, m)
	}
	return all, nil
}

func (r *Repository) deleteBySlackUID(slackUID string) error {
	res, err := r.db.Exec("delete * from Members where slack_uid = ?", slackUID)
	if err != nil {
		return fmt.Errorf("deleteBySlackUID: %w", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("deleteBySlackUID RowsAffected: %w", err)
	}
	if rowsAffected != 1 {
		return fmt.Errorf("deleteBySlackUID %d rows were affected. expecting 1", rowsAffected)
	}
	return nil
}

func (r *Repository) updateMealsEaten(id int64, mealsEaten int64) (*Member, error) {
	if id == 0 {
		return nil, fmt.Errorf("updateMealsEaten: invalid updated ID")
	}
	res, err := r.db.Exec("update Members set updated_at = ?, meals_eaten = ? where id = ?", time.Now().Unix(), mealsEaten, id)
	if err != nil {
		return nil, fmt.Errorf("updateMealsEaten: %w", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("updateMealsEaten: %w", err)
	}
	if rowsAffected == 0 {
		return nil, fmt.Errorf("updateMealsEaten: %d rows affected. Expecting 1", rowsAffected)
	}
	updated, err := r.get(id)
	if err != nil {
		return nil, fmt.Errorf("updateMealsEaten get: %w", err)
	}
	return updated, nil
}

func (r *Repository) updateMealsCooked(id int64, mealsCooked int64) (*Member, error) {
	if id == 0 {
		return nil, fmt.Errorf("updateMealsCooked: invalid updated ID")
	}
	res, err := r.db.Exec("update Members set updated_at = ?, meals_cooked = ? where id = ?", time.Now().Unix(), mealsCooked, id)
	if err != nil {
		return nil, fmt.Errorf("updateMealsCooked: %w", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("updateMealsCooked: %w", err)
	}
	if rowsAffected == 0 {
		return nil, fmt.Errorf("updateMealsCooked: %d rows affected. Expecting 1", rowsAffected)
	}
	updated, err := r.get(id)
	if err != nil {
		return nil, fmt.Errorf("updateMealsCooked get: %w", err)
	}
	return updated, nil
}
