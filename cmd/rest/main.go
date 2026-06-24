package main

import (
	"context"
	"log"

	"go-api-boilerplate/internal/adapter/repositories"
	"go-api-boilerplate/internal/adapter/rest"
	"go-api-boilerplate/internal/adapter/rest/book"
	"go-api-boilerplate/internal/application"
	"go-api-boilerplate/internal/config"
	"go-api-boilerplate/internal/infra"

	_ "go-api-boilerplate/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

// @title           API Demo
// @version         1.0
// @description     This is a sample server API Demo.

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /
type restApp struct {
	Router *gin.Engine
	db     *pgxpool.Pool
}

func initializeRestApp(ctx context.Context) (*restApp, error) {
	cfg := config.Config()

	db, err := infra.NewPostgresPool(ctx, cfg.Postgres(), cfg.Debug())
	if err != nil {
		return nil, err
	}

	bookRepo := repositories.NewPostgresBookRepo(db)
	bookService := application.NewBookService(bookRepo)
	bookHandler := book.NewBookHandler(bookService)

	router := gin.New()
	router.Use(gin.Recovery())
	if cfg.Debug() {
		router.Use(gin.Logger())
	}
	rest.SetupRouter(router, bookHandler)

	return &restApp{Router: router, db: db}, nil
}

func (a *restApp) Close() {
	a.db.Close()
}

func main() {
	app, err := initializeRestApp(context.Background())
	if err != nil {
		log.Fatalf("failed to create app: %v", err)
	}
	defer app.Close()

	app.Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	app.Router.Run(":8080")
}
