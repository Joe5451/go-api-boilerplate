package interceptors

import (
	"context"
	"log"

	"go-api-boilerplate/internal/adapter/grpc/grpckit"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ErrorHandler() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		resp, err := handler(ctx, req)
		if err == nil {
			return resp, nil
		}

		st := grpckit.ToStatus(err)
		if status.Code(st) == codes.Internal {
			log.Printf("[GRPC_INTERNAL_ERROR] %s: %+v\n", info.FullMethod, err)
		}

		return nil, st
	}
}
