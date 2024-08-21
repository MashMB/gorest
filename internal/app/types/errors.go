package types

import "time"

type FieldError struct {
	Field    string `json:"field"`
	Code     string `json:"code"`
	Message  string `json:"message"`
	Value    string `json:"value"`
	Expected string `json:"expected"`
}

func NewFieldError(field, code, message string) FieldError {
	return FieldError{
		Field:   field,
		Code:    code,
		Message: message,
	}
}

func NewDetailedFieldError(field, code, message, value, expected string) FieldError {
	return FieldError{
		Field:    field,
		Code:     code,
		Message:  message,
		Value:    value,
		Expected: expected,
	}
}

type ApiErrorDto struct {
	Timestamp string       `json:"timestamp"`
	Code      string       `json:"code"`
	Message   string       `json:"message"`
	Details   []FieldError `json:"details"`
}

func NewApiErrorDto(code, message string, details ...FieldError) ApiErrorDto {
	if len(details) == 0 {
		details = []FieldError{}
	}

	return ApiErrorDto{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Code:      code,
		Message:   message,
		Details:   details,
	}
}

type ApiError struct {
	Status  int
	Cause   string
	Details []FieldError
}

func (e *ApiError) Error() string {
	return e.Cause
}

func NewApiError(status int, cause string, details ...FieldError) *ApiError {
	if len(details) == 0 {
		details = []FieldError{}
	}

	return &ApiError{
		Status:  status,
		Cause:   cause,
		Details: details,
	}
}
