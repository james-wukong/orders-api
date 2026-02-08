// Package handlers contains HTTP handlers for restaurant-related endpoints.
package handlers

import (
	"net/http"

	"github.com/james-wukong/orders-api/internal/interfaces/http/dto"
	"github.com/james-wukong/orders-api/internal/usecase/restaurant"

	"github.com/gin-gonic/gin"
)

type RestaurantHandler struct {
	createRestaurantUC *restaurant.CreateRestaurantUseCase
	// getRestaurantUC    *restaurant.GetRestaurantUseCase
}

func NewRestaurantHandler(
	c *restaurant.CreateRestaurantUseCase,
	// g *restaurant.GetRestaurantUseCase,
) *RestaurantHandler {
	return &RestaurantHandler{
		createRestaurantUC: c,
	}
}

// Register satisfies the RouterRegister interface
func (h *RestaurantHandler) Register(v1 *gin.RouterGroup) {
	userGroup := v1.Group("/restaurants")
	{
		userGroup.POST("/register", h.Create)
		// userGroup.GET("/:id", h.GetProfile)
	}
}

func (h *RestaurantHandler) Create(c *gin.Context) {
	var req dto.CreateRestaurantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.createRestaurantUC.Execute(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Map domain entity to response DTO
	c.JSON(http.StatusCreated, dto.MapToRestaurantResponse(res))
}
