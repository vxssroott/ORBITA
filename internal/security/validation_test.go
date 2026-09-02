package security

import "testing"

func TestValidateIdentifier(t *testing.T) {
	valid := []string{
		"ORBITA-01",
		"spacecraft_001",
		"ground.station",
	}

	for _, value := range valid {
		if err := ValidateIdentifier(value); err != nil {
			t.Fatalf("expected %q to be valid: %v", value, err)
		}
	}

	invalid := []string{
		"",
		" ",
		"spacecraft/01",
		"spacecraft 01",
	}

	for _, value := range invalid {
		if err := ValidateIdentifier(value); err == nil {
			t.Fatalf("expected %q to be invalid", value)
		}
	}
}
