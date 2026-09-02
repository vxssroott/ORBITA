package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/vxssroott/ORBITA/internal/security"
)

type Server struct {
	mux     *http.ServeMux
	maxBody int64
}

func NewServer(maxBody int64) *Server {
	if maxBody <= 0 {
		maxBody = 1 << 20
	}

	s := &Server{
		mux:     http.NewServeMux(),
		maxBody: maxBody,
	}

	s.routes()

	return s
}

func (s *Server) routes() {
	s.mux.HandleFunc("/healthz", s.health)
	s.mux.HandleFunc("/readyz", s.ready)
	s.mux.HandleFunc("/v1/status", s.status)
}

func (s *Server) Handler() http.Handler {
	return security.RequestID(
		security.SecurityHeaders(
			http.MaxBytesHandler(s.mux, s.maxBody),
		),
	)
}

func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"status": "ok",
	})
}

func (s *Server) ready(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"status": "ready",
	})
}

func (s *Server) status(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"service":    "orbita-api",
		"status":     "operational",
		"version":    "v1",
		"request_id": security.GetRequestID(r.Context()),
		"timestamp":  time.Now().UTC(),
	})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(value)
}
