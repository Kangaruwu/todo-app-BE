package handlers

import (
	"strconv"

	"go-backend-todo/internal/models"
	"go-backend-todo/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// TodoHandler struct chứa các dependencies
type TodoHandler struct {
	todoService service.TodoService
}

// NewTodoHandler tạo một instance mới của todo handler
func NewTodoHandler(todoService service.TodoService) *TodoHandler {
	return &TodoHandler{
		todoService: todoService,
	}
}

// GetTodos lấy danh sách todos với pagination và filter
func (h *TodoHandler) GetTodos(c *fiber.Ctx) error {
	// TODO: Lấy userID từ JWT token hoặc session
	// Tạm thời hardcode để test
	temp := "eef14407-d602-40ad-9c79-2241c3d06deb"
	userID := uuid.Must(uuid.Parse(temp)) 
	//userID := uuid.New() // Thay bằng userID thực từ auth

	// Parse query parameters
	limitStr := c.Query("limit", "10")
	offsetStr := c.Query("offset", "0")
	completedStr := c.Query("completed")

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

	var completed *bool
	if completedStr != "" {
		completedBool, err := strconv.ParseBool(completedStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid completed parameter",
			})
		}
		completed = &completedBool
	}

	todos, total, err := h.todoService.GetTodosWithPagination(c.Context(), userID, limit, offset, completed)
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
	})
}

// CreateTodo tạo todo mới
func (h *TodoHandler) CreateTodo(c *fiber.Ctx) error {
	// TODO: Lấy userID từ JWT token hoặc session
	userID := uuid.New() // Thay bằng userID thực từ auth

	var req models.CreateTodoRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	todo, err := h.todoService.CreateTodo(c.Context(), req, userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Todo created successfully",
		"data":    todo,
	})
}

// GetTodo lấy todo theo ID
func (h *TodoHandler) GetTodo(c *fiber.Ctx) error {
	// TODO: Lấy userID từ JWT token hoặc session
	userID := uuid.New() // Thay bằng userID thực từ auth

	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid todo ID",
		})
	}

	todo, err := h.todoService.GetTodoByID(c.Context(), id, userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data": todo,
	})
}

// UpdateTodo cập nhật todo
func (h *TodoHandler) UpdateTodo(c *fiber.Ctx) error {
	// TODO: Lấy userID từ JWT token hoặc session
	userID := uuid.New() // Thay bằng userID thực từ auth

	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid todo ID",
		})
	}

	var req models.UpdateTodoRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	todo, err := h.todoService.UpdateTodo(c.Context(), id, req, userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Todo updated successfully",
		"data":    todo,
	})
}

// DeleteTodo xóa todo
func (h *TodoHandler) DeleteTodo(c *fiber.Ctx) error {
	// TODO: Lấy userID từ JWT token hoặc session
	userID := uuid.New() // Thay bằng userID thực từ auth

	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid todo ID",
		})
	}

	err = h.todoService.DeleteTodo(c.Context(), id, userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Todo deleted successfully",
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
