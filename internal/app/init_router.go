// Package app initializes the application routers and dependencies
// for the every module.
package app

import (
	infraPostgres "github.com/james-wukong/orders-api/internal/infrastructure/postgres/persistence"
	infraRedis "github.com/james-wukong/orders-api/internal/infrastructure/redis/cache"
	infraCache "github.com/james-wukong/orders-api/internal/infrastructure/repository"
	"github.com/james-wukong/orders-api/internal/infrastructure/security"
	"github.com/james-wukong/orders-api/internal/interfaces/http/handlers"
	restaurantUC "github.com/james-wukong/orders-api/internal/usecase/restaurant"
	userUC "github.com/james-wukong/orders-api/internal/usecase/user"
)

func (a *App) initRestaurantRouter() *handlers.RestaurantHandler {
	// 1. Repository Layer: Infrastructure implementation of Domain interfaces ---
	// This variable satisfies the user.Repository interface
	redisRepo := infraRedis.NewRestaurantCache(a.Redis, a.Log)
	repo := infraPostgres.NewRestaurantRepository(a.Database.DB, a.Log)
	cachedRepo := infraCache.NewCachedRestaurantRepository(repo, redisRepo, a.Log)

	// 2. UseCase Layer: Business Logic ---
	// THIS IS HOW YOU CREATE THE createRestaurantUC VARIABLE
	// We pass the repository into the UseCase constructor
	createUC := restaurantUC.NewCreateRestaurantUseCase(cachedRepo)

	// 3. Handler Layer: HTTP Handlers ---
	return handlers.NewRestaurantHandler(createUC)
}

func (a *App) initUserRouter() *handlers.UserHandler {
	// 1. Repository Layer: Infrastructure implementation of Domain interfaces ---
	redisRepo := infraRedis.NewUserCache(a.Redis, a.Log)
	repo := infraPostgres.NewUserRepository(a.Database.DB, a.Log)
	cachedRepo := infraCache.NewCachedUserRepository(repo, redisRepo, a.Log)

	// 2. Service Layer: Domain Services ---
	hasher := security.NewBcryptHasher(security.DefaultBcryptCost)

	// 3. UseCase Layer: Business Logic ---
	createUC := userUC.NewCreateUserUseCase(cachedRepo, hasher)

	// 4. Handler Layer: HTTP Handlers ---
	return handlers.NewUserHandler(createUC)
}
