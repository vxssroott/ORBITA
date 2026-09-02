package state

import (
	"testing"
	"time"

	"github.com/vxssroott/ORBITA/pkg/protocol"
)

func TestApply(t *testing.T) {
	engine := NewEngine()

	envelope := protocol.TelemetryEnvelope{
		SchemaVersion: "v1",
		SpacecraftID:  "ORBITA-01",
		Timestamp:     time.Now().UTC(),
		Parameters: map[string]any{
			"battery_voltage": 28.5,
			"temperature":     35.0,
			"signal_strength": 80.0,
		},
	}

	result, err := engine.Apply(envelope)
	if err != nil {
		t.Fatalf("apply failed: %v", err)
	}

	if result.Health != "nominal" {
		t.Fatalf("expected nominal health, got %s", result.Health)
	}

	if result.Parameters["battery_voltage"] != 28.5 {
		t.Fatal("battery voltage was not persisted")
	}
}
