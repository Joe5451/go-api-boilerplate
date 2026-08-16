package test

import (
	"context"
	"fmt"
	"log"
	"net"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	grpcadapter "go-api-boilerplate/internal/adapter/grpc"
	grpcbook "go-api-boilerplate/internal/adapter/grpc/book"
	pg "go-api-boilerplate/internal/adapter/repositories/postgres"
	"go-api-boilerplate/internal/adapter/rest"
	"go-api-boilerplate/internal/adapter/rest/book"
	"go-api-boilerplate/internal/application"
	"go-api-boilerplate/internal/config"
	"go-api-boilerplate/internal/infra"
	"go-api-boilerplate/proto"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	grpclib "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

var (
	dbPool *pgxpool.Pool
)

type RestApp struct {
	Router *gin.Engine
	db     *pgxpool.Pool
}

func Setup() {
	ctx := context.Background()

	pgContainer := createPostgresContainer(ctx)
	pgHost, err := pgContainer.Host(ctx)
	if err != nil {
		log.Fatalf("failed to get postgres container host: %v", err)
	}

	pgPort, err := pgContainer.MappedPort(ctx, "5432")
	if err != nil {
		log.Fatalf("failed to get postgres container port: %v", err)
	}

	dbPool, err = pgxpool.New(ctx, fmt.Sprintf(
		"postgres://testuser:testpassword@%s:%s/testdb?sslmode=disable",
		pgHost,
		pgPort.Port(),
	))
	if err != nil {
		log.Fatalf("failed to create postgres pool: %v", err)
	}

	viper.Set("DEBUG", false)
	viper.Set("POSTGRES_HOST", pgHost)
	viper.Set("POSTGRES_PORT", pgPort.Port())
	viper.Set("POSTGRES_USER", "testuser")
	viper.Set("POSTGRES_PASSWORD", "testpassword")
	viper.Set("POSTGRES_DBNAME", "testdb")
	viper.Set("POSTGRES_SCHEMA", "public")
}

func InitRestApp(ctx context.Context) (*RestApp, error) {
	cfg := config.Config()

	db, err := infra.NewPostgresPool(ctx, cfg.Postgres(), cfg.Debug())
	if err != nil {
		return nil, err
	}

	bookRepo := pg.NewPostgresBookRepo(db)
	bookService := application.NewBookService(bookRepo)
	bookHandler := book.NewBookHandler(bookService)

	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(gin.Recovery())
	if cfg.Debug() {
		router.Use(gin.Logger())
	}
	rest.SetupRouter(router, bookHandler)

	return &RestApp{Router: router, db: db}, nil
}

func (a *RestApp) Close() {
	a.db.Close()
}

type GrpcApp struct {
	Client proto.BookServiceClient
	conn   *grpclib.ClientConn
	server *grpclib.Server
	db     *pgxpool.Pool
}

func InitGrpcApp(ctx context.Context) (*GrpcApp, error) {
	cfg := config.Config()

	db, err := infra.NewPostgresPool(ctx, cfg.Postgres(), cfg.Debug())
	if err != nil {
		return nil, err
	}

	bookRepo := pg.NewPostgresBookRepo(db)
	bookService := application.NewBookService(bookRepo)
	bookServer := grpcbook.NewBookServer(bookService)

	lis := bufconn.Listen(1024 * 1024)
	server := grpcadapter.NewServer(bookServer, false)
	go func() {
		if err := server.Serve(lis); err != nil {
			log.Printf("bufconn server stopped: %v", err)
		}
	}()

	conn, err := grpclib.NewClient(
		"passthrough:///bufnet",
		grpclib.WithContextDialer(func(ctx context.Context, _ string) (net.Conn, error) {
			return lis.DialContext(ctx)
		}),
		grpclib.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		server.Stop()
		db.Close()
		return nil, err
	}

	return &GrpcApp{
		Client: proto.NewBookServiceClient(conn),
		conn:   conn,
		server: server,
		db:     db,
	}, nil
}

func (a *GrpcApp) Close() {
	a.server.Stop()
	_ = a.conn.Close()
	a.db.Close()
}

func Teardown() {
	if dbPool != nil {
		dbPool.Close()
	}

	// TODO: Terminate container
}

func DB() *pgxpool.Pool {
	return dbPool
}

func CleanupDatabase(t *testing.T) {
	if dbPool == nil {
		t.Fatal("database pool not initialized")
	}

	_, err := dbPool.Exec(
		context.Background(),
		"TRUNCATE books RESTART IDENTITY CASCADE",
	)
	if err != nil {
		t.Logf("warning: failed to truncate: %v", err)
	}
}

func createPostgresContainer(ctx context.Context) *postgres.PostgresContainer {
	pgContainer, err := postgres.Run(ctx,
		"postgres:15.3-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpassword"),
		postgres.WithInitScripts(getInitSQLPath()),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second),
		),
	)
	if err != nil {
		log.Fatalf("failed to start postgres container: %v", err)
	}
	return pgContainer
}

func getInitSQLPath() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(filename), "..", "init.sql")
}
