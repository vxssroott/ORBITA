package telemetry

import (
	"fmt"
	"math"
	"time"

	"github.com/vxssroott/ORBITA/pkg/protocol"
)

func ValidateEnvelope(envelope protocol.TelemetryEnvelope) error {
	if envelope.SpacecraftID == "" {
		return fmt.Errorf("spacecraft ID is required")
	}

	if envelope.Timestamp.IsZero() {
		return fmt.Errorf("timestamp is required")
	}

	if envelope.Timestamp.After(time.Now().Add(5 * time.Minute)) {
		return fmt.Errorf("timestamp is too far in the future")
	}

	if envelope.Sequence == 0 {
		return fmt.Errorf("sequence must be greater than zero")
	}

	if len(envelope.Parameters) == 0 {
		return fmt.Errorf("telemetry parameters are required")
	}

	for key, value := range envelope.Parameters {
		if key == "" {
			return fmt.Errorf("telemetry parameter name is required")
		}

		switch typed := value.(type) {
		case float64:
			if math.IsNaN(typed) || math.IsInf(typed, 0) {
				return fmt.Errorf("telemetry parameter %q contains a non-finite value", key)
			}
		case float32:
			if math.IsNaN(float64(typed)) || math.IsInf(float64(typed), 0) {
				return fmt.Errorf("telemetry parameter %q contains a non-finite value", key)
			}
		case int:
		case int8:
		case int16:
		case int32:
		case int64:
		case uint:
		case uint8:
		case uint16:
		case uint32:
		case uint64:
		case string:
		case bool:
		case nil:
			return fmt.Errorf("telemetry parameter %q has a nil value", key)
		default:
			return fmt.Errorf("telemetry parameter %q has unsupported type %T", key, value)
		}
	}

	return nil
}
