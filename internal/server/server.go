package server

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/stockyard-dev/stockyard-trailhead/internal/store"
)

const resourceName = "habits"

type Server struct {
	db      *store.DB
	mux     *http.ServeMux
	limits  Limits
	dataDir string
	pCfg    map[string]json.RawMessage
}

func New(db *store.DB, limits Limits, dataDir string) *Server {
	s := &Server{
		db:      db,
		mux:     http.NewServeMux(),
		limits:  limits,
		dataDir: dataDir,
	}
	s.loadPersonalConfig()

	// Habits CRUD
	s.mux.HandleFunc("GET /api/habits", s.listHabits)
	s.mux.HandleFunc("POST /api/habits", s.createHabit)
	s.mux.HandleFunc("GET /api/habits/{id}", s.getHabit)
	s.mux.HandleFunc("PUT /api/habits/{id}", s.updateHabit)
	s.mux.HandleFunc("DELETE /api/habits/{id}", s.deleteHabit)

	// Check-in actions
	s.mux.HandleFunc("POST /api/habits/{id}/check", s.checkIn)
	s.mux.HandleFunc("POST /api/habits/{id}/uncheck", s.uncheck)
	s.mux.HandleFunc("GET /api/habits/{id}/history", s.history)

	// Views and stats
	s.mux.HandleFunc("GET /api/today", s.today)
	s.mux.HandleFunc("GET /api/stats", s.stats)
	s.mux.HandleFunc("GET /api/health", s.health)

	// Personalization
	s.mux.HandleFunc("GET /api/config", s.configHandler)

	// Extras
	s.mux.HandleFunc("GET /api/extras/{resource}", s.listExtras)
	s.mux.HandleFunc("GET /api/extras/{resource}/{id}", s.getExtras)
	s.mux.HandleFunc("PUT /api/extras/{resource}/{id}", s.putExtras)

	// Tier
	s.mux.HandleFunc("GET /api/tier", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, map[string]any{
			"tier":        s.limits.Tier,
			"upgrade_url": "https://stockyard.dev/trailhead/",
		})
	})

	// Dashboard
	s.mux.HandleFunc("GET /ui", s.dashboard)
	s.mux.HandleFunc("GET /ui/", s.dashboard)
	s.mux.HandleFunc("GET /", s.root)

	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

// ─── helpers ──────────────────────────────────────────────────────

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(v)
}

func writeErr(w http.ResponseWriter, code int, msg string) {
	writeJSON(w, code, map[string]string{"error": msg})
}

func orEmpty[T any](s []T) []T {
	if s == nil {
		return []T{}
	}
	return s
}

func (s *Server) root(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	http.Redirect(w, r, "/ui", http.StatusFound)
}

// ─── personalization ──────────────────────────────────────────────

func (s *Server) loadPersonalConfig() {
	path := filepath.Join(s.dataDir, "config.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}
	var cfg map[string]json.RawMessage
	if err := json.Unmarshal(data, &cfg); err != nil {
		log.Printf("trailhead: warning: could not parse config.json: %v", err)
		return
	}
	s.pCfg = cfg
	log.Printf("trailhead: loaded personalization from %s", path)
}

func (s *Server) configHandler(w http.ResponseWriter, r *http.Request) {
	if s.pCfg == nil {
		writeJSON(w, 200, map[string]any{})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s.pCfg)
}

// ─── extras ───────────────────────────────────────────────────────

func (s *Server) listExtras(w http.ResponseWriter, r *http.Request) {
	resource := r.PathValue("resource")
	all := s.db.AllExtras(resource)
	out := make(map[string]json.RawMessage, len(all))
	for id, data := range all {
		out[id] = json.RawMessage(data)
	}
	writeJSON(w, 200, out)
}

func (s *Server) getExtras(w http.ResponseWriter, r *http.Request) {
	resource := r.PathValue("resource")
	id := r.PathValue("id")
	data := s.db.GetExtras(resource, id)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(data))
}

func (s *Server) putExtras(w http.ResponseWriter, r *http.Request) {
	resource := r.PathValue("resource")
	id := r.PathValue("id")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeErr(w, 400, "read body")
		return
	}
	var probe map[string]any
	if err := json.Unmarshal(body, &probe); err != nil {
		writeErr(w, 400, "invalid json")
		return
	}
	if err := s.db.SetExtras(resource, id, string(body)); err != nil {
		writeErr(w, 500, "save failed")
		return
	}
	writeJSON(w, 200, map[string]string{"ok": "saved"})
}

// ─── habits ───────────────────────────────────────────────────────

func (s *Server) listHabits(w http.ResponseWriter, r *http.Request) {
	incArch := r.URL.Query().Get("archived") == "true"
	writeJSON(w, 200, map[string]any{"habits": orEmpty(s.db.ListHabits(incArch))})
}

func (s *Server) createHabit(w http.ResponseWriter, r *http.Request) {
	if s.limits.MaxItems > 0 && len(s.db.ListHabits(false)) >= s.limits.MaxItems {
		writeErr(w, 402, "Free tier limit reached. Upgrade at https://stockyard.dev/trailhead/")
		return
	}
	var h store.Habit
	if err := json.NewDecoder(r.Body).Decode(&h); err != nil {
		writeErr(w, 400, "invalid json")
		return
	}
	if h.Name == "" {
		writeErr(w, 400, "name required")
		return
	}
	if err := s.db.CreateHabit(&h); err != nil {
		writeErr(w, 500, "create failed")
		return
	}
	writeJSON(w, 201, s.db.GetHabit(h.ID))
}

func (s *Server) getHabit(w http.ResponseWriter, r *http.Request) {
	h := s.db.GetHabit(r.PathValue("id"))
	if h == nil {
		writeErr(w, 404, "not found")
		return
	}
	writeJSON(w, 200, h)
}

// updateHabit accepts a full or partial habit. All empty string fields are
// preserved from the existing record. Archived has special handling: a
// PATCH-style update can't distinguish "not sent" from "false", so we use
// a separate endpoint for archive/unarchive instead.
//
// Same partial-update preservation pattern as the other tools.
func (s *Server) updateHabit(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	existing := s.db.GetHabit(id)
	if existing == nil {
		writeErr(w, 404, "not found")
		return
	}

	// Decode into a flexible map first so we can detect which fields were
	// actually sent (vs zero values from omission).
	var raw map[string]json.RawMessage
	if err := json.NewDecoder(r.Body).Decode(&raw); err != nil {
		writeErr(w, 400, "invalid json")
		return
	}

	patch := store.Habit{
		Name:        existing.Name,
		Description: existing.Description,
		Frequency:   existing.Frequency,
		Color:       existing.Color,
		Archived:    existing.Archived,
	}

	if v, ok := raw["name"]; ok {
		var s string
		json.Unmarshal(v, &s)
		if s != "" {
			patch.Name = s
		}
	}
	if v, ok := raw["description"]; ok {
		var s string
		json.Unmarshal(v, &s)
		patch.Description = s
	}
	if v, ok := raw["frequency"]; ok {
		var s string
		json.Unmarshal(v, &s)
		if s != "" {
			patch.Frequency = s
		}
	}
	if v, ok := raw["color"]; ok {
		var s string
		json.Unmarshal(v, &s)
		if s != "" {
			patch.Color = s
		}
	}
	if v, ok := raw["archived"]; ok {
		var b bool
		json.Unmarshal(v, &b)
		patch.Archived = b
	}

	if err := s.db.UpdateHabit(id, &patch); err != nil {
		writeErr(w, 500, "update failed")
		return
	}
	writeJSON(w, 200, s.db.GetHabit(id))
}

func (s *Server) deleteHabit(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	s.db.DeleteHabit(id)
	s.db.DeleteExtras(resourceName, id)
	writeJSON(w, 200, map[string]string{"deleted": "ok"})
}

// ─── check-ins ────────────────────────────────────────────────────

func (s *Server) checkIn(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req struct {
		Date string `json:"date"`
		Note string `json:"note"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	if err := s.db.CheckIn(id, req.Date, req.Note); err != nil {
		writeErr(w, 500, "checkin failed")
		return
	}
	writeJSON(w, 200, s.db.GetHabit(id))
}

func (s *Server) uncheck(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req struct {
		Date string `json:"date"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	s.db.Uncheck(id, req.Date)
	writeJSON(w, 200, s.db.GetHabit(id))
}

func (s *Server) history(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, 200, map[string]any{"checkins": orEmpty(s.db.ListCheckIns(r.PathValue("id"), 90))})
}

// ─── views & stats ────────────────────────────────────────────────

func (s *Server) today(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, 200, s.db.Today())
}

func (s *Server) stats(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, 200, s.db.Stats())
}

func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	st := s.db.Stats()
	writeJSON(w, 200, map[string]any{
		"status":  "ok",
		"service": "trailhead",
		"habits":  st.Habits,
		"streaks": st.ActiveStreaks,
	})
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
