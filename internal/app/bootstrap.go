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
	"github.com/james-wukong/orders-api/internal/infrastructure/logger"
	infraPostgre "github.com/james-wukong/orders-api/internal/infrastructure/postgres"
	infraRedis "github.com/james-wukong/orders-api/internal/infrastructure/redis"
	router "github.com/james-wukong/orders-api/internal/interfaces/http"
	"github.com/james-wukong/orders-api/internal/interfaces/http/middleware"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type App struct {
	HTTPServer *http.Server
	Database   *DBWrapper
	Redis      *redis.Client
}

type DBWrapper struct {
	DB *gorm.DB
}

type DBPoolWrapper struct {
	Pool *pgxpool.Pool
}

// Initialize logger with console output only
var conLog = logger.New(logger.LogConfig{
	EnableConsole: true,
	FilePath:      "logs/app.log",
	MaxSize:       5,     // Rotate every 5MB
	MaxBackups:    10,    // Keep last 10 files
	Compress:      false, // Save disk space
})

func Bootstrap(ctx context.Context) (*App, error) {
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

	// 1. Setup Gin
	r := gin.Default()
	if cfg.App.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	r.Use(
		middleware.RateLimiterMiddleware(),
		middleware.CORSMiddleware(),
		middleware.RecoveryMiddleware(),
		middleware.SetClientMiddleware(),
	)
	v1 := r.Group("/api/v1")

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
	}

	// 2. Init Handlers
	rHandler := application.initRestaurantRouter(db)

	// 3. Register everything dynamically
	routerManager := router.NewRouter(r)
	routerManager.RegisterModules(v1,
		// uHandler,
		// oHandler,
		rHandler,
	// Adding a new module (e.g. PaymentHandler) is now just one line here!
	)

	return application, nil
}
