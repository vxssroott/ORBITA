package rules

import "testing"

func TestEngineEvaluate(t *testing.T) {
	max := 30.0

	engine := NewEngine([]Rule{
		{
			Name:      "battery-high",
			Parameter: "battery_voltage",
			Max:       &max,
			Severity:  SeverityCritical,
			EventType: "battery_threshold",
			Message:   "battery voltage exceeds safe threshold",
		},
	})

	events := engine.Evaluate(
		"ORBITA-01",
		map[string]any{
			"battery_voltage": 35.0,
		},
	)

	if len(events) != 1 {
		t.Fatalf("expected one event, got %d", len(events))
	}

	if events[0].Severity != "critical" {
		t.Fatalf("expected critical severity, got %s", events[0].Severity)
	}
}
