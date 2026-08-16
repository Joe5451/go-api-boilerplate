package grpckit

import (
	"errors"
	"fmt"

	"go-api-boilerplate/internal/domain"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Error struct {
	Code    codes.Code
	Message string
}

func (e *Error) Error() string {
	return fmt.Sprintf("grpc error: code: %s, message: %s", e.Code, e.Message)
}

func NewValidationError(msg string) *Error {
	return &Error{
		Code:    codes.InvalidArgument,
		Message: msg,
	}
}

func ToStatus(err error) error {
	if err == nil {
		return nil
	}

	if _, ok := status.FromError(err); ok {
		return err
	}

	var grpcErr *Error
	if errors.As(err, &grpcErr) {
		return status.Error(grpcErr.Code, grpcErr.Message)
	}

	switch {
	case errors.Is(err, domain.ErrBookNotFound):
		return status.Error(codes.NotFound, domain.ErrBookNotFound.Error())
	case errors.Is(err, domain.ErrTitleRequired), errors.Is(err, domain.ErrAuthorRequired):
		return status.Error(codes.InvalidArgument, err.Error())
	default:
		return status.Error(codes.Internal, "internal server error")
	}
}
