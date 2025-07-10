package handlers

import (
	"go-backend-todo/internal/models"
	"go-backend-todo/internal/service"

	"encoding/json"

	"github.com/gofiber/fiber/v2"
)

// AuthHandler handles authentication HTTP requests
type AuthHandler struct {
	authService service.AuthService
}

// NewAuthHandler creates a new instance of auth handler
func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Login handles user login
// @Summary User login
// @Description Authenticate user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "Login credentials"
// @Success 200 {object} models.LoginResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	body := c.Body()

	if len(body) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "Request body cannot be empty",
			"data":    nil,
		})
	}

	var req models.LoginRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "Invalid request body",
			"data":    nil,
		})
	}

	loginResponse, err := h.authService.Login(c.Context(), &req)
	if err != nil {
		if err.Error() == "invalid credentials" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  fiber.StatusBadRequest,
				"message": "Invalid email or password",
				"data":    nil,
			})
		}
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  fiber.StatusUnauthorized,
			"message":  "Failed to login user",
			"data":    nil,
		})
	}

	return c.JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"message": "Login successful",
		"data":    loginResponse,
	})
}

// Register handles user registration
// @Summary User registration
// @Description Register a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.RegisterRequest true "User registration data"
// @Success 201 {object} models.RegisterResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 409 {object} models.ErrorResponse
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	body := c.Body()

	if len(body) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "Request body cannot be empty",
			"data":    nil,
		})
	}

	var req models.RegisterRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "Invalid request body",
			"data":    nil,
		})
	}

	if err := h.authService.Register(c.Context(), &req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.StatusInternalServerError,
			"message": "Failed to register user",
			"data":    nil,
		})
	}

	return c.JSON(fiber.Map{
		"status":  fiber.StatusCreated,
		"message": "User registered successfully",
		"data":    nil,
	})
}

// ConfirmEmail handles email confirmation
// @Summary Confirm email address
// @Description Confirm user email with verification token
// @Tags auth
// @Accept json
// @Produce json
// @Param token path string true "Email confirmation token"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /auth/confirm-email/{token} [get]
func (h *AuthHandler) ConfirmEmail(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "ConfirmEmail endpoint not implemented",
	})
}

// RecoverPassword handles password recovery
// @Summary Password recovery
// @Description Send password recovery email
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.RecoverPasswordRequest true "Recovery email data"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /auth/recover-password [post]
func (h *AuthHandler) RecoverPassword(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "RecoverPassword endpoint not implemented",
	})
}

// ResetPassword handles password reset
// @Summary Reset password
// @Description Reset user password with recovery token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.ResetPasswordRequest true "Password reset data"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /auth/reset-password [post]
func (h *AuthHandler) ResetPassword(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "ResetPassword endpoint not implemented",
	})
}
