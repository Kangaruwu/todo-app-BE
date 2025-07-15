package middlewares

import (
	"time"

	"go-backend-todo/internal/config"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// JWTManager handles JWT token operations
type JWTManager struct {
	cfg *config.Config
}

// NewJWTManager creates a new JWT manager
func NewJWTManager(cfg *config.Config) *JWTManager {
	return &JWTManager{cfg: cfg}
}

// JWTClaims represents JWT claims
type JWTClaims struct {
	UserID                string `json:"user_id"`
	Username              string `json:"username"`
	Email                 string `json:"email"`
	Role                  string `json:"role"`
	EmailValidationStatus string `json:"email_validation_status,omitempty"`
	jwt.RegisteredClaims
}

// GenerateAccessToken generates access token
func (j *JWTManager) GenerateAccessToken(userID uuid.UUID, username, email, role, emailStatus string) (string, error) {
	claims := JWTClaims{
		UserID:   userID.String(),
		Username: username,
		Email:    email,
		Role:     role,
		EmailValidationStatus: emailStatus, 
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(j.cfg.JWT.AccessExpiryHour))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    j.cfg.App.Name,
			Subject:   userID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.cfg.JWT.AccessSecret))
}

// GenerateRefreshToken generates refresh token
func (j *JWTManager) GenerateRefreshToken(userID uuid.UUID) (string, error) {
	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * time.Duration(j.cfg.JWT.RefreshExpiryDay))),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Subject:   userID.String(),
		Issuer:    j.cfg.App.Name,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.cfg.JWT.RefreshSecret))
}

// GenerateVerificationToken generates email verification token
func (j *JWTManager) GenerateVerificationToken(email string) (string, error) {
	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Subject:   email,
		Issuer:    j.cfg.App.Name,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.cfg.JWT.VerificationSecret))
}

// GenerateRecoveryToken generates password recovery token
func (j *JWTManager) GenerateRecoveryToken(email string) (string, error) {
	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Subject:   email,
		Issuer:    j.cfg.App.Name,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.cfg.JWT.RecoverySecret))
}

// ParseAccessToken parses and validates access token
func (j *JWTManager) ParseAccessToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.cfg.JWT.AccessSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}

// ParseRefreshToken parses and validates refresh token
func (j *JWTManager) ParseRefreshToken(tokenString string) (*jwt.RegisteredClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.cfg.JWT.RefreshSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}

// GetUserIDFromContext gets user ID from context
func GetUserIDFromContext(c *fiber.Ctx) (uuid.UUID, error) {
	userIDStr := c.Locals("user_id")
	if userIDStr == nil {
		return uuid.Nil, fiber.NewError(fiber.StatusUnauthorized, "User ID not found in context")
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		return uuid.Nil, fiber.NewError(fiber.StatusBadRequest, "Invalid user ID format")
	}

	return userID, nil
}
