package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/karima-store/internal/config"
	"github.com/karima-store/internal/middleware"
	"github.com/karima-store/internal/services"
)

type AuthHandler struct {
	authService services.AuthService
	config      *config.Config
}

func NewAuthHandler(authService services.AuthService, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		config:      cfg,
	}
}

// Register initiates the registration flow
// @Summary Initiate Registration
// @Description Redirects to Kratos registration UI or returns init flow URL
// @Tags auth
// @Produce json
// @Success 303 {string} string "Redirect to Kratos"
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	// For API clients, we might want to return the flow URL instead of redirecting
	// For now, we'll redirect to the browser-based UI flow
	return c.Redirect(h.config.KratosPublicURL+"/self-service/registration/browser", fiber.StatusSeeOther)
}

// Login initiates the login flow
// @Summary Initiate Login
// @Description Redirects to Kratos login UI or returns init flow URL
// @Tags auth
// @Produce json
// @Success 303 {string} string "Redirect to Kratos"
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	return c.Redirect(h.config.KratosPublicURL+"/self-service/login/browser", fiber.StatusSeeOther)
}

// Logout initiates the logout flow
// @Summary Initiate Logout
// @Description Redirects to Kratos logout UI
// @Tags auth
// @Produce json
// @Success 303 {string} string "Redirect to Kratos"
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	return c.Redirect(h.config.KratosPublicURL+"/self-service/browser/flows/logout", fiber.StatusSeeOther)
}

// Me returns the current authenticated user's details
// @Summary Get Current User
// @Description detailed user info merging Kratos identity and local DB user
// @Tags auth
// @Produce json
// @Success 200 {object} models.User
// @Router /auth/me [get]
func (h *AuthHandler) Me(c *fiber.Ctx) error {
	// This endpoint is protected by KratosMiddleware, so we should have locals
	kratosSession, ok := c.Locals("session").(*middleware.KratosSession)
	if !ok || kratosSession == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Not authenticated",
		})
	}

	email, _ := c.Locals("user_email").(string)

	// Sync User (Lazy Sync)
	user, err := h.authService.SyncUser(&kratosSession.Identity, email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to sync user data: " + err.Error(),
		})
	}

	// Update LastLoginAt if needed (logic can be added to service)

	return c.JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"user":    user,
			"session": kratosSession,
		},
	})
}
