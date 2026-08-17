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
	sb.WriteString("validation failed on field '")
	sb.WriteString(q.err.Field())
	sb.WriteString("'")
	sb.WriteString(", condition: ")
	sb.WriteString(q.err.ActualTag())
	if q.err.Param() != "" {
		sb.WriteString(" { ")
		sb.WriteString(q.err.Param())
		sb.WriteString(" }")
	}
	if q.err.Value() != nil && q.err.Value() != "" {
		sb.WriteString(fmt.Sprintf(", actual: %v", q.err.Value()))
	}
	return sb.String()
}
