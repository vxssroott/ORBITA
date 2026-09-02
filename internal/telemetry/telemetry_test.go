package telemetry

import (
	"testing"
	"time"

	"github.com/vxssroott/ORBITA/pkg/protocol"
)

func TestValidator(t *testing.T) {
	validator := NewValidator()

	valid := protocol.TelemetryEnvelope{
		SchemaVersion: "v1",
		SpacecraftID:  "ORBITA-01",
		Timestamp:     time.Now().UTC(),
		Parameters: map[string]any{
			"battery_voltage": 28.5,
		},
	}

	if err := validator.Validate(&valid); err != nil {
		t.Fatalf("valid envelope rejected: %v", err)
	}

	invalid := valid
	invalid.SpacecraftID = ""

	if err := validator.Validate(&invalid); err == nil {
		t.Fatal("invalid envelope accepted")
	}
}

func TestIsFresh(t *testing.T) {
	now := time.Now().UTC()

	if !IsFresh(now.Add(-10*time.Second), time.Minute, now) {
		t.Fatal("expected telemetry to be fresh")
	}

	if IsFresh(now.Add(-2*time.Minute), time.Minute, now) {
		t.Fatal("expected telemetry to be stale")
	}
}
