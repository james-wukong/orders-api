package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/james-wukong/orders-api/internal/interfaces/http/dto"
	userUC "github.com/james-wukong/orders-api/internal/usecase/user"
)

type UserHandler struct {
	// createUserUC *user.CreateUserUseCase
	createUserUC *userUC.CreateUserUseCase
}

func NewUserHandler(
	c *userUC.CreateUserUseCase,
) *UserHandler {
	return &UserHandler{
		createUserUC: c,
	}
}

// Register satisfies the RouterRegister interface
func (h *UserHandler) Register(v1 *gin.RouterGroup) {
	userGroup := v1.Group("/user")
	{
		userGroup.POST("/register", h.Create)
		// userGroup.GET("/:id", h.GetProfile)
	}
}

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
