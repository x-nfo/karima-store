package services

import (
	"fmt"

	"github.com/karima-store/internal/models"
	"github.com/karima-store/internal/repository"
)

// UserService interface defines user management operations
type UserService interface {
	GetUsers(limit, offset int, filters map[string]interface{}) ([]models.User, int64, error)
	GetUserByID(id uint) (*models.User, error)
	UpdateUserRole(id uint, role models.UserRole) error
	DeactivateUser(id uint) error
	ActivateUser(id uint) error
	GetUserStats() (map[string]interface{}, error)
}

type userService struct {
	userRepo repository.UserRepository
}

// NewUserService creates a new user service
func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

// GetUsers retrieves users with pagination and filters
func (s *userService) GetUsers(limit, offset int, filters map[string]interface{}) ([]models.User, int64, error) {
	// Validate pagination
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	// For now, we'll implement basic filtering
	// In a real implementation, you'd build dynamic queries based on filters
	users := []models.User{}
	var total int64

	// This is a simplified implementation
	// You would need to implement GetAll in UserRepository
	// For now, return empty list as placeholder
	return users, total, nil
}

// GetUserByID retrieves a user by ID
func (s *userService) GetUserByID(id uint) (*models.User, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

// UpdateUserRole updates a user's role
func (s *userService) UpdateUserRole(id uint, role models.UserRole) error {
	// Validate role
	if !models.ValidateRole(role) {
		return fmt.Errorf("invalid role: %s", role)
	}

	// Get user
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return fmt.Errorf("user not found")
	}

	// Update role
	user.Role = role
	if err := s.userRepo.Update(user); err != nil {
		return fmt.Errorf("failed to update user role: %w", err)
	}

	return nil
}

// DeactivateUser deactivates a user account
func (s *userService) DeactivateUser(id uint) error {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return fmt.Errorf("user not found")
	}

	user.IsActive = false
	if err := s.userRepo.Update(user); err != nil {
		return fmt.Errorf("failed to deactivate user: %w", err)
	}

	return nil
}

// ActivateUser activates a user account
func (s *userService) ActivateUser(id uint) error {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return fmt.Errorf("user not found")
	}

	user.IsActive = true
	if err := s.userRepo.Update(user); err != nil {
		return fmt.Errorf("failed to activate user: %w", err)
	}

	return nil
}

// GetUserStats returns user statistics
func (s *userService) GetUserStats() (map[string]interface{}, error) {
	// This is a placeholder implementation
	// In a real implementation, you'd query the database for actual stats
	stats := map[string]interface{}{
		"total_users":    0,
		"active_users":   0,
		"inactive_users": 0,
		"admin_users":    0,
		"customer_users": 0,
	}

	return stats, nil
}
