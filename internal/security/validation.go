package security

import (
	"errors"
	"strings"
)

var ErrInvalidIdentifier = errors.New("invalid identifier")

func ValidateIdentifier(value string) error {
	value = strings.TrimSpace(value)

	if value == "" {
		return ErrInvalidIdentifier
	}

	if len(value) > 128 {
		return ErrInvalidIdentifier
	}

	for _, character := range value {
		if (character >= 'a' && character <= 'z') ||
			(character >= 'A' && character <= 'Z') ||
			(character >= '0' && character <= '9') ||
			character == '-' ||
			character == '_' ||
			character == '.' {
			continue
		}

		return ErrInvalidIdentifier
	}

	return nil
}
