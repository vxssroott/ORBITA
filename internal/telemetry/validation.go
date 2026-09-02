package telemetry

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/vxssroott/ORBITA/pkg/protocol"
)

func ValidateEnvelope(envelope protocol.TelemetryEnvelope) error {
	if strings.TrimSpace(envelope.SchemaVersion) == "" {
		return fmt.Errorf("schema version is required")
	}

	if strings.TrimSpace(envelope.SpacecraftID) == "" {
		return fmt.Errorf("spacecraft id is required")
	}

	if envelope.Timestamp.IsZero() {
		return fmt.Errorf("timestamp is required")
	}

	if envelope.Timestamp.Location() == nil {
		return fmt.Errorf("timestamp location is invalid")
	}

	if envelope.Timestamp.After(time.Now().UTC().Add(5 * time.Minute)) {
		return fmt.Errorf("telemetry timestamp is too far in the future")
	}

	if envelope.Parameters == nil {
		return fmt.Errorf("parameters are required")
	}

	if len(envelope.Parameters) > 512 {
		return fmt.Errorf("too many telemetry parameters")
	}

	for name, value := range envelope.Parameters {
		if strings.TrimSpace(name) == "" {
			return fmt.Errorf("telemetry parameter name cannot be empty")
		}

		switch typed := value.(type) {
		case float64:
			if math.IsNaN(typed) || math.IsInf(typed, 0) {
				return fmt.Errorf("parameter %q contains invalid numeric value", name)
			}
		case float32:
			if math.IsNaN(float64(typed)) || math.IsInf(float64(typed), 0) {
				return fmt.Errorf("parameter %q contains invalid numeric value", name)
			}
		}
	}

	return nil
}
