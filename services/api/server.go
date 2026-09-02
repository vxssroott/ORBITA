package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/vxssroott/ORBITA/internal/security"
	"github.com/vxssroott/ORBITA/internal/state"
	"github.com/vxssroott/ORBITA/internal/telemetry"
)

type Server struct {
	telemetry *telemetry.Store
	state     *state.Engine
	mux       *http.ServeMux
}

func NewServer(
	telemetryStore *telemetry.Store,
	stateEngine *state.Engine,
) *Server {
	s := &Server{
		telemetry: telemetryStore,
		state:     stateEngine,
		mux:       http.NewServeMux(),
	}

	s.routes()

	return s
}

func (s *Server) routes() {
	s.mux.HandleFunc("/health", s.health)
	s.mux.HandleFunc("/ready", s.ready)
	s.mux.HandleFunc("/v1/spacecraft/state", s.spacecraftState)
	s.mux.HandleFunc("/v1/spacecraft/telemetry", s.spacecraftTelemetry)
}

func (s *Server) Handler() http.Handler {
	return security.SecurityHeaders(s.mux)
}

func (s *Server) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"status":    "ok",
		"timestamp": time.Now().UTC(),
		"service":   "orbita-api",
	})
}

func (s *Server) ready(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status": "ready",
	})
}

func (s *Server) spacecraftState(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	if err := security.ValidateIdentifier(id); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "valid spacecraft id is required",
		})
		return
	}

	value, ok := s.state.Get(id)

	if !ok {
		writeJSON(w, http.StatusNotFound, map[string]string{
			"error": "spacecraft state not found",
		})
		return
	}

	writeJSON(w, http.StatusOK, value)
}

func (s *Server) spacecraftTelemetry(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	if err := security.ValidateIdentifier(id); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "valid spacecraft id is required",
		})
		return
	}

	value, ok := s.telemetry.Latest(id)

	if !ok {
		writeJSON(w, http.StatusNotFound, map[string]string{
			"error": "telemetry not found",
		})
		return
	}

	writeJSON(w, http.StatusOK, value)
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(value); err != nil {
		return
	}
}
