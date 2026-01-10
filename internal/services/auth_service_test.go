package services

import (
	"testing"

	"github.com/karima-store/internal/config"
	"github.com/karima-store/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
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

func setupAuthService() (*authService, *MockUserRepository) {
	mockRepo := new(MockUserRepository)
	cfg := &config.Config{
		JWTSecret: "test-secret",
	}
	service := NewAuthService(mockRepo, cfg).(*authService)
	return service, mockRepo
}

func TestAuthService_Register(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		service, mockRepo := setupAuthService()
		email := "test@example.com"
		password := "password123"
		fullName := "Test User"

		mockRepo.On("FindByEmail", email).Return(nil, nil)
		mockRepo.On("Create", mock.AnythingOfType("*models.User")).Return(nil)

		user, token, err := service.Register(email, password, fullName)

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, email, user.Email)
		assert.Equal(t, fullName, user.FullName)
		assert.NotEmpty(t, token)
		mockRepo.AssertExpectations(t)
	})

	t.Run("EmailAlreadyExists", func(t *testing.T) {
		service, mockRepo := setupAuthService()
		email := "test@example.com"
		existingUser := &models.User{Email: email}

		mockRepo.On("FindByEmail", email).Return(existingUser, nil)

		user, token, err := service.Register(email, "password", "name")

		if assert.Error(t, err) {
			assert.Equal(t, "email already registered", err.Error())
		}
		assert.Nil(t, user)
		assert.Empty(t, token)
		mockRepo.AssertExpectations(t)
	})
}

func TestAuthService_Login(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		service, mockRepo := setupAuthService()
		email := "test@example.com"
		password := "password123"
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

		existingUser := &models.User{
			Email:    email,
			Password: string(hashedPassword),
			IsActive: true,
		}

		mockRepo.On("FindByEmail", email).Return(existingUser, nil)
		mockRepo.On("Update", mock.AnythingOfType("*models.User")).Return(nil)

		user, token, err := service.Login(email, password)

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.NotEmpty(t, token)
		mockRepo.AssertExpectations(t)
	})

	t.Run("InvalidCredentials", func(t *testing.T) {
		service, mockRepo := setupAuthService()
		email := "test@example.com"
		mockRepo.On("FindByEmail", email).Return(nil, nil)

		user, token, err := service.Login(email, "wrongpassword")

		if assert.Error(t, err) {
			assert.Equal(t, "invalid credentials", err.Error())
		}
		assert.Nil(t, user)
		assert.Empty(t, token)
	})
}

func TestAuthService_FindOrCreateByOAuth(t *testing.T) {
	t.Run("CreateNewUser", func(t *testing.T) {
		service, mockRepo := setupAuthService()
		email := "oauth@example.com"
		provider := "google"
		providerID := "12345"

		mockRepo.On("FindByEmail", email).Return(nil, nil)
		mockRepo.On("Create", mock.AnythingOfType("*models.User")).Return(nil)

		user, token, err := service.FindOrCreateByOAuth(provider, email, providerID, "OAuth User", "avatar.jpg")

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, email, user.Email)
		assert.NotEmpty(t, token)
		mockRepo.AssertExpectations(t)
	})

	t.Run("ExistingUser", func(t *testing.T) {
		service, mockRepo := setupAuthService()
		email := "oauth@example.com"
		existingUser := &models.User{
			Email:  email,
			Avatar: "old.jpg",
		}

		mockRepo.On("FindByEmail", email).Return(existingUser, nil)
		mockRepo.On("Update", mock.AnythingOfType("*models.User")).Return(nil) // Updates avatar

		user, token, err := service.FindOrCreateByOAuth("google", email, "123", "User", "new.jpg")

		assert.NoError(t, err)
		assert.Equal(t, existingUser, user)
		assert.NotEmpty(t, token)
		mockRepo.AssertExpectations(t)
	})
}
