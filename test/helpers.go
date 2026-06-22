package test

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"go-api-boilerplate/internal/adapter/repositories"
	"go-api-boilerplate/internal/adapter/rest"
	"go-api-boilerplate/internal/adapter/rest/book"
	"go-api-boilerplate/internal/application"
	"go-api-boilerplate/internal/config"
	"go-api-boilerplate/internal/infra"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
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

	bookRepo := repositories.NewPostgresBookRepo(db)
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
