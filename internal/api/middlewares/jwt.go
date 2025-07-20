package middlewares

import (
	"context"
	"log"
	"time"

	"go-backend-todo/internal/config"
	user_repository "go-backend-todo/internal/repository/user"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// JWTManager handles JWT token operations
type JWTManager struct {
	cfg      *config.Config
	userRepo user_repository.UserRepository
}

// NewJWTManager creates a new JWT manager
func NewJWTManager(cfg *config.Config, userRepo user_repository.UserRepository) *JWTManager {
	return &JWTManager{
		cfg:      cfg,
		userRepo: userRepo,
	}
}

// JWTClaims represents JWT claims
type JWTClaims struct {
	UserID                string `json:"user_id"`
	Username              string `json:"username"`
	Email                 string `json:"email"`
	Role                  string `json:"role"`
	EmailValidationStatus string `json:"email_validation_status,omitempty"`
	TokenVersion          int    `json:"token_version"`
	jwt.RegisteredClaims
}

// GenerateAccessToken generates access token
func (j *JWTManager) GenerateAccessToken(userID uuid.UUID, username, email, role, emailStatus string, tokenVersion int) (string, error) {
	claims := JWTClaims{
		UserID:                userID.String(),
		Username:              username,
		Email:                 email,
		Role:                  role,
		EmailValidationStatus: emailStatus,
		TokenVersion:          tokenVersion,
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
		// Validate token version
		userID, err := uuid.Parse(claims.UserID)
		if err != nil {
			return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid user ID in token")
		}

		// Get current token version from database
		currentVersion, err := j.userRepo.GetTokenVersion(context.Background(), userID)
		if err != nil {
			return nil, fiber.NewError(fiber.StatusUnauthorized, "Failed to validate token")
		}

		// Check if token version matches current version
		if claims.TokenVersion != currentVersion {
			log.Println("Token version mismatch:", claims.TokenVersion, "expected:", currentVersion)
			return nil, fiber.NewError(fiber.StatusUnauthorized, "Token has been revoked")
		}

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
