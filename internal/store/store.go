package store

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

type DB struct { db *sql.DB }

type Habit struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Frequency    string   `json:"frequency"`
	Streak       int      `json:"streak"`
	Target       int      `json:"target"`
	Status       string   `json:"status"`
	CreatedAt    string   `json:"created_at"`
}

func Open(dataDir string) (*DB, error) {
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, err
	}
	dsn := filepath.Join(dataDir, "trailhead.db") + "?_journal_mode=WAL&_busy_timeout=5000"
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS habits (
			id TEXT PRIMARY KEY,\n\t\t\tname TEXT DEFAULT '',\n\t\t\tfrequency TEXT DEFAULT 'daily',\n\t\t\tstreak INTEGER DEFAULT 0,\n\t\t\ttarget INTEGER DEFAULT 0,\n\t\t\tstatus TEXT DEFAULT 'active',
			created_at TEXT DEFAULT (datetime('now'))
		)`)
	if err != nil {
		return nil, fmt.Errorf("migrate: %w", err)
	}
	return &DB{db: db}, nil
}

func (d *DB) Close() error { return d.db.Close() }

func genID() string { return fmt.Sprintf("%d", time.Now().UnixNano()) }

func (d *DB) Create(e *Habit) error {
	e.ID = genID()
	e.CreatedAt = time.Now().UTC().Format(time.RFC3339)
	_, err := d.db.Exec(`INSERT INTO habits (id, name, frequency, streak, target, status, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		e.ID, e.Name, e.Frequency, e.Streak, e.Target, e.Status, e.CreatedAt)
	return err
}

func (d *DB) Get(id string) *Habit {
	row := d.db.QueryRow(`SELECT id, name, frequency, streak, target, status, created_at FROM habits WHERE id=?`, id)
	var e Habit
	if err := row.Scan(&e.ID, &e.Name, &e.Frequency, &e.Streak, &e.Target, &e.Status, &e.CreatedAt); err != nil {
		return nil
	}
	return &e
}

func (d *DB) List() []Habit {
	rows, err := d.db.Query(`SELECT id, name, frequency, streak, target, status, created_at FROM habits ORDER BY created_at DESC`)
	if err != nil {
		return nil
	}
	defer rows.Close()
	var result []Habit
	for rows.Next() {
		var e Habit
		if err := rows.Scan(&e.ID, &e.Name, &e.Frequency, &e.Streak, &e.Target, &e.Status, &e.CreatedAt); err != nil {
			continue
		}
		result = append(result, e)
	}
	return result
}

func (d *DB) Delete(id string) error {
	_, err := d.db.Exec(`DELETE FROM habits WHERE id=?`, id)
	return err
}

func (d *DB) Count() int {
	var n int
	d.db.QueryRow(`SELECT COUNT(*) FROM habits`).Scan(&n)
	return n
}
