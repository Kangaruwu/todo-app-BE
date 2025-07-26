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

// RefreshJWTClaims represents JWT claims for refresh token
type RefreshJWTClaims struct {
	TokenVersion int `json:"token_version"`
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
func (j *JWTManager) GenerateRefreshToken(userID uuid.UUID, tokenVersion int) (string, error) {
	claims := JWTClaims{
		TokenVersion: tokenVersion,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * time.Duration(j.cfg.JWT.RefreshExpiryDay))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   userID.String(),
			Issuer:    j.cfg.App.Name,
		},
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
func (j *JWTManager) ParseRefreshToken(tokenString string) (*RefreshJWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &RefreshJWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.cfg.JWT.RefreshSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*RefreshJWTClaims); ok && token.Valid {
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


// RefreshAccessToken refreshes the access token, also returning new refresh token
// This function rotates the refresh token on every refresh
func (j *JWTManager) RefreshAccessToken(refreshToken string) (string, string, error) {
	// Parse the refresh token
	claims, err := j.ParseRefreshToken(refreshToken)
	if err != nil {
		return "", "", fiber.NewError(fiber.StatusUnauthorized, "Invalid refresh token: "+err.Error())
	}

	// Get user ID from claims
	userID := claims.Subject

	// Check token version
	currentVersion, err := j.userRepo.GetTokenVersion(context.Background(), uuid.MustParse(userID))
	if err != nil {
		return "", "", fiber.NewError(fiber.StatusUnauthorized, "Failed to get current token version: "+err.Error())
	}
	if currentVersion != claims.TokenVersion {
		return "", "", fiber.NewError(fiber.StatusUnauthorized, "Token has been revoked or version mismatch")
	}

	// Increment token version in the database
	err = j.userRepo.IncrementTokenVersion(context.Background(), uuid.MustParse(userID))
	if err != nil {
		return "", "", fiber.NewError(fiber.StatusUnauthorized, "Failed to increment token version: "+err.Error())
	}

	user, err := j.userRepo.GetByID(context.Background(), uuid.MustParse(userID))
	if err != nil {
		return "", "", fiber.NewError(fiber.StatusUnauthorized, "Failed to get user: "+err.Error())
	}

	// Generate new access token
	newAccessToken, err := j.GenerateAccessToken(uuid.MustParse(userID), claims.Subject, claims.Subject, string(user.Role), string(user.Status), user.TokenVersion)
	if err != nil {
		return "", "", fiber.NewError(fiber.StatusInternalServerError, "Failed to generate new access token: "+err.Error())
	}

	// Generate new refresh token (Refresh token rotates on every refresh)
	newRefreshToken, err := j.GenerateRefreshToken(uuid.MustParse(userID), user.TokenVersion)
	if err != nil {
		return "", "", fiber.NewError(fiber.StatusInternalServerError, "Failed to generate new refresh token: "+err.Error())
	}

	return newAccessToken, newRefreshToken, nil
}