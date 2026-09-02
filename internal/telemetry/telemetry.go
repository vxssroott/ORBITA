package telemetry

import (
	"errors"
	"sync"
	"time"

	"github.com/vxssroott/ORBITA/pkg/protocol"
)

var (
	ErrNilEnvelope       = errors.New("telemetry envelope is nil")
	ErrMissingSpacecraft = errors.New("spacecraft id is required")
	ErrMissingTimestamp  = errors.New("telemetry timestamp is required")
)

type Validator struct{}

func NewValidator() *Validator {
	return &Validator{}
}

func (v *Validator) Validate(envelope *protocol.TelemetryEnvelope) error {
	if envelope == nil {
		return ErrNilEnvelope
	}

	if envelope.SpacecraftID == "" {
		return ErrMissingSpacecraft
	}

	if envelope.Timestamp.IsZero() {
		return ErrMissingTimestamp
	}

	if envelope.Parameters == nil {
		return errors.New("telemetry parameters are required")
	}

	return nil
}

type Store struct {
	mu     sync.RWMutex
	latest map[string]protocol.TelemetryEnvelope
}

func NewStore() *Store {
	return &Store{
		latest: make(map[string]protocol.TelemetryEnvelope),
	}
}

func (s *Store) Put(envelope protocol.TelemetryEnvelope) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.latest[envelope.SpacecraftID] = envelope
}

func (s *Store) Latest(spacecraftID string) (protocol.TelemetryEnvelope, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	value, ok := s.latest[spacecraftID]
	return value, ok
}

func IsFresh(timestamp time.Time, maxAge time.Duration, now time.Time) bool {
	if timestamp.IsZero() || maxAge <= 0 {
		return false
	}

	age := now.Sub(timestamp)

	return age >= 0 && age <= maxAge
}
