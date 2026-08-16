package interceptors

import (
	"context"
	"errors"
	"testing"

	"go-api-boilerplate/internal/domain"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestErrorHandler(t *testing.T) {
	interceptor := ErrorHandler()
	info := &grpc.UnaryServerInfo{FullMethod: "/proto.BookService/GetBook"}

	tests := []struct {
		name     string
		handler  grpc.UnaryHandler
		wantCode codes.Code
	}{
		{
			name: "success",
			handler: func(ctx context.Context, req interface{}) (interface{}, error) {
				return "ok", nil
			},
			wantCode: codes.OK,
		},
		{
			name: "not_found",
			handler: func(ctx context.Context, req interface{}) (interface{}, error) {
				return nil, domain.ErrBookNotFound
			},
			wantCode: codes.NotFound,
		},
		{
			name: "validation",
			handler: func(ctx context.Context, req interface{}) (interface{}, error) {
				return nil, domain.ErrTitleRequired
			},
			wantCode: codes.InvalidArgument,
		},
		{
			name: "internal",
			handler: func(ctx context.Context, req interface{}) (interface{}, error) {
				return nil, errors.New("db down")
			},
			wantCode: codes.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := interceptor(context.Background(), nil, info, tt.handler)
			if status.Code(err) != tt.wantCode {
				t.Fatalf("expected %s, got %s", tt.wantCode, status.Code(err))
			}
		})
	}
}

func TestRecovery(t *testing.T) {
	interceptor := Recovery()
	_, err := interceptor(
		context.Background(),
		nil,
		&grpc.UnaryServerInfo{FullMethod: "/proto.BookService/GetBook"},
		func(ctx context.Context, req interface{}) (interface{}, error) {
			panic("boom")
		},
	)
	if status.Code(err) != codes.Internal {
		t.Fatalf("expected Internal, got %s", status.Code(err))
	}
}
