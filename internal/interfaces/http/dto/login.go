package dto

import (
	"time"

	"github.com/google/uuid"
	uSession "github.com/james-wukong/orders-api/internal/domain/user_session"
)

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token       string       `json:"token"`
	SessionInfo *SessionInfo `json:"session_info"`
}

type SessionInfo struct {
	UserID          uuid.UUID `json:"user_id"`
	OperatingSystem string    `json:"operating_system"`
	ClientDevice    string    `json:"client_device"`
	UserAgent       string    `json:"user_agent"`
	IPAddress       string    `json:"ip_address"`
	ExpiresAt       time.Time `json:"expires_at"`
}

func MapToLoginResponse(entity *uSession.UserSessions) *LoginResponse {
	return &LoginResponse{
		Token: entity.Token,
		SessionInfo: &SessionInfo{
			UserID:          entity.UserID,
			OperatingSystem: entity.OperatingSystem,
			ClientDevice:    entity.ClientDevice,
			UserAgent:       entity.UserAgent,
			IPAddress:       entity.IPAddress,
			ExpiresAt:       time.Now().Add(uSession.TokenTTL),
		},
	}
}
