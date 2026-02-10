package dto

import (
	"time"

	"github.com/google/uuid"
	session "github.com/james-wukong/orders-api/internal/domain/user_session"
	"github.com/james-wukong/orders-api/internal/pkg/utils"
)

// CreateSessionRequest defines how session data is sent over the API
type CreateSessionRequest struct {
	UserID          uuid.UUID `json:"user_id"`
	Token           string    `json:"token"`
	OperatingSystem string    `json:"operating_system"`
	ClientDevice    string    `json:"client_device"`
	UseAgent        string    `json:"use_agent"`
	DeviceInfo      string    `json:"device_info"`
	IPAddress       string    `json:"ip_address"`
	ExpiresAt       time.Time `json:"expires_at"`
}

func MapToCreateSessionRequest(entity *session.UserSessions) CreateSessionRequest {
	token, _ := utils.GenerateRandomToken(12)

	return CreateSessionRequest{
		UserID:          entity.ID,
		Token:           token,
		DeviceInfo:      entity.DeviceInfo,
		OperatingSystem: entity.OperatingSystem,
		ClientDevice:    entity.ClientDevice,
		UseAgent:        entity.UseAgent,
		IPAddress:       entity.IPAddress,
		ExpiresAt:       time.Now().Add(48 * time.Hour), // Set expiration to 48 hours from now
	}
}
