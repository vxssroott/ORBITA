package protocol

import "fmt"

type ErrorCode string

const (
	ErrorInvalidInput ErrorCode = "INVALID_INPUT"
	ErrorNotFound     ErrorCode = "NOT_FOUND"
	ErrorConflict     ErrorCode = "CONFLICT"
	ErrorUnauthorized ErrorCode = "UNAUTHORIZED"
	ErrorForbidden    ErrorCode = "FORBIDDEN"
	ErrorUnavailable  ErrorCode = "UNAVAILABLE"
	ErrorInternal     ErrorCode = "INTERNAL"
)

type OperationalError struct {
	Code    ErrorCode
	Message string
	Err     error
}

func (e *OperationalError) Error() string {
	if e.Err == nil {
		return e.Message
	}

	return fmt.Sprintf("%s: %v", e.Message, e.Err)
}

func (e *OperationalError) Unwrap() error {
	return e.Err
}

func NewOperationalError(code ErrorCode, message string, err error) error {
	return &OperationalError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}
