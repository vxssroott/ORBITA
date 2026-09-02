package state

type Health string

const (
	HealthUnknown  Health = "unknown"
	HealthNominal  Health = "nominal"
	HealthDegraded Health = "degraded"
	HealthCritical Health = "critical"
)

func CalculateHealth(parameters map[string]float64) Health {
	health := HealthNominal

	for key, value := range parameters {
		switch key {
		case "battery_voltage":
			if value < 20 {
				return HealthCritical
			}

		case "temperature":
			if value > 80 || value < -20 {
				return HealthCritical
			}

		case "signal_strength":
			if value < 20 {
				health = HealthDegraded
			}
		}
	}

	return health
}
