// Package app initializes the application components such as database and cache clients.
package app

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/james-wukong/orders-api/internal/config"
	infraPostgre "github.com/james-wukong/orders-api/internal/infrastructure/postgres"
	infraRedis "github.com/james-wukong/orders-api/internal/infrastructure/redis"
	router "github.com/james-wukong/orders-api/internal/interfaces/http"
	"github.com/james-wukong/orders-api/internal/interfaces/http/middleware"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"gorm.io/gorm"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type App struct {
	HTTPServer *http.Server
	Database   *DBWrapper
	Redis      *redis.Client
	Log        *zerolog.Logger
	AppConfig  *config.AppConfig
	JWTConfig  *config.JWTConfig
}

type DBWrapper struct {
	DB *gorm.DB
}

type DBPoolWrapper struct {
	Pool *pgxpool.Pool
}

func Bootstrap(ctx context.Context, log *zerolog.Logger) (*App, error) {
	cfg := config.InitConfig()

	// Initialize Postgres connection pool
	db, err := infraPostgre.NewGormDB(ctx, cfg.Database.Postgres)
	if err != nil {
		return nil, err
	}
	// Initialize Redis client
	redisClient, err := infraRedis.New(cfg.Caches.Redis)
	if err != nil {
		return nil, err
	}

	// 1. Setup Gin and middleware
	r := gin.Default()
	mw := middleware.NewManager(log, db, redisClient)
	if cfg.App.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	r.Use(
		mw.RateLimiterMiddleware(),
		mw.CORSMiddleware(),
		mw.RecoveryMiddleware(),
		mw.SetClientMiddleware(),
	)
	v1 := r.Group("/api/v1")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// setup http server
	server := &http.Server{
		Addr:           fmt.Sprintf(":%d", cfg.App.Port),
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	application := &App{
		HTTPServer: server,
		Database:   &DBWrapper{DB: db},
		Redis:      redisClient,
		Log:        log,
		AppConfig:  &cfg.App,
		JWTConfig:  &cfg.JWT,
	}

	// 2. Init Handlers
	rHandler := application.initRestaurantRouter()
	uHandler := application.initUserRouter()

	// 3. Register everything dynamically
	routerManager := router.NewRouter(r, mw)
	routerManager.RegisterModules(v1,
		rHandler,
		uHandler,
	// Adding a new module (e.g. PaymentHandler) is now just one line here!
	)

	return application, nil
}
