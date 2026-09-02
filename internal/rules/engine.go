package rules

import (
	"github.com/vxssroott/ORBITA/pkg/protocol"
)

type Severity string

const (
	SeverityWarning  Severity = "warning"
	SeverityCritical Severity = "critical"
)

type Rule struct {
	Name      string
	Parameter string
	Min       *float64
	Max       *float64
	Severity  Severity
	EventType string
	Message   string
}

type Engine struct {
	rules []Rule
}

func NewEngine(rules []Rule) *Engine {
	copyRules := make([]Rule, len(rules))
	copy(copyRules, rules)

	return &Engine{rules: copyRules}
}

func (e *Engine) Evaluate(
	spacecraftID string,
	parameters map[string]any,
) []protocol.OperationalEvent {
	events := make([]protocol.OperationalEvent, 0)

	for _, rule := range e.rules {
		raw, exists := parameters[rule.Parameter]
		if !exists {
			continue
		}

		value, ok := toFloat64(raw)
		if !ok {
			continue
		}

		violated := false

		if rule.Min != nil && value < *rule.Min {
			violated = true
		}

		if rule.Max != nil && value > *rule.Max {
			violated = true
		}

		if !violated {
			continue
		}

		severity := protocol.SeverityWarning

		if rule.Severity == SeverityCritical {
			severity = protocol.SeverityCritical
		}

		events = append(events, protocol.OperationalEvent{
			Type:         rule.EventType,
			Severity:     severity,
			SpacecraftID: spacecraftID,
			Message:      rule.Message,
			Metadata: map[string]any{
				"parameter": rule.Parameter,
				"value":     value,
				"rule":      rule.Name,
			},
		})
	}

	return events
}

func toFloat64(value any) (float64, bool) {
	switch v := value.(type) {
	case float64:
		return v, true
	case float32:
		return float64(v), true
	case int:
		return float64(v), true
	case int64:
		return float64(v), true
	case int32:
		return float64(v), true
	default:
		return 0, false
	}
}
