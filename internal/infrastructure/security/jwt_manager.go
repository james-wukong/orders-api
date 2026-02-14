package security

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTManager struct {
	secretKey string
	tokenTTL  time.Duration
}

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func NewJWTManager(secretKey string, tokenTTL time.Duration) *JWTManager {
	return &JWTManager{secretKey: secretKey, tokenTTL: tokenTTL}
}

// GenerateToken will create a jwt token based on user id
func (m *JWTManager) GenerateToken(userID uuid.UUID) (string, error) {
	claims := Claims{
		UserID: userID.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.tokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(m.secretKey))
}

func (m *JWTManager) VerifyToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(
		tokenStr,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(m.secretKey), nil
		},
	)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, err
	}

	return claims, nil
}
