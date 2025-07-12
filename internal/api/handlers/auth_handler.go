package handlers

import (
	"encoding/json"
	"go-backend-todo/internal/api/middlewares"
	"go-backend-todo/internal/api/responses"
	"go-backend-todo/internal/models"
	"go-backend-todo/internal/service"

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
// @Success 200 {object} models.LoginResponse "Login successful"
// @Failure 400 {object} map[string]interface{} "Invalid request data"
// @Failure 401 {object} map[string]interface{} "Invalid credentials"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req models.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return responses.BadRequest(c, "Invalid request body format")
	}

	// Validate request
	if err := middlewares.ValidateStruct(&req); err != nil {
		return responses.BadRequestWithError(c, "Validation failed", err)
	}

	user, err := h.authService.Login(c.Context(), &req)
	if err != nil {
		return responses.Unauthorized(c, "Invalid email or password")
	}

	// Generate tokens
	accessToken, err := h.jwtManager.GenerateAccessToken(user.UserID, user.Username, user.Email, string(user.Role))
	if err != nil {
		return responses.InternalServerError(c, "Failed to generate access token")
	}

	refreshToken, err := h.jwtManager.GenerateRefreshToken(user.UserID)
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
// @Success 201 {object} models.RegisterResponse
// @Failure 400 {object} responses.ErrorResponse
// @Failure 409 {object} responses.ErrorResponse
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	body := c.Body()

	if len(body) == 0 {
		return responses.BadRequest(c, "Request body cannot be empty")
	}

	var req models.RegisterRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return responses.BadRequest(c, "Invalid request body")
	}

	if err := h.authService.Register(c.Context(), &req); err != nil {
		return responses.InternalServerError(c, "Failed to register user")
	}

	return responses.Created(c, "User registered successfully", nil)
}

// ConfirmEmail handles email confirmation
// @Summary Confirm email address
// @Description Confirm user email with verification token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param token path string true "Email confirmation token"
// @Success 200 {object} responses.SuccessResponse
// @Failure 400 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Router /auth/confirm-email/{token} [get]
func (h *AuthHandler) ConfirmEmail(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "ConfirmEmail endpoint not implemented",
	})
}

// RecoverPassword handles password recovery
// @Summary Password recovery
// @Description Send password recovery email
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body models.RecoverPasswordRequest true "Recovery email data"
// @Success 200 {object} responses.SuccessResponse
// @Failure 400 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Router /auth/recover-password [post]
func (h *AuthHandler) RecoverPassword(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "RecoverPassword endpoint not implemented",
	})
}

// ResetPassword handles password reset
// @Summary Reset password
// @Description Reset user password with recovery token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body models.ResetPasswordRequest true "Password reset data"
// @Success 200 {object} responses.SuccessResponse
// @Failure 400 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Router /auth/reset-password [post]
func (h *AuthHandler) ResetPassword(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "ResetPassword endpoint not implemented",
	})
}
