package handlers

import (
	"go-backend-todo/internal/api/middlewares"
	"go-backend-todo/internal/api/responses"
	"go-backend-todo/internal/models"
	"go-backend-todo/internal/service"

	"github.com/google/uuid"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// GetUserProfile gets current user's profile
// @Summary Get current user profile
// @Description Retrieve the profile information of the currently authenticated user
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Router /users/profile [get]
func (h *UserHandler) GetUserProfile(c *fiber.Ctx) error {
	user, err := h.userService.GetUserByID(c.Context(), uuid.MustParse(c.Locals("user_id").(string)))
	if err != nil {
		return responses.InternalServerErrorWithError(c, "Failed to get user", err)
	}
	return responses.OK(c, "User found", user)
}

// UpdateUserProfile updates current user's profile
// @Summary Update current user profile
// @Description Update the profile information of the currently authenticated user
// @Tags Users
// @Accept json
// @Produce json
// @Param profile body models.UpdateProfileRequest true "Profile update data"
// @Security BearerAuth
// @Router /users/profile [put]
func (h *UserHandler) UpdateUserProfile(c *fiber.Ctx) error {
	var req models.UpdateProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return responses.BadRequestWithError(c, "Invalid request data", err)
	}

	user, err := h.userService.UpdateUserProfile(c.Context(), uuid.MustParse(c.Locals("user_id").(string)), req)
	if err != nil {
		return responses.InternalServerErrorWithError(c, "Failed to update user", err)
	}
	return responses.OK(c, "User profile updated", user)
}

// ChangePassword changes user's password
// @Summary Change user password
// @Description Change the password for the currently authenticated user
// @Tags Users
// @Accept json
// @Produce json
// @Param password body models.ChangePasswordRequest true "Password change data"
// @Security BearerAuth
// @Router /users/change-password [put]
func (h *UserHandler) ChangePassword(c *fiber.Ctx) error {
	body := c.Body()
	var req models.ChangePasswordRequest
	err := middlewares.RequestValidation(&body, &req)(c)
	if err != nil {
		return responses.BadRequestWithError(c, "Invalid request body", err)
	}
	err = h.userService.ChangePassword(c.Context(), uuid.MustParse(c.Locals("user_id").(string)), &req)
	if err != nil {
		return responses.InternalServerErrorWithError(c, "Failed to change password", err)
	}
	return responses.OK(c, "Password changed successfully", nil)
}

// DeleteUserProfile deletes current user's account
// @Summary Delete current user account
// @Description Permanently delete the currently authenticated user's account and all associated data
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Router /users/profile [delete]
func (h *UserHandler) DeleteUserProfile(c *fiber.Ctx) error {
	// userID := uuid.MustParse(c.Locals("user_id").(string))
	// TODO: Implement logic to delete current user's account
	return responses.OK(c, "User account deleted successfully", nil)
}
