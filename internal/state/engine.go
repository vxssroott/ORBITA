package state

import (
	"errors"
	"sync"
	"time"

	"github.com/vxssroott/ORBITA/pkg/protocol"
)

var ErrMissingSpacecraftID = errors.New("spacecraft id is required")

type Engine struct {
	mu     sync.RWMutex
	states map[string]protocol.SpacecraftState
}

func NewEngine() *Engine {
	return &Engine{
		states: make(map[string]protocol.SpacecraftState),
	}
}

func (e *Engine) Apply(envelope protocol.TelemetryEnvelope) (protocol.SpacecraftState, error) {
	if envelope.SpacecraftID == "" {
		return protocol.SpacecraftState{}, ErrMissingSpacecraftID
	}

	e.mu.Lock()
	defer e.mu.Unlock()

	state, exists := e.states[envelope.SpacecraftID]

	if !exists {
		state = protocol.SpacecraftState{
			SpacecraftID: envelope.SpacecraftID,
			Mode:         "unknown",
			Health:       "unknown",
			Parameters:   make(map[string]float64),
		}
	}

	for key, value := range envelope.Parameters {
		switch v := value.(type) {
		case float64:
			state.Parameters[key] = v
		case int:
			state.Parameters[key] = float64(v)
		case int64:
			state.Parameters[key] = float64(v)
		}
	}

	state.UpdatedAt = envelope.Timestamp

	if mode, ok := envelope.Parameters["mode"]; ok {
		if value, ok := mode.(string); ok {
			state.Mode = value
		}
	}

	state.Health = CalculateHealth(state.Parameters)

	e.states[envelope.SpacecraftID] = state

	return state, nil
}

func (e *Engine) Get(spacecraftID string) (protocol.SpacecraftState, bool) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	state, ok := e.states[spacecraftID]
	return state, ok
}

func IsStateCurrent(state protocol.SpacecraftState, maxAge time.Duration, now time.Time) bool {
	if state.UpdatedAt.IsZero() || maxAge <= 0 {
		return false
	}

	age := now.Sub(state.UpdatedAt)
	return age >= 0 && age <= maxAge
}
