package services

import (
	"fmt"
	"time"

	"github.com/karima-store/internal/models"
	"github.com/karima-store/internal/repository"
)

type AuthService interface {
	SyncUser(kratosIdentity *models.KratosIdentity, email string) (*models.User, error)
	GetUserByID(id uint) (*models.User, error)
}

type authService struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{userRepo: userRepo}
}

// SyncUser ensures that the Kratos user exists in the local database
func (s *authService) SyncUser(kratosIdentity *models.KratosIdentity, email string) (*models.User, error) {
	// 1. Try to find by Kratos ID
	user, err := s.userRepo.FindByKratosID(kratosIdentity.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to check user by kratos ID: %w", err)
	}

	// 2. If not found by Kratos ID, try by Email (migration/linking case)
	if user == nil {
		user, err = s.userRepo.FindByEmail(email)
		if err != nil {
			return nil, fmt.Errorf("failed to check user by email: %w", err)
		}
	}

	// 3. User still not found? Create new user
	if user == nil {
		newUser := &models.User{
			KratosID:   kratosIdentity.ID,
			Email:      email,
			Role:       models.RoleCustomer, // Default role
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
			IsActive:   true,
			IsVerified: true,   // Assumed verified if Kratos lets them login (simplification)
			FullName:   "User", // Placeholder, ideally get from traits
		}

		// Attempt to extract name from traits if available
		if nameMap, ok := kratosIdentity.Traits["name"].(map[string]interface{}); ok {
			first, _ := nameMap["first"].(string)
			last, _ := nameMap["last"].(string)
			newUser.FullName = fmt.Sprintf("%s %s", first, last)
		} else if nameStr, ok := kratosIdentity.Traits["name"].(string); ok {
			newUser.FullName = nameStr
		}

		if err := s.userRepo.Create(newUser); err != nil {
			return nil, fmt.Errorf("failed to create user: %w", err)
		}
		return newUser, nil
	}

	// 4. Update existing user with Kratos ID if missing (Link Account)
	if user.KratosID == "" {
		user.KratosID = kratosIdentity.ID
		if err := s.userRepo.Update(user); err != nil {
			return nil, fmt.Errorf("failed to link kratos ID to user: %w", err)
		}
	}

	return user, nil
}

func (s *authService) GetUserByID(id uint) (*models.User, error) {
	return s.userRepo.FindByID(id)
}
