package events

import (
	"fmt"
	"sync"
	"time"

	"github.com/vxssroott/ORBITA/pkg/protocol"
)

type Engine struct {
	mu     sync.RWMutex
	events []protocol.OperationalEvent
}

func NewEngine() *Engine {
	return &Engine{
		events: make([]protocol.OperationalEvent, 0, 128),
	}
}

func (e *Engine) Emit(event protocol.OperationalEvent) protocol.OperationalEvent {
	e.mu.Lock()
	defer e.mu.Unlock()

	if event.ID == "" {
		event.ID = fmt.Sprintf("evt-%d", time.Now().UnixNano())
	}

	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now().UTC()
	}

	e.events = append(e.events, event)

	return event
}

func (e *Engine) Recent(limit int) []protocol.OperationalEvent {
	e.mu.RLock()
	defer e.mu.RUnlock()

	if limit <= 0 || limit > len(e.events) {
		limit = len(e.events)
	}

	start := len(e.events) - limit
	result := make([]protocol.OperationalEvent, limit)

	copy(result, e.events[start:])

	return result
}
