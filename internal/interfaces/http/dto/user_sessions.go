package dto

import (
	"time"

	"github.com/google/uuid"
)

// CreateSessionRequest defines how session data is sent over the API
type CreateSessionRequest struct {
	UserID          uuid.UUID         `json:"user_id"`
	Token           string            `json:"token"`
	OperatingSystem string            `json:"operating_system"`
	ClientDevice    string            `json:"client_device"`
	UserAgent       string            `json:"use_agent"`
	DeviceInfo      map[string]string `json:"device_info"`
	IPAddress       string            `json:"ip_address"`
	ExpiresAt       time.Time         `json:"expires_at"`
}
