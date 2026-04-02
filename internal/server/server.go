package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/stockyard-dev/stockyard-trailhead/internal/store"
)

type Server struct { db *store.DB; mux *http.ServeMux }

func New(db *store.DB) *Server {
	s := &Server{db: db, mux: http.NewServeMux()}
	s.mux.HandleFunc("GET /api/habits", s.listHabits)
	s.mux.HandleFunc("POST /api/habits", s.createHabit)
	s.mux.HandleFunc("GET /api/habits/{id}", s.getHabit)
	s.mux.HandleFunc("PUT /api/habits/{id}", s.updateHabit)
	s.mux.HandleFunc("DELETE /api/habits/{id}", s.deleteHabit)
	s.mux.HandleFunc("POST /api/habits/{id}/check", s.checkIn)
	s.mux.HandleFunc("POST /api/habits/{id}/uncheck", s.uncheck)
	s.mux.HandleFunc("GET /api/habits/{id}/history", s.history)
	s.mux.HandleFunc("GET /api/today", s.today)
	s.mux.HandleFunc("GET /api/stats", s.stats)
	s.mux.HandleFunc("GET /api/health", s.health)
	s.mux.HandleFunc("GET /ui", s.dashboard)
	s.mux.HandleFunc("GET /ui/", s.dashboard)
	s.mux.HandleFunc("GET /", s.root)
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) { s.mux.ServeHTTP(w, r) }
func writeJSON(w http.ResponseWriter, code int, v any) { w.Header().Set("Content-Type","application/json"); w.WriteHeader(code); json.NewEncoder(w).Encode(v) }
func writeErr(w http.ResponseWriter, code int, msg string) { writeJSON(w, code, map[string]string{"error": msg}) }
func (s *Server) root(w http.ResponseWriter, r *http.Request) { if r.URL.Path != "/" { http.NotFound(w, r); return }; http.Redirect(w, r, "/ui", http.StatusFound) }

func (s *Server) listHabits(w http.ResponseWriter, r *http.Request) {
	incArch := r.URL.Query().Get("archived") == "true"
	writeJSON(w, 200, map[string]any{"habits": orEmpty(s.db.ListHabits(incArch))})
}
func (s *Server) createHabit(w http.ResponseWriter, r *http.Request) {
	var h store.Habit; json.NewDecoder(r.Body).Decode(&h)
	if h.Name == "" { writeErr(w, 400, "name required"); return }
	if err := s.db.CreateHabit(&h); err != nil { writeErr(w, 500, err.Error()); return }
	writeJSON(w, 201, s.db.GetHabit(h.ID))
}
func (s *Server) getHabit(w http.ResponseWriter, r *http.Request) {
	h := s.db.GetHabit(r.PathValue("id")); if h == nil { writeErr(w, 404, "not found"); return }; writeJSON(w, 200, h)
}
func (s *Server) updateHabit(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id"); ex := s.db.GetHabit(id); if ex == nil { writeErr(w, 404, "not found"); return }
	var h store.Habit; json.NewDecoder(r.Body).Decode(&h)
	if h.Name == "" { h.Name = ex.Name }; if h.Frequency == "" { h.Frequency = ex.Frequency }
	if h.Color == "" { h.Color = ex.Color }
	s.db.UpdateHabit(id, &h); writeJSON(w, 200, s.db.GetHabit(id))
}
func (s *Server) deleteHabit(w http.ResponseWriter, r *http.Request) { s.db.DeleteHabit(r.PathValue("id")); writeJSON(w, 200, map[string]string{"deleted":"ok"}) }

func (s *Server) checkIn(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id"); var req struct{ Date string `json:"date"`; Note string `json:"note"` }
	json.NewDecoder(r.Body).Decode(&req)
	if err := s.db.CheckIn(id, req.Date, req.Note); err != nil { writeErr(w, 500, err.Error()); return }
	writeJSON(w, 200, s.db.GetHabit(id))
}
func (s *Server) uncheck(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id"); var req struct{ Date string `json:"date"` }; json.NewDecoder(r.Body).Decode(&req)
	s.db.Uncheck(id, req.Date); writeJSON(w, 200, s.db.GetHabit(id))
}
func (s *Server) history(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, 200, map[string]any{"checkins": orEmpty(s.db.ListCheckIns(r.PathValue("id"), 90))})
}
func (s *Server) today(w http.ResponseWriter, r *http.Request) { writeJSON(w, 200, s.db.Today()) }
func (s *Server) stats(w http.ResponseWriter, r *http.Request) { writeJSON(w, 200, s.db.Stats()) }
func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	st := s.db.Stats(); writeJSON(w, 200, map[string]any{"status":"ok","service":"trailhead","habits":st.Habits,"streaks":st.ActiveStreaks})
}
func orEmpty[T any](s []T) []T { if s == nil { return []T{} }; return s }
func init() { log.SetFlags(log.LstdFlags | log.Lshortfile) }
