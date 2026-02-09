// Package app initializes the application routers and dependencies
// for the every module.
package app

import (
	infraPostgres "github.com/james-wukong/orders-api/internal/infrastructure/postgres/persistence"
	"github.com/james-wukong/orders-api/internal/infrastructure/postgres/security"
	"github.com/james-wukong/orders-api/internal/interfaces/http/handlers"
	restaurantUC "github.com/james-wukong/orders-api/internal/usecase/restaurant"
	userUC "github.com/james-wukong/orders-api/internal/usecase/user"
	"gorm.io/gorm"
)

func (a *App) initRestaurantRouter(db *gorm.DB) *handlers.RestaurantHandler {
	// 1. Repository Layer: Infrastructure implementation of Domain interfaces ---
	// This variable satisfies the user.Repository interface
	repo := infraPostgres.NewRestaurantRepository(db)

	// 2. UseCase Layer: Business Logic ---
	// THIS IS HOW YOU CREATE THE createRestaurantUC VARIABLE
	// We pass the repository into the UseCase constructor
	createUC := restaurantUC.NewCreateRestaurantUseCase(repo)

	// 3. Handler Layer: HTTP Handlers ---
	return handlers.NewRestaurantHandler(createUC)
}

func (a *App) initUserRouter(db *gorm.DB) *handlers.UserHandler {
	// 1. Repository Layer: Infrastructure implementation of Domain interfaces ---
	repo := infraPostgres.NewUserRepository(db)

	// 2. Service Layer: Domain Services ---
	hasher := security.NewBcryptHasher(security.DefaultBcryptCost)

	// 3. UseCase Layer: Business Logic ---
	createUC := userUC.NewCreateUserUseCase(repo, hasher)

	// 4. Handler Layer: HTTP Handlers ---
	return handlers.NewUserHandler(createUC)
}
