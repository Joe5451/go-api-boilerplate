package grpc

import (
	"go-api-boilerplate/internal/adapter/grpc/book"
	"go-api-boilerplate/internal/adapter/grpc/interceptors"
	"go-api-boilerplate/proto"

	grpclib "google.golang.org/grpc"
)

func NewServer(bookServer *book.BookServer, debug bool) *grpclib.Server {
	unaryInterceptors := []grpclib.UnaryServerInterceptor{
		interceptors.Recovery(),
		interceptors.ErrorHandler(),
	}
	if debug {
		unaryInterceptors = []grpclib.UnaryServerInterceptor{
			interceptors.Recovery(),
			interceptors.Logging(),
			interceptors.ErrorHandler(),
		}
	}

	s := grpclib.NewServer(grpclib.ChainUnaryInterceptor(unaryInterceptors...))
	proto.RegisterBookServiceServer(s, bookServer)
	return s
}
