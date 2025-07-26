package handlers

import (
	"go-backend-todo/internal/api/middlewares"
	"go-backend-todo/internal/api/responses"
	"go-backend-todo/internal/models"
	"go-backend-todo/internal/service"

	"log"

	"github.com/gofiber/fiber/v2"
)

// AuthHandler handles authentication HTTP requests
type AuthHandler struct {
	authService service.AuthService
	jwtManager  *middlewares.JWTManager
}

// NewAuthHandler creates a new instance of auth handler
func NewAuthHandler(authService service.AuthService, jwtManager *middlewares.JWTManager) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		jwtManager:  jwtManager,
	}
}

// Login handles user login
// @Summary User login
// @Description Authenticate user with email and password
// @Tags Authentication
// @Accept json
// @Produce json
// @Param credentials body models.LoginRequest true "Login credentials"
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req models.LoginRequest
	body := c.Body()
	if err := middlewares.RequestValidation(&body, &req)(c); err != nil {
		return responses.BadRequestWithError(c, "Invalid request body", err)
	}

	user, err := h.authService.Login(c.Context(), &req)
	if err != nil {
		return responses.Unauthorized(c, "Invalid email or password")
	}

	// Generate tokens
	accessToken, err := h.jwtManager.GenerateAccessToken(user.UserID, user.Username, user.Email, string(user.Role), string(user.Status), user.TokenVersion)
	if err != nil {
		return responses.InternalServerError(c, "Failed to generate access token")
	}

	refreshToken, err := h.jwtManager.GenerateRefreshToken(user.UserID, user.TokenVersion)
	if err != nil {
		return responses.InternalServerError(c, "Failed to generate refresh token")
	}

	response := models.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return responses.OK(c, "Login successful", response)
}

// Register handles user registration
// @Summary User registration
// @Description Register a new user account
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body models.RegisterRequest true "User registration data"
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	body := c.Body()

	if len(body) == 0 {
		return responses.BadRequest(c, "Request body cannot be empty")
	}

	var req models.RegisterRequest
	err := middlewares.RequestValidation(&body, &req)(c)
	if err != nil {
		return responses.BadRequestWithError(c, "Invalid request body", err)
	}

	verificationToken, err := h.jwtManager.GenerateVerificationToken(req.Email)
	if err != nil {
		return responses.InternalServerError(c, "Failed to generate verification token: "+err.Error())
	}

	if err := h.authService.Register(c.Context(), &req, verificationToken); err != nil {
		return responses.InternalServerError(c, "Failed to register user: "+err.Error())
	}

	log.Println("User registered successfully:", req.Username)

	return responses.Created(c, "User registered successfully", nil)
}

// VerifyEmail handles email verification
// @Summary Verify email address
// @Description Verify user email with verification token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param token path string true "Email verification token"
// @Router /auth/verify-email/{token} [get]
func (h *AuthHandler) VerifyEmail(c *fiber.Ctx) error {
	token := c.Params("token")
	if token == "" {
		return responses.BadRequest(c, "Token is required")
	}

	if err := h.authService.VerifyEmail(c.Context(), token); err != nil {
		return responses.InternalServerError(c, "Failed to verify email: "+err.Error())
	}

	return responses.OK(c, "Email verified successfully", nil)
}

// RecoverPassword handles password recovery
// @Summary Password recovery
// @Description Send password reset email
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body models.RecoverPasswordRequest true "Recovery email data"
// @Router /auth/recover-password [post]
func (h *AuthHandler) RecoverPassword(c *fiber.Ctx) error {
	var req models.RecoverPasswordRequest
	body := c.Body()
	if err := middlewares.RequestValidation(&body, &req)(c); err != nil {
		return responses.BadRequestWithError(c, "Invalid request body", err)
	}

	recoverToken, err := h.jwtManager.GenerateRecoveryToken(req.Email)
	if err != nil {
		return responses.InternalServerError(c, "Failed to generate recovery token: "+err.Error())
	}

	if err := h.authService.RecoverPassword(c.Context(), &req, recoverToken); err != nil {
		return responses.InternalServerError(c, "Failed to send recovery email: "+err.Error())
	}

	return responses.OK(c, "Recovery email sent successfully", nil)
}

// ResetPassword handles password reset
// @Summary Reset password
// @Description Reset user password with recovery token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body models.ResetPasswordRequest true "Password reset data"
// @Router /auth/reset-password [post]
func (h *AuthHandler) ResetPassword(c *fiber.Ctx) error {
	var req models.ResetPasswordRequest
	body := c.Body()
	if err := middlewares.RequestValidation(&body, &req)(c); err != nil {
		return responses.BadRequestWithError(c, "Invalid request body", err)
	}
	if req.Token == "" {
		return responses.BadRequest(c, "Recovery token is required")
	}
	if req.NewPassword == "" {
		return responses.BadRequest(c, "New password is required")
	}
	if err := h.authService.ResetPassword(c.Context(), &req); err != nil {
		return responses.InternalServerError(c, "Failed to reset password: "+err.Error())
	}
	return responses.OK(c, "Password reset successfully", nil)
}

// Logout handles user logout
// @Summary User logout
// @Description Log out the currently authenticated user
// @Tags Authentication
// @Accept json
// @Produce json
// @Security BearerAuth
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	c.Locals("user_id", nil)
	c.Locals("username", nil)
	c.Locals("email", nil)
	c.Locals("role", nil)
	c.Locals("email_validation_status", nil)
	c.Locals("claims", nil)
	return responses.OK(c, "Logged out successfully", nil)
}

// RefreshAccessToken handles access token refresh
// @Summary Refresh access token
// @Description Refresh the access token for the currently authenticated user
// @Tags Authentication
// @Accept json
// @Produce json
// @Param token body models.RefreshAccessTokenRequest true "Refresh token data"
// @Router /auth/refresh-token [post]
func (h *AuthHandler) RefreshAccessToken(c *fiber.Ctx) error {
	var req models.RefreshAccessTokenRequest
	body := c.Body()
	if err := middlewares.RequestValidation(&body, &req)(c); err != nil {
		return responses.BadRequestWithError(c, "Invalid request body", err)
	}
	if req.RefreshToken == "" {
		return responses.BadRequest(c, "Refresh token is required")
	}
	
	newAccessToken, newRefreshToken, err := h.jwtManager.RefreshAccessToken(req.RefreshToken)
	if err != nil {
		return responses.InternalServerError(c, "Failed to refresh access token: "+err.Error())
	}
	
	response := models.RefreshAccessTokenResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}

	return responses.OK(c, "Access token refreshed successfully", response)
}