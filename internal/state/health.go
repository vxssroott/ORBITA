package state

import (
	"fmt"
	"math"

	"github.com/vxssroott/ORBITA/pkg/protocol"
)

const (
	HealthUnknown  = "unknown"
	HealthNominal  = "nominal"
	HealthDegraded = "degraded"
	HealthWarning  = "warning"
	HealthCritical = "critical"
)

func CalculateHealth(values map[string]float64) string {
	if len(values) == 0 {
		return HealthUnknown
	}

	health := HealthNominal

	for key, value := range values {
		if math.IsNaN(value) || math.IsInf(value, 0) {
			return HealthCritical
		}

		switch key {
		case "battery_voltage":
			if value < 20 {
				return HealthCritical
			}

		case "battery":
			if value < 15 {
				if health == HealthNominal {
					health = HealthWarning
				}
			}

		case "temperature":
			if value > 80 || value < -20 {
				return HealthCritical
			}

			if value > 70 || value < -10 {
				if health == HealthNominal {
					health = HealthWarning
				}
			}

		case "signal_strength":
			if value < 20 {
				if health == HealthNominal {
					health = HealthDegraded
				}
			}
		}
	}

	return health
}

func ValidateState(state protocol.SpacecraftState) error {
	if state.SpacecraftID == "" {
		return fmt.Errorf("spacecraft ID is required")
	}

	if state.UpdatedAt.IsZero() {
		return fmt.Errorf("state update time is required")
	}

	if state.Parameters == nil {
		return fmt.Errorf("state parameters are required")
	}

	switch state.Health {
	case HealthUnknown, HealthNominal, HealthDegraded, HealthWarning, HealthCritical:
	default:
		return fmt.Errorf("invalid spacecraft health %q", state.Health)
	}

	return nil
}
