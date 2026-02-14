package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/james-wukong/orders-api/internal/interfaces/http/dto"
	"github.com/james-wukong/orders-api/internal/interfaces/http/middleware"
	"github.com/james-wukong/orders-api/internal/pkg/utils"
	userUC "github.com/james-wukong/orders-api/internal/usecase/user"
)

type UserHandler struct {
	// createUserUC *user.CreateUserUseCase
	createUserUC *userUC.CreateUserUseCase
	loginUC      *userUC.LoginUseCase
}

func NewUserHandler(
	c *userUC.CreateUserUseCase,
	l *userUC.LoginUseCase,
) *UserHandler {
	return &UserHandler{
		createUserUC: c,
		loginUC:      l,
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

// Create godoc
//
//	@Summary		Create
//	@Description	register a user
//	@Tags			register users
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.CreateUserRequest	true	"User Credentials"
//	@Success		201		{object}	dto.UserResponse
//	@Failure		400		{object}	map[string]any	"{'error':'error message'}"
//	@Failure		500		{object}	map[string]any	"{'error':'error message'}"
//	@Router			/user/register [post]
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

// Login godoc
//
//	@Summary		Login
//	@Description	login a user
//	@Tags			login
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.LoginRequest	true	"User Credentials"
//	@Success		200		{object}	dto.LoginResponse
//	@Failure		400		{object}	map[string]any	"{'error':'error message'}"
//	@Failure		401		{object}	map[string]any	"{'error':'error message'}"
//	@Router			/user/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req dto.LoginRequest

	// 1. Map request to DTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// get device info from request
	extractor := utils.NewMetadataExtractor()
	metadata := extractor.GetMetadata(c)
	mJSONBytes, err := json.Marshal(metadata.Location)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Login user and get token
	resp, err := h.loginUC.Execute(c.Request.Context(),
		&req,
		metadata.IP, metadata.OS,
		metadata.Device, metadata.UserAgent,
		string(mJSONBytes),
	)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid credentials",
		})
		return
	}
	// Map to response
	c.JSON(http.StatusOK, dto.MapToLoginResponse(resp))

}
