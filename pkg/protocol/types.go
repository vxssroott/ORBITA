package protocol

import "time"

type TelemetryEnvelope struct {
	SchemaVersion string         `json:"schema_version"`
	SpacecraftID  string         `json:"spacecraft_id"`
	Timestamp     time.Time      `json:"timestamp"`
	Sequence      uint64         `json:"sequence"`
	Parameters    map[string]any `json:"parameters"`
}

type TelemetryValue struct {
	Name      string    `json:"name"`
	Value     float64   `json:"value"`
	Unit      string    `json:"unit"`
	Timestamp time.Time `json:"timestamp"`
}

type Spacecraft struct {
	ID          string `json:"spacecraft_id"`
	Platform    string `json:"platform"`
	Version     string `json:"version"`
	Description string `json:"description"`
}

type SpacecraftState struct {
	SpacecraftID string             `json:"spacecraft_id"`
	UpdatedAt    time.Time          `json:"updated_at"`
	Mode         string             `json:"mode"`
	Health       string             `json:"health"`
	Parameters   map[string]float64 `json:"parameters"`
}

type EventSeverity string

const (
	SeverityInfo     EventSeverity = "info"
	SeverityWarning  EventSeverity = "warning"
	SeverityCritical EventSeverity = "critical"
)

type OperationalEvent struct {
	ID           string         `json:"event_id"`
	Type         string         `json:"event_type"`
	Severity     EventSeverity  `json:"severity"`
	SpacecraftID string         `json:"spacecraft_id"`
	Timestamp    time.Time      `json:"timestamp"`
	Message      string         `json:"message"`
	Metadata     map[string]any `json:"metadata,omitempty"`
}

type Command struct {
	ID           string         `json:"command_id"`
	SpacecraftID string         `json:"spacecraft_id"`
	Type         string         `json:"command_type"`
	Parameters   map[string]any `json:"parameters,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
}

type CommandResult struct {
	CommandID    string    `json:"command_id"`
	Accepted     bool      `json:"accepted"`
	Acknowledged bool      `json:"acknowledged"`
	Timestamp    time.Time `json:"timestamp"`
	Message      string    `json:"message"`
}
