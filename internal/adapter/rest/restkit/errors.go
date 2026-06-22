package restkit

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type (
	ErrorCode string
	Error     struct {
		Status  int
		Code    ErrorCode
		Message string
	}
)

const (
	ErrCodeInternalError    ErrorCode = "INTERNAL_ERROR"
	ErrCodeValidationFailed ErrorCode = "VALIDATION_FAILED"
	ErrCodeNotFound         ErrorCode = "NOT_FOUND"
)

var ErrorMessages = map[string]string{
	"internal_error": "Internal server error.",
	"not_found":      "Resource not found.",
}

func (e *Error) Error() string {
	return fmt.Sprintf("rest error: status: %d, code: %s, message: %s", e.Status, e.Code, e.Message)
}

func NewInternalServerError() *Error {
	return &Error{
		Status:  http.StatusInternalServerError,
		Code:    ErrCodeInternalError,
		Message: ErrorMessages["internal_error"],
	}
}

func NewValidationError(err error) *Error {
	if err == nil {
		return nil
	}

	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return &Error{
			Status:  http.StatusBadRequest,
			Code:    ErrCodeValidationFailed,
			Message: err.Error(),
		}
	}

	for _, fieldErr := range validationErrors {
		return &Error{
			Status:  http.StatusBadRequest,
			Code:    ErrCodeValidationFailed,
			Message: NewFieldError(fieldErr).String(),
		}
	}

	return nil
}

func NewNotFoundError(msg string) *Error {
	return &Error{
		Status:  http.StatusNotFound,
		Code:    ErrCodeNotFound,
		Message: msg,
	}
}
