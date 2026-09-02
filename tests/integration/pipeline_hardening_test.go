package integration

import (
	"testing"
	"time"

	"github.com/vxssroott/ORBITA/internal/events"
	"github.com/vxssroott/ORBITA/internal/rules"
	"github.com/vxssroott/ORBITA/internal/state"
	"github.com/vxssroott/ORBITA/internal/telemetry"
	"github.com/vxssroott/ORBITA/pkg/protocol"
)

func TestTelemetryValidationRejectsMalformedPackets(t *testing.T) {
	validator := telemetry.NewValidator()

	valid := protocol.TelemetryEnvelope{
		SchemaVersion: "1.0",
		SpacecraftID:  "NIGCOMSAT-SIM-01",
		Timestamp:     time.Now().UTC(),
		Sequence:      1,
		Parameters: map[string]any{
			"temperature":     25.0,
			"battery_voltage": 28.0,
			"signal_strength": 92.0,
		},
	}

	cases := []struct {
		name     string
		envelope *protocol.TelemetryEnvelope
	}{
		{
			name:     "nil envelope",
			envelope: nil,
		},
		{
			name: "missing spacecraft",
			envelope: func() *protocol.TelemetryEnvelope {
				v := valid
				v.SpacecraftID = ""
				return &v
			}(),
		},
		{
			name: "missing timestamp",
			envelope: func() *protocol.TelemetryEnvelope {
				v := valid
				v.Timestamp = time.Time{}
				return &v
			}(),
		},
		{
			name: "missing parameters",
			envelope: func() *protocol.TelemetryEnvelope {
				v := valid
				v.Parameters = nil
				return &v
			}(),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if err := validator.Validate(tc.envelope); err == nil {
				t.Fatalf("expected malformed telemetry to be rejected")
			}
		})
	}

	if err := validator.Validate(&valid); err != nil {
		t.Fatalf("valid telemetry rejected: %v", err)
	}
}

func TestTelemetryStoreRetainsLatestEnvelope(t *testing.T) {
	store := telemetry.NewStore()

	first := protocol.TelemetryEnvelope{
		SchemaVersion: "1.0",
		SpacecraftID:  "NIGCOMSAT-SIM-01",
		Timestamp:     time.Now().UTC(),
		Sequence:      1,
		Parameters:    map[string]any{"temperature": 25.0},
	}

	second := first
	second.Sequence = 2
	second.Timestamp = first.Timestamp.Add(time.Second)
	second.Parameters = map[string]any{"temperature": 26.0}

	store.Put(first)
	store.Put(second)

	latest, ok := store.Latest(first.SpacecraftID)
	if !ok {
		t.Fatalf("expected latest telemetry")
	}

	if latest.Sequence != 2 {
		t.Fatalf("expected latest sequence 2, got %d", latest.Sequence)
	}
}

func TestStateRecoveryFromCriticalAnomaly(t *testing.T) {
	engine := state.NewEngine()

	critical := protocol.TelemetryEnvelope{
		SpacecraftID: "NIGCOMSAT-SIM-01",
		Timestamp:    time.Now().UTC(),
		Sequence:     1,
		Parameters: map[string]any{
			"temperature":     95.0,
			"battery_voltage": 17.5,
			"signal_strength": 12.0,
		},
	}

	recovered := critical
	recovered.Sequence = 2
	recovered.Timestamp = critical.Timestamp.Add(time.Second)
	recovered.Parameters = map[string]any{
		"temperature":     25.0,
		"battery_voltage": 28.0,
		"signal_strength": 92.0,
	}

	criticalState, err := engine.Apply(critical)
	if err != nil {
		t.Fatalf("critical telemetry rejected: %v", err)
	}

	if criticalState.Health != "critical" {
		t.Fatalf("expected critical health, got %q", criticalState.Health)
	}

	recoveredState, err := engine.Apply(recovered)
	if err != nil {
		t.Fatalf("recovery telemetry rejected: %v", err)
	}

	if recoveredState.Health != "nominal" {
		t.Fatalf("expected nominal recovery, got %q", recoveredState.Health)
	}
}

func TestRepeatedAnomalyDetectionIsDeterministic(t *testing.T) {
	maxTemperature := 80.0
	minBattery := 20.0
	minSignal := 20.0

	engine := rules.NewEngine([]rules.Rule{
		{
			Name:      "temperature-critical-high",
			Parameter: "temperature",
			Max:       &maxTemperature,
			Severity:  rules.SeverityCritical,
			EventType: "thermal_anomaly",
			Message:   "spacecraft temperature exceeds safe operating threshold",
		},
		{
			Name:      "battery-critical-low",
			Parameter: "battery_voltage",
			Min:       &minBattery,
			Severity:  rules.SeverityCritical,
			EventType: "battery_anomaly",
			Message:   "spacecraft battery voltage is below safe operating threshold",
		},
		{
			Name:      "signal-degraded-low",
			Parameter: "signal_strength",
			Min:       &minSignal,
			Severity:  rules.SeverityWarning,
			EventType: "communications_degraded",
			Message:   "spacecraft signal strength is below operational threshold",
		},
	})

	parameters := map[string]any{
		"temperature":     95.0,
		"battery_voltage": 17.5,
		"signal_strength": 12.0,
	}

	first := engine.Evaluate("NIGCOMSAT-SIM-01", parameters)
	second := engine.Evaluate("NIGCOMSAT-SIM-01", parameters)

	if len(first) != 3 || len(second) != 3 {
		t.Fatalf("expected 3 events per evaluation, got %d and %d", len(first), len(second))
	}

	for i := range first {
		if first[i].Type != second[i].Type {
			t.Fatalf("non-deterministic event type at index %d: %q != %q", i, first[i].Type, second[i].Type)
		}

		if first[i].Severity != second[i].Severity {
			t.Fatalf("non-deterministic severity at index %d", i)
		}

		if first[i].Message != second[i].Message {
			t.Fatalf("non-deterministic message at index %d", i)
		}
	}
}

func TestEventEngineRecordsOperationalEvents(t *testing.T) {
	engine := events.NewEngine()

	event := protocol.OperationalEvent{
		Type:         "thermal_anomaly",
		Severity:     protocol.SeverityCritical,
		SpacecraftID: "NIGCOMSAT-SIM-01",
		Message:      "spacecraft temperature exceeds safe operating threshold",
	}

	recorded := engine.Emit(event)

	if recorded.ID == "" {
		t.Fatal("event ID was not assigned")
	}

	if recorded.Timestamp.IsZero() {
		t.Fatal("event timestamp was not assigned")
	}

	recent := engine.Recent(10)
	if len(recent) != 1 {
		t.Fatalf("expected 1 recorded event, got %d", len(recent))
	}

	if recent[0].ID != recorded.ID {
		t.Fatalf("recorded event mismatch")
	}
}
