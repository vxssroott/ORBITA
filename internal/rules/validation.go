package rules

import (
	"fmt"
	"strings"
)

func ValidateRule(rule Rule) error {
	if strings.TrimSpace(rule.Name) == "" {
		return fmt.Errorf("rule name is required")
	}

	if strings.TrimSpace(rule.Parameter) == "" {
		return fmt.Errorf("rule parameter is required")
	}

	if rule.Min == nil && rule.Max == nil {
		return fmt.Errorf("rule must define a minimum or maximum")
	}

	if rule.Min != nil && rule.Max != nil && *rule.Min > *rule.Max {
		return fmt.Errorf("rule minimum cannot exceed maximum")
	}

	if strings.TrimSpace(rule.EventType) == "" {
		return fmt.Errorf("rule event type is required")
	}

	if strings.TrimSpace(rule.Message) == "" {
		return fmt.Errorf("rule message is required")
	}

	switch rule.Severity {
	case SeverityWarning, SeverityCritical:
	default:
		return fmt.Errorf("unsupported rule severity")
	}

	return nil
}

func ValidateRules(rules []Rule) error {
	if len(rules) > 4096 {
		return fmt.Errorf("rule set exceeds maximum size")
	}

	for _, rule := range rules {
		if err := ValidateRule(rule); err != nil {
			return err
		}
	}

	return nil
}
