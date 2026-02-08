// Package app initializes the application routers and dependencies
// for the every module.
package app

import (
	infraPostgres "github.com/james-wukong/orders-api/internal/infrastructure/postgres"
	"github.com/james-wukong/orders-api/internal/interfaces/http/handlers"
	restaurantUC "github.com/james-wukong/orders-api/internal/usecase/restaurant"
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

	return handlers.NewRestaurantHandler(createUC)
}
