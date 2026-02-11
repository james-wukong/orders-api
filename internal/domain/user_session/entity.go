// Package usersession defines the UserSessions entity for database interactions
package usersession

import (
	"time"

	"github.com/google/uuid"
	"github.com/james-wukong/orders-api/internal/pkg/utils"
)

const (
	TokenLength  = 32
	ExpireLength = 48 * time.Hour
)

// UserSessions represents the user_sessions table in the database
type UserSessions struct {
	ID              uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID          uuid.UUID `gorm:"type:uuid;not null"`
	Token           string    `gorm:"type:varchar(500);not null"`
	DeviceInfo      string    `gorm:"type:text"`
	IPAddress       string    `gorm:"type:varchar(45)"`
	OperatingSystem string    `gorm:"type:varchar(100)"`
	ClientDevice    string    `gorm:"type:varchar(100)"`
	UseAgent        string    `gorm:"type:varchar(100)"`
	ExpiresAt       time.Time `gorm:"not null"`
	CreatedAt       time.Time `gorm:"default:current_timestamp"`
}

// TableName overrides GORM's default pluralization
// func (UserSessions) TableName() string {
// 	return "user_sessions"
// }

func NewUserSession(userID uuid.UUID, deviceInfo string, ipAddress string, operatingSystem string, clientDevice string, userAgent string, expiresAt time.Time) *UserSessions {
	token, _ := utils.GenerateRandomToken(TokenLength) // Generate a random token for the session
	return &UserSessions{
		UserID:          userID,
		Token:           token,
		DeviceInfo:      deviceInfo,
		IPAddress:       ipAddress,
		OperatingSystem: operatingSystem,
		ClientDevice:    clientDevice,
		UseAgent:        userAgent,
		ExpiresAt:       time.Now().Add(ExpireLength),
	}
}
