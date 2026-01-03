package services

import (
	"errors"
	"testing"

	"github.com/karima-store/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository for testing
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) FindByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) FindByKratosID(kratosID string) (*models.User, error) {
	args := m.Called(kratosID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) Create(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) Update(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) FindByID(id uint) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func TestAuthService_SyncUser(t *testing.T) {
	// Test Case 1: Existing user by Kratos ID
	t.Run("Existing_User_By_KratosID", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		service := NewAuthService(mockRepo)

		kratosID := "test-kratos-id"
		email := "test@example.com"
		identity := &models.KratosIdentity{
			ID:     kratosID,
			Traits: map[string]interface{}{"email": email},
		}

		existingUser := &models.User{
			KratosID: kratosID,
			Email:    email,
		}

		mockRepo.On("FindByKratosID", kratosID).Return(existingUser, nil)

		user, err := service.SyncUser(identity, email)

		assert.NoError(t, err)
		assert.Equal(t, existingUser, user)
		mockRepo.AssertExpectations(t)
	})

	// Test Case 2: New user (not found by ID or Email)
	t.Run("New_User", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		service := NewAuthService(mockRepo)

		kratosID := "new-kratos-id"
		email := "new@example.com"
		identity := &models.KratosIdentity{
			ID:     kratosID,
			Traits: map[string]interface{}{"email": email, "name": "New User"},
		}

		mockRepo.On("FindByKratosID", kratosID).Return(nil, nil)
		mockRepo.On("FindByEmail", email).Return(nil, nil)
		mockRepo.On("Create", mock.AnythingOfType("*models.User")).Return(nil)

		user, err := service.SyncUser(identity, email)

		assert.NoError(t, err)
		assert.Equal(t, kratosID, user.KratosID)
		assert.Equal(t, email, user.Email)
		assert.Equal(t, "New User", user.FullName)
		mockRepo.AssertExpectations(t)
	})

	// Test Case 3: Link existing legacy user (found by Email, no Kratos ID)
	t.Run("Link_Legacy_User", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		service := NewAuthService(mockRepo)

		kratosID := "link-kratos-id"
		email := "legacy@example.com"
		identity := &models.KratosIdentity{
			ID:     kratosID,
			Traits: map[string]interface{}{"email": email},
		}

		legacyUser := &models.User{
			Email:    email,
			KratosID: "",
			FullName: "Legacy User",
		}

		mockRepo.On("FindByKratosID", kratosID).Return(nil, nil)
		mockRepo.On("FindByEmail", email).Return(legacyUser, nil)
		mockRepo.On("Update", legacyUser).Return(nil)

		user, err := service.SyncUser(identity, email)

		assert.NoError(t, err)
		assert.Equal(t, legacyUser, user)
		assert.Equal(t, kratosID, user.KratosID)
		mockRepo.AssertExpectations(t)
	})

	// Test Case 4: Error checking Kratos ID
	t.Run("Error_Checking_KratosID", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		service := NewAuthService(mockRepo)

		kratosID := "error-id"
		email := "error@example.com"
		identity := &models.KratosIdentity{ID: kratosID}

		mockRepo.On("FindByKratosID", kratosID).Return(nil, errors.New("db error"))

		user, err := service.SyncUser(identity, email)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "failed to check user by kratos ID")
		mockRepo.AssertExpectations(t)
	})
}
