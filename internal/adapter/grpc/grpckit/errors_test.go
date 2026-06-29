package grpckit

import (
	"errors"
	"testing"

	"go-api-boilerplate/internal/domain"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestToStatus(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		wantCode codes.Code
	}{
		{name: "nil", err: nil, wantCode: codes.OK},
		{name: "not_found", err: domain.ErrBookNotFound, wantCode: codes.NotFound},
		{name: "title_required", err: domain.ErrTitleRequired, wantCode: codes.InvalidArgument},
		{name: "author_required", err: domain.ErrAuthorRequired, wantCode: codes.InvalidArgument},
		{name: "validation", err: NewValidationError("id must be greater than 0"), wantCode: codes.InvalidArgument},
		{name: "already_status", err: status.Error(codes.NotFound, "book not found"), wantCode: codes.NotFound},
		{name: "unknown", err: errors.New("boom"), wantCode: codes.Internal},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ToStatus(tt.err)
			if status.Code(got) != tt.wantCode {
				t.Fatalf("expected %s, got %s", tt.wantCode, status.Code(got))
			}
		})
	}
}
