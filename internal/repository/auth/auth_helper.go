package auth_repository

import (
	"context"
	_ "embed"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/resend/resend-go/v2"
	"go-backend-todo/internal/config"
	"go-backend-todo/internal/utils"
	"golang.org/x/crypto/bcrypt"
	"go-backend-todo/internal/models"
)

//go:embed assets/ConfirmationEmail.html
var confirmationEmailTemplate string

func SendVerificationEmail(ctx context.Context, username, email string, confirmationToken string) error {
	client := resend.NewClient(config.GetEnv("RESEND_API_KEY", ""))

	params := &resend.SendEmailRequest{
		From:    config.GetEnv("MY_EMAIL", ""),
		To:      []string{email},
		Subject: "[The Ultimate Todo] Email Confirmation",
		Html:    GenerateEmailValidationHTML(username, GenerateEmailValidationLink(confirmationToken)),
	}

	sent, err := client.Emails.SendWithContext(ctx, params)

	if err != nil {
		return err
	}
	fmt.Println(sent.Id)

	return nil
}

func GenerateEmailValidationToken() (string, error) {
	guuid := uuid.New()
	if guuid == uuid.Nil {
		return "", fmt.Errorf("failed to generate UUID for email validation token")
	}
	salt := utils.RandInRange(bcrypt.MinCost, bcrypt.MaxCost)
	confirmationToken, err := bcrypt.GenerateFromPassword([]byte(guuid.String()), salt)
	if err != nil {
		return "", fmt.Errorf("failed to generate password hash for email validation token: %w", err)
	}

	return string(confirmationToken), nil
}

func GenerateEmailValidationLink(token string) string {
	baseURL := config.GetEnv("BASE_URL", "http://localhost:8080")
	return fmt.Sprintf("%s/confirm-email?token=%s", baseURL, token)
}

func GenerateEmailValidationHTML(username, activationLink string) string {
	return fmt.Sprintf(confirmationEmailTemplate,
		username,
		activationLink,
		activationLink,
		2025,
	)
}

func GenerateAccessToken(userID uuid.UUID, username, email, role string) (string, error) {
	claims := models.JWTClaims{
		UserID:   userID.String(),
		Username: username,
		Email:    email,
		Role:     role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(), 
			IssuedAt:  time.Now().Unix(),
			NotBefore: time.Now().Unix(),
			Issuer:    "go-backend-todo",
			Subject:   userID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(config.GetEnv("JWT_ACCESS_SECRET", "sfsdfsdfsdfsdfsd")))
	if err != nil {
		return "", fmt.Errorf("failed to generate access token: %w", err)
	}
	return signedToken, nil
}

func GenerateRefreshToken(userID uuid.UUID, username, email, role string) (string, error) {
	claims := models.JWTClaims{
		UserID:   userID.String(),
		Username: username,
		Email:    email,
		Role:     role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(), 
			IssuedAt:  time.Now().Unix(),
			NotBefore: time.Now().Unix(),
			Issuer:    "go-backend-todo",
			Subject:   userID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(config.GetEnv("JWT_REFRESH_SECRET", "bruh")))
	if err != nil {
		return "", fmt.Errorf("failed to generate refresh token: %w", err)
	}
	return signedToken, nil
}