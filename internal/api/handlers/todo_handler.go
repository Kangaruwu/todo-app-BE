package handlers

import (
	"strconv"

	"go-backend-todo/internal/api/middlewares"
	"go-backend-todo/internal/api/responses"
	"go-backend-todo/internal/models"
	"go-backend-todo/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// TodoHandler struct contain dependencies
type TodoHandler struct {
	todoService service.TodoService
}

// NewTodoHandler create a new instance of todo handler
func NewTodoHandler(todoService service.TodoService) *TodoHandler {
	return &TodoHandler{
		todoService: todoService,
	}
}

// GetTodos lấy danh sách todos với pagination và filter
// @Summary Get user's todos with pagination and filters
// @Description Retrieve a paginated list of todos for the authenticated user with optional filters
// @Tags Todos
// @Accept json
// @Produce json
// @Param limit query int false "Number of items per page (default: 10)" minimum(1) maximum(100)
// @Param offset query int false "Number of items to skip (default: 0)" minimum(0)
// @Param completed query bool false "Filter by completion status"
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "Paginated list of todos"
// @Failure 400 {object} map[string]string "Invalid query parameters"
// @Failure 401 {object} map[string]string "Unauthorized - missing or invalid token"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /todos [get]
func (h *TodoHandler) GetTodos(c *fiber.Ctx) error {
	// Lấy userID từ JWT token
	userID, err := middlewares.GetUserIDFromContext(c)
	if err != nil {
		return responses.Unauthorized(c, "User not authenticated")
	}

	// Parse query parameters
	limitStr := c.Query("limit", "10")
	offsetStr := c.Query("offset", "0")
	completedStr := c.Query("completed")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		return responses.BadRequest(c, "Invalid limit parameter (1-100)")
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		return responses.BadRequest(c, "Invalid offset parameter")
	}

	var completed *bool
	if completedStr != "" {
		completedBool, err := strconv.ParseBool(completedStr)
		if err != nil {
			return responses.BadRequest(c, "Invalid completed parameter")
		}
		completed = &completedBool
	}

	todos, total, err := h.todoService.GetTodosWithPagination(c.Context(), userID, limit, offset, completed)
	if err != nil {
		return responses.InternalServerErrorWithError(c, "Failed to get todos", err)
	}

	page := (offset / limit) + 1
	return responses.OKWithPagination(c, "Todos retrieved successfully", todos, page, limit, total)
}

// CreateTodo tạo todo mới
// @Summary Create a new todo
// @Description Create a new todo item for the authenticated user
// @Tags Todos
// @Accept json
// @Produce json
// @Param todo body models.CreateTodoRequest true "Todo creation data"
// @Security BearerAuth
// @Success 201 {object} map[string]interface{} "Todo created successfully"
// @Failure 400 {object} map[string]string "Invalid request data"
// @Failure 401 {object} map[string]string "Unauthorized - missing or invalid token"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /todos [post]
func (h *TodoHandler) CreateTodo(c *fiber.Ctx) error {
	// Lấy userID từ JWT token
	userID, err := middlewares.GetUserIDFromContext(c)
	if err != nil {
		return responses.Unauthorized(c, "User not authenticated")
	}

	var req models.CreateTodoRequest
	if err := c.BodyParser(&req); err != nil {
		return responses.BadRequest(c, "Invalid request body format")
	}

	// Validate request
	if err := middlewares.ValidateStruct(&req); err != nil {
		return responses.BadRequestWithError(c, "Validation failed", err)
	}

	todo, err := h.todoService.CreateTodo(c.Context(), req, userID)
	if err != nil {
		return responses.InternalServerErrorWithError(c, "Failed to create todo", err)
	}

	return responses.Created(c, "Todo created successfully", todo)
}

// GetTodo lấy todo theo ID
// @Summary Get todo by ID
// @Description Retrieve a specific todo item by its unique identifier
// @Tags Todos
// @Accept json
// @Produce json
// @Param id path string true "Todo ID" format(uuid)
// @Security BearerAuth
// @Router /todos/{id} [get]
func (h *TodoHandler) GetTodo(c *fiber.Ctx) error {
	// Lấy userID từ JWT token
	userID, err := middlewares.GetUserIDFromContext(c)
	if err != nil {
		return responses.Unauthorized(c, "User not authenticated")
	}

	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return responses.BadRequest(c, "Invalid todo ID format")
	}

	todo, err := h.todoService.GetTodoByID(c.Context(), id, userID)
	if err != nil {
		return responses.NotFound(c, "Todo not found")
	}

	return responses.OK(c, "Todo retrieved successfully", todo)
}

// UpdateTodo cập nhật todo
// @Summary Update todo
// @Description Update an existing todo item's information
// @Tags Todos
// @Accept json
// @Produce json
// @Param id path string true "Todo ID" format(uuid)
// @Param todo body models.UpdateTodoRequest true "Todo update data"
// @Security BearerAuth
// @Router /todos/{id} [put]
func (h *TodoHandler) UpdateTodo(c *fiber.Ctx) error {
	// Lấy userID từ JWT token
	userID, err := middlewares.GetUserIDFromContext(c)
	if err != nil {
		return responses.Unauthorized(c, "User not authenticated")
	}

	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return responses.BadRequest(c, "Invalid todo ID format")
	}

	var req models.UpdateTodoRequest
	if err := c.BodyParser(&req); err != nil {
		return responses.BadRequest(c, "Invalid request body format")
	}

	// Validate request
	if err := middlewares.ValidateStruct(&req); err != nil {
		return responses.BadRequestWithError(c, "Validation failed", err)
	}

	todo, err := h.todoService.UpdateTodo(c.Context(), id, req, userID)
	if err != nil {
		return responses.InternalServerErrorWithError(c, "Failed to update todo", err)
	}

	return responses.OK(c, "Todo updated successfully", todo)
}

// DeleteTodo xóa todo
// @Summary Delete todo
// @Description Permanently delete a todo item
// @Tags Todos
// @Accept json
// @Produce json
// @Param id path string true "Todo ID" format(uuid)
// @Security BearerAuth
// @Router /todos/{id} [delete]
func (h *TodoHandler) DeleteTodo(c *fiber.Ctx) error {
	// Lấy userID từ JWT token
	userID, err := middlewares.GetUserIDFromContext(c)
	if err != nil {
		return responses.Unauthorized(c, "User not authenticated")
	}

	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return responses.BadRequest(c, "Invalid todo ID format")
	}

	err = h.todoService.DeleteTodo(c.Context(), id, userID)
	if err != nil {
		return responses.InternalServerErrorWithError(c, "Failed to delete todo", err)
	}

	return responses.OK(c, "Todo deleted successfully", nil)
}

// ToggleTodoStatus toggles todo completion status
// @Summary Toggle todo completion status
// @Description Toggle the completion status of a todo item (completed/incomplete)
// @Tags Todos
// @Accept json
// @Produce json
// @Param id path string true "Todo ID" format(uuid)
// @Security BearerAuth
// @Router /todos/{id}/toggle [patch]
func (h *TodoHandler) ToggleTodoStatus(c *fiber.Ctx) error {
	// TODO: Lấy userID từ JWT token hoặc session
	// userID := uuid.New() // Thay bằng userID thực từ auth

	// idStr := c.Params("id")
	// id, err := uuid.Parse(idStr)
	// if err != nil {
	// 	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
	// 		"error": "Invalid todo ID",
	// 	})
	// }

	// todo, err := h.todoService.ToggleTodoStatus(c.Context(), id, userID)
	// if err != nil {
	// 	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
	// 		"error": err.Error(),
	// 	})
	// }

	// return c.JSON(fiber.Map{
	// 	"message": "Todo status toggled successfully",
	// 	"data":    todo,
	// })
	return c.JSON(fiber.Map{
		"message": "Toggle todo status is not implemented yet",
		"data":    nil,
	})
}

// GetTodosByStatus gets todos filtered by completion status
// @Summary Get todos by completion status
// @Description Retrieve todos filtered by their completion status (completed or incomplete)
// @Tags Todos
// @Accept json
// @Produce json
// @Param status path string true "Todo status" Enums(completed, incomplete)
// @Param limit query int false "Number of items per page (default: 10)" minimum(1) maximum(100)
// @Param offset query int false "Number of items to skip (default: 0)" minimum(0)
// @Security BearerAuth
// @Router /todos/status/{status} [get]
func (h *TodoHandler) GetTodosByStatus(c *fiber.Ctx) error {
	// TODO: Lấy userID từ JWT token hoặc session
	userID := uuid.New() // Thay bằng userID thực từ auth

	statusStr := c.Params("status")
	var completed bool

	switch statusStr {
	case "completed":
		completed = true
	case "incomplete":
		completed = false
	default:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid status. Use 'completed' or 'incomplete'",
		})
	}

	// Parse query parameters
	limitStr := c.Query("limit", "10")
	offsetStr := c.Query("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid limit parameter",
		})
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid offset parameter",
		})
	}

	todos, total, err := h.todoService.GetTodosWithPagination(c.Context(), userID, limit, offset, &completed)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data":   todos,
		"total":  total,
		"limit":  limit,
		"offset": offset,
		"status": statusStr,
	})
}

// GetTodoStats gets user's todo statistics
// @Summary Get todo statistics
// @Description Retrieve statistics about user's todos (total, completed, incomplete)
// @Tags Todos
// @Accept json
// @Produce json
// @Security BearerAuth
// @Router /todos/stats [get]
func (h *TodoHandler) GetTodoStats(c *fiber.Ctx) error {
	// TODO: Lấy userID từ JWT token hoặc session
	// userID := uuid.New() // Thay bằng userID thực từ auth

	// stats, err := h.todoService.GetTodoStats(c.Context(), userID)
	// if err != nil {
	// 	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
	// 		"error": err.Error(),
	// 	})
	// }

	return c.JSON(fiber.Map{
		"message": "Todo statistics are not implemented yet",
	})
}

// Backward compatibility - keep old function signatures for existing routes
func GetTodos(c *fiber.Ctx) error {
	// TODO: Implement dependency injection để get handler instance
	return c.JSON(fiber.Map{
		"message": "Please use the new handler with dependency injection",
		"data":    []interface{}{},
	})
}

func CreateTodo(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "Please use the new handler with dependency injection",
	})
}

func GetTodo(c *fiber.Ctx) error {
	id := c.Params("id")
	return c.JSON(fiber.Map{
		"message": "Please use the new handler with dependency injection",
		"id":      id,
	})
}

func UpdateTodo(c *fiber.Ctx) error {
	id := c.Params("id")
	return c.JSON(fiber.Map{
		"message": "Please use the new handler with dependency injection",
		"id":      id,
	})
}

func DeleteTodo(c *fiber.Ctx) error {
	id := c.Params("id")
	return c.JSON(fiber.Map{
		"message": "Please use the new handler with dependency injection",
		"id":      id,
	})
}
