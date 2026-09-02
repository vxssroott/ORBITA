package audit

import (
	"sync"
	"time"
)

type Entry struct {
	ID        string
	Timestamp time.Time
	Actor     string
	Action    string
	Resource  string
	Outcome   string
	Metadata  map[string]string
}

type Log struct {
	mu      sync.RWMutex
	entries []Entry
}

func NewLog() *Log {
	return &Log{
		entries: make([]Entry, 0, 256),
	}
}

func (l *Log) Append(entry Entry) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if entry.Timestamp.IsZero() {
		entry.Timestamp = time.Now().UTC()
	}

	l.entries = append(l.entries, entry)
}

func (l *Log) List(limit int) []Entry {
	l.mu.RLock()
	defer l.mu.RUnlock()

	if limit <= 0 || limit > len(l.entries) {
		limit = len(l.entries)
	}

	start := len(l.entries) - limit

	result := make([]Entry, limit)
	copy(result, l.entries[start:])

	return result
}
