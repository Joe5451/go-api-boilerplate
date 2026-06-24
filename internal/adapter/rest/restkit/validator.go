package restkit

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

type fieldError struct {
	err validator.FieldError
}

func NewFieldError(err validator.FieldError) *fieldError {
	return &fieldError{err: err}
}

func (q fieldError) String() string {
	var sb strings.Builder
	sb.WriteString("validation failed on field '" + q.err.Field() + "'")
	sb.WriteString(", condition: " + q.err.ActualTag())
	if q.err.Param() != "" {
		sb.WriteString(" { " + q.err.Param() + " }")
	}
	if q.err.Value() != nil && q.err.Value() != "" {
		sb.WriteString(fmt.Sprintf(", actual: %v", q.err.Value()))
	}
	return sb.String()
}
