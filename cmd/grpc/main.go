package main

import (
	"context"
	"fmt"
	"log"
	"net"

	grpcadapter "go-api-boilerplate/internal/adapter/grpc"
	"go-api-boilerplate/internal/adapter/grpc/book"
	"go-api-boilerplate/internal/adapter/repositories/postgres"
	"go-api-boilerplate/internal/application"
	"go-api-boilerplate/internal/config"
	"go-api-boilerplate/internal/infra"

	"github.com/jackc/pgx/v5/pgxpool"
	grpclib "google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type grpcApp struct {
	server *grpclib.Server
	db     *pgxpool.Pool
}

func initializeGrpcApp(ctx context.Context) (*grpcApp, error) {
	cfg := config.Config()

	db, err := infra.NewPostgresPool(ctx, cfg.Postgres(), cfg.Debug())
	if err != nil {
		return nil, err
	}

	bookRepo := postgres.NewPostgresBookRepo(db)
	bookService := application.NewBookService(bookRepo)
	bookServer := book.NewBookServer(bookService)

	server := grpcadapter.NewServer(bookServer, cfg.Debug())
	if cfg.Debug() {
		reflection.Register(server)
	}

	return &grpcApp{server: server, db: db}, nil
}

func (a *grpcApp) Close() {
	a.db.Close()
}

func main() {
	cfg := config.Config()

	app, err := initializeGrpcApp(context.Background())
	if err != nil {
		log.Fatalf("failed to create app: %v", err)
	}
	defer app.Close()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Grpc().Port))
	if err != nil {
		log.Fatalf("failed to listen on gRPC port: %v", err)
	}

	log.Printf("gRPC server starting on :%d", cfg.Grpc().Port)
	if err := app.server.Serve(lis); err != nil {
		log.Fatalf("failed to serve gRPC: %v", err)
	}
}
