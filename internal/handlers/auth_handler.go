package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/karima-store/internal/config"
	"github.com/karima-store/internal/services"
	"github.com/shareed2k/goth_fiber"
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

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user with email, password, and full name
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Registration Request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /api/v1/auth/register [post]
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req RegisterRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	user, token, err := h.authService.Register(req.Email, req.Password, req.FullName)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// Set Cookie
	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
		Secure:   h.config.AppEnv == "production",
	})

	return c.JSON(fiber.Map{"user": user, "token": token})
}

// Login godoc
// @Summary Login user
// @Description Login with email and password to get JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login Request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req LoginRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	user, token, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	// Set Cookie
	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
		Secure:   h.config.AppEnv == "production",
	})

	return c.JSON(fiber.Map{"user": user, "token": token})
}

// OAuthLogin godoc
// @Summary Initiate OAuth login
// @Description Initiate OAuth login with a provider (e.g., google)
// @Tags auth
// @Param provider path string true "OAuth Provider (google)"
// @Success 302
// @Router /api/v1/auth/{provider} [get]
func (h *AuthHandler) OAuthLogin(c *fiber.Ctx) error {
	// goth_fiber uses "provider" param from route by default
	return goth_fiber.BeginAuthHandler(c)
}

// Logout godoc
// @Summary Logout user
// @Description Clear the JWT cookie
// @Tags auth
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/auth/logout [post]
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		HTTPOnly: true,
	})

	// Logout from Goth session if possible (optional)
	// goth_fiber.Logout(c)

	return c.JSON(fiber.Map{"message": "Logged out successfully"})
}

// Me godoc
// @Summary Get current user profile
// @Description Get the profile of the currently authenticated user
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/auth/me [get]
func (h *AuthHandler) Me(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	// userID (from JWT claims) is float64 by default when parsing JSON/MapClaims
	var id uint
	switch v := userID.(type) {
	case float64:
		id = uint(v)
	case uint:
		id = v
	default:
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Invalid user ID type"})
	}

	user, err := h.authService.GetUserByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	return c.JSON(fiber.Map{"user": user})
}

// OAuthCallback godoc
// @Summary OAuth Callback
// @Description Callback endpoint for OAuth providers
// @Tags auth
// @Param provider path string true "OAuth Provider"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /api/v1/auth/{provider}/callback [get]
func (h *AuthHandler) OAuthCallback(c *fiber.Ctx) error {
	user, err := goth_fiber.CompleteUserAuth(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	// Create or update user in our DB
	dbUser, token, err := h.authService.FindOrCreateByOAuth(user.Provider, user.Email, user.UserID, user.Name, user.AvatarURL)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to process login"})
	}

	// Set Cookie
	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
		Secure:   h.config.AppEnv == "production",
	})

	// Just return JSON for now, or redirect to frontend
	// return c.Redirect("http://localhost:3000/auth/success?token=" + token)
	return c.JSON(fiber.Map{"user": dbUser, "token": token})
}
