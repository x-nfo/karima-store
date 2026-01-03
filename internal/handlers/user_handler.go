package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/karima-store/internal/models"
	"github.com/karima-store/internal/services"
	"github.com/karima-store/internal/utils"
)

type UserHandler struct {
	userService services.UserService
}

func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// GetUsers godoc
// @Summary Get all users
// @Description Get list of all users (admin only)
// @Tags users
// @Accept json
// @Produce json
// @Param limit query int false "Limit" default(20)
// @Param offset query int false "Offset" default(0)
// @Security KratosSession
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /users [get]
func (h *UserHandler) GetUsers(c *fiber.Ctx) error {
	// Parse query params
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	// Get users
	users, total, err := h.userService.GetUsers(limit, offset, nil)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to get users", err)
	}

	return utils.SendSuccess(c, fiber.Map{
		"users":  users,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	}, "Users retrieved successfully")
}

// GetUser godoc
// @Summary Get user by ID
// @Description Get user details by ID (admin or self)
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Security KratosSession
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /users/{id} [get]
func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	// Parse user ID
	userID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Invalid user ID", err)
	}

	// Get user
	user, err := h.userService.GetUserByID(uint(userID))
	if err != nil {
		return utils.SendError(c, fiber.StatusNotFound, "User not found", err)
	}

	return utils.SendSuccess(c, fiber.Map{
		"user": user,
	}, "User retrieved successfully")
}

// GetCurrentUser godoc
// @Summary Get current authenticated user
// @Description Get current user's profile
// @Tags users
// @Accept json
// @Produce json
// @Security KratosSession
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /users/me [get]
func (h *UserHandler) GetCurrentUser(c *fiber.Ctx) error {
	// Get user from context (set by auth middleware)
	user := c.Locals("user")
	if user == nil {
		return utils.SendError(c, fiber.StatusUnauthorized, "User not authenticated", nil)
	}

	return utils.SendSuccess(c, fiber.Map{
		"user": user,
	}, "Current user retrieved successfully")
}

// UpdateUserRole godoc
// @Summary Update user role
// @Description Update a user's role (admin only)
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param body body map[string]string true "Role update request"
// @Security KratosSession
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /users/{id}/role [put]
func (h *UserHandler) UpdateUserRole(c *fiber.Ctx) error {
	// Parse user ID
	userID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Invalid user ID", err)
	}

	// Parse request body
	var req struct {
		Role string `json:"role" validate:"required,oneof=admin customer"`
	}
	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	// Update role
	role := models.UserRole(req.Role)
	if err := h.userService.UpdateUserRole(uint(userID), role); err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to update user role", err)
	}

	return utils.SendSuccess(c, nil, "User role updated successfully")
}

// DeactivateUser godoc
// @Summary Deactivate user
// @Description Deactivate a user account (admin only)
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Security KratosSession
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /users/{id}/deactivate [put]
func (h *UserHandler) DeactivateUser(c *fiber.Ctx) error {
	// Parse user ID
	userID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Invalid user ID", err)
	}

	// Deactivate user
	if err := h.userService.DeactivateUser(uint(userID)); err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to deactivate user", err)
	}

	return utils.SendSuccess(c, nil, "User deactivated successfully")
}

// ActivateUser godoc
// @Summary Activate user
// @Description Activate a user account (admin only)
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Security KratosSession
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /users/{id}/activate [put]
func (h *UserHandler) ActivateUser(c *fiber.Ctx) error {
	// Parse user ID
	userID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Invalid user ID", err)
	}

	// Activate user
	if err := h.userService.ActivateUser(uint(userID)); err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to activate user", err)
	}

	return utils.SendSuccess(c, nil, "User activated successfully")
}

// GetUserStats godoc
// @Summary Get user statistics
// @Description Get user statistics (admin only)
// @Tags users
// @Accept json
// @Produce json
// @Security KratosSession
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /users/stats [get]
func (h *UserHandler) GetUserStats(c *fiber.Ctx) error {
	stats, err := h.userService.GetUserStats()
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to get user stats", err)
	}

	return utils.SendSuccess(c, stats, "User statistics retrieved successfully")
}
