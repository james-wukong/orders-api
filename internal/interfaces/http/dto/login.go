package dto

import "time"

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token       string      `json:"token"`
	Email       string      `json:"email"`
	ExpiresAt   time.Time   `json:"expires_at"`
	SessionInfo SessionInfo `json:"session_info"`
}

type SessionInfo struct {
	OperatingSystem string `json:"operating_system"`
	ClientDevice    string `json:"client_device"`
	UseAgent        string `json:"use_agent"`
	// DeviceInfo      string `json:"device_info"`
	IPAddress string `json:"ip_address"`
}
