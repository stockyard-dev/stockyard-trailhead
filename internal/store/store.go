package store

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

type DB struct{ db *sql.DB }

type Habit struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Frequency   string `json:"frequency"` // daily, weekly
	Color       string `json:"color,omitempty"`
	Archived    bool   `json:"archived"`
	CreatedAt   string `json:"created_at"`
	Streak      int    `json:"streak"`
	BestStreak  int    `json:"best_streak"`
	TotalChecks int    `json:"total_checks"`
	CheckedToday bool  `json:"checked_today"`
}

type CheckIn struct {
	ID        string `json:"id"`
	HabitID   string `json:"habit_id"`
	Date      string `json:"date"` // YYYY-MM-DD
	Note      string `json:"note,omitempty"`
	CreatedAt string `json:"created_at"`
}

func Open(dataDir string) (*DB, error) {
	if err := os.MkdirAll(dataDir, 0755); err != nil { return nil, err }
	dsn := filepath.Join(dataDir, "trailhead.db") + "?_journal_mode=WAL&_busy_timeout=5000"
	db, err := sql.Open("sqlite", dsn)
	if err != nil { return nil, err }
	for _, q := range []string{
		`CREATE TABLE IF NOT EXISTS habits (id TEXT PRIMARY KEY, name TEXT NOT NULL, description TEXT DEFAULT '', frequency TEXT DEFAULT 'daily', color TEXT DEFAULT '#c45d2c', archived INTEGER DEFAULT 0, created_at TEXT DEFAULT (datetime('now')))`,
		`CREATE TABLE IF NOT EXISTS checkins (id TEXT PRIMARY KEY, habit_id TEXT NOT NULL REFERENCES habits(id) ON DELETE CASCADE, date TEXT NOT NULL, note TEXT DEFAULT '', created_at TEXT DEFAULT (datetime('now')))`,
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_checkins_unique ON checkins(habit_id, date)`,
		`CREATE INDEX IF NOT EXISTS idx_checkins_habit ON checkins(habit_id)`,
	} {
		if _, err := db.Exec(q); err != nil { return nil, fmt.Errorf("migrate: %w", err) }
	}
	return &DB{db: db}, nil
}

func (d *DB) Close() error { return d.db.Close() }
func genID() string { return fmt.Sprintf("%d", time.Now().UnixNano()) }
func now() string { return time.Now().UTC().Format(time.RFC3339) }
func today() string { return time.Now().Format("2006-01-02") }

func (d *DB) hydrateHabit(h *Habit) {
	d.db.QueryRow(`SELECT COUNT(*) FROM checkins WHERE habit_id=?`, h.ID).Scan(&h.TotalChecks)
	var n int
	d.db.QueryRow(`SELECT COUNT(*) FROM checkins WHERE habit_id=? AND date=?`, h.ID, today()).Scan(&n)
	h.CheckedToday = n > 0
	h.Streak = d.calcStreak(h.ID)
	h.BestStreak = d.calcBestStreak(h.ID)
}

func (d *DB) calcStreak(habitID string) int {
	date := time.Now()
	streak := 0
	for i := 0; i < 365; i++ {
		ds := date.Format("2006-01-02")
		var n int
		d.db.QueryRow(`SELECT COUNT(*) FROM checkins WHERE habit_id=? AND date=?`, habitID, ds).Scan(&n)
		if n == 0 {
			if i == 0 { date = date.AddDate(0, 0, -1); continue } // allow today not checked yet
			break
		}
		streak++
		date = date.AddDate(0, 0, -1)
	}
	return streak
}

func (d *DB) calcBestStreak(habitID string) int {
	rows, _ := d.db.Query(`SELECT date FROM checkins WHERE habit_id=? ORDER BY date ASC`, habitID)
	if rows == nil { return 0 }
	defer rows.Close()
	var dates []string
	for rows.Next() { var dt string; rows.Scan(&dt); dates = append(dates, dt) }
	if len(dates) == 0 { return 0 }
	best, cur := 1, 1
	for i := 1; i < len(dates); i++ {
		prev, _ := time.Parse("2006-01-02", dates[i-1])
		curr, _ := time.Parse("2006-01-02", dates[i])
		if curr.Sub(prev).Hours() <= 24 { cur++ } else { cur = 1 }
		if cur > best { best = cur }
	}
	return best
}

func (d *DB) CreateHabit(h *Habit) error {
	h.ID = genID(); h.CreatedAt = now()
	if h.Frequency == "" { h.Frequency = "daily" }
	if h.Color == "" { h.Color = "#c45d2c" }
	_, err := d.db.Exec(`INSERT INTO habits (id,name,description,frequency,color,created_at) VALUES (?,?,?,?,?,?)`,
		h.ID, h.Name, h.Description, h.Frequency, h.Color, h.CreatedAt)
	return err
}

func (d *DB) GetHabit(id string) *Habit {
	var h Habit; var archived int
	if err := d.db.QueryRow(`SELECT id,name,description,frequency,color,archived,created_at FROM habits WHERE id=?`, id).Scan(&h.ID, &h.Name, &h.Description, &h.Frequency, &h.Color, &archived, &h.CreatedAt); err != nil { return nil }
	h.Archived = archived == 1; d.hydrateHabit(&h); return &h
}

func (d *DB) ListHabits(includeArchived bool) []Habit {
	q := `SELECT id,name,description,frequency,color,archived,created_at FROM habits`
	if !includeArchived { q += ` WHERE archived=0` }
	q += ` ORDER BY name ASC`
	rows, _ := d.db.Query(q)
	if rows == nil { return nil }
	defer rows.Close()
	var out []Habit
	for rows.Next() {
		var h Habit; var archived int
		rows.Scan(&h.ID, &h.Name, &h.Description, &h.Frequency, &h.Color, &archived, &h.CreatedAt)
		h.Archived = archived == 1; d.hydrateHabit(&h)
		out = append(out, h)
	}
	return out
}

func (d *DB) UpdateHabit(id string, h *Habit) error {
	archived := 0; if h.Archived { archived = 1 }
	_, err := d.db.Exec(`UPDATE habits SET name=?,description=?,frequency=?,color=?,archived=? WHERE id=?`,
		h.Name, h.Description, h.Frequency, h.Color, archived, id)
	return err
}

func (d *DB) DeleteHabit(id string) error {
	d.db.Exec(`DELETE FROM checkins WHERE habit_id=?`, id)
	_, err := d.db.Exec(`DELETE FROM habits WHERE id=?`, id)
	return err
}

func (d *DB) CheckIn(habitID, date, note string) error {
	if date == "" { date = today() }
	_, err := d.db.Exec(`INSERT OR REPLACE INTO checkins (id,habit_id,date,note,created_at) VALUES (?,?,?,?,?)`,
		genID(), habitID, date, note, now())
	return err
}

func (d *DB) Uncheck(habitID, date string) error {
	if date == "" { date = today() }
	_, err := d.db.Exec(`DELETE FROM checkins WHERE habit_id=? AND date=?`, habitID, date)
	return err
}

func (d *DB) ListCheckIns(habitID string, days int) []CheckIn {
	if days <= 0 { days = 30 }
	since := time.Now().AddDate(0, 0, -days).Format("2006-01-02")
	rows, _ := d.db.Query(`SELECT id,habit_id,date,note,created_at FROM checkins WHERE habit_id=? AND date>=? ORDER BY date DESC`, habitID, since)
	if rows == nil { return nil }
	defer rows.Close()
	var out []CheckIn
	for rows.Next() {
		var c CheckIn; rows.Scan(&c.ID, &c.HabitID, &c.Date, &c.Note, &c.CreatedAt)
		out = append(out, c)
	}
	return out
}

// ── Today's view: all habits with today's check status ──

type DayView struct {
	Date   string  `json:"date"`
	Habits []Habit `json:"habits"`
	Done   int     `json:"done"`
	Total  int     `json:"total"`
}

func (d *DB) Today() DayView {
	habits := d.ListHabits(false)
	done := 0
	for _, h := range habits { if h.CheckedToday { done++ } }
	return DayView{Date: today(), Habits: habits, Done: done, Total: len(habits)}
}

type Stats struct {
	Habits    int     `json:"habits"`
	TotalChecks int   `json:"total_checks"`
	ActiveStreaks int `json:"active_streaks"`
	CompletionRate float64 `json:"completion_rate"`
}

func (d *DB) Stats() Stats {
	var s Stats
	d.db.QueryRow(`SELECT COUNT(*) FROM habits WHERE archived=0`).Scan(&s.Habits)
	d.db.QueryRow(`SELECT COUNT(*) FROM checkins`).Scan(&s.TotalChecks)
	habits := d.ListHabits(false)
	for _, h := range habits { if h.Streak > 0 { s.ActiveStreaks++ } }
	if s.Habits > 0 {
		done := 0; for _, h := range habits { if h.CheckedToday { done++ } }
		s.CompletionRate = float64(done) / float64(s.Habits) * 100
	}
	return s
}
