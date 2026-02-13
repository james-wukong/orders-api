package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/james-wukong/orders-api/internal/interfaces/http/dto"
	"github.com/james-wukong/orders-api/internal/interfaces/http/middleware"
	userUC "github.com/james-wukong/orders-api/internal/usecase/user"
)

type UserHandler struct {
	// createUserUC *user.CreateUserUseCase
	createUserUC *userUC.CreateUserUseCase
	loginUC      *userUC.LoginUseCase
}

func NewUserHandler(
	c *userUC.CreateUserUseCase,
) *UserHandler {
	return &UserHandler{
		createUserUC: c,
	}
}

// Register satisfies the RouterRegister interface
func (h *UserHandler) Register(mw *middleware.Manager, v1 *gin.RouterGroup) {
	userGroup := v1.Group("/user")
	{
		// userGroup.Use(mw.Authenticate())
		userGroup.POST("/register", h.Create)
		userGroup.POST("/login", h.Login)
		// userGroup.GET("/:id", mw.Authenticate(), h.GetProfile)
	}
}

// Create handles user registration requests
func (h *UserHandler) Create(c *gin.Context) {
	// 1. Map Request DTO to use case input
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2. Execute Use Case
	res, err := h.createUserUC.Execute(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 3. Map domain entity to response DTO
	c.JSON(http.StatusCreated, dto.MapToUserResponse(res))
}

func (h *UserHandler) Login(c *gin.Context) {
	// 1. Map Request DTO to use case input
	// var req dto.LoginRequest
	// if err := c.ShouldBindJSON(&req); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }

	// // 2. Execute Use Case
	// res, err := h.loginUC.Execute(c.Request.Context(), req)
	// if err != nil {
	// 	c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	// 	return
	// }

}
