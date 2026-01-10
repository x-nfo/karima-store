package services

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/karima-store/internal/config"
	"github.com/karima-store/internal/models"
	"github.com/karima-store/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(email, password, fullName string) (*models.User, string, error)
	Login(email, password string) (*models.User, string, error)
	FindOrCreateByOAuth(provider, email, providerID, fullName, avatar string) (*models.User, string, error)
	GetUserByID(id uint) (*models.User, error)
}

type authService struct {
	userRepo repository.UserRepository
	cfg      *config.Config
}

func NewAuthService(userRepo repository.UserRepository, cfg *config.Config) AuthService {
	return &authService{
		userRepo: userRepo,
		cfg:      cfg,
	}
}

func (s *authService) Register(email, password, fullName string) (*models.User, string, error) {
	// Check if user exists
	existingUser, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, "", err
	}
	if existingUser != nil {
		return nil, "", errors.New("email already registered")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", err
	}

	// Create user
	user := &models.User{
		Email:      email,
		Password:   string(hashedPassword),
		FullName:   fullName,
		Role:       models.RoleCustomer,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		IsActive:   true,
		IsVerified: false, // Email verification logic needed separate
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, "", err
	}

	// Generate Token
	token, err := s.generateToken(user)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (s *authService) Login(email, password string) (*models.User, string, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, "", err
	}
	if user == nil {
		return nil, "", errors.New("invalid credentials")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	if !user.IsActive {
		return nil, "", errors.New("account is inactive")
	}

	// Generate Token
	token, err := s.generateToken(user)
	if err != nil {
		return nil, "", err
	}

	// Update Last Login
	now := time.Now()
	user.LastLoginAt = &now
	_ = s.userRepo.Update(user)

	return user, token, nil
}

func (s *authService) FindOrCreateByOAuth(provider, email, providerID, fullName, avatar string) (*models.User, string, error) {
	// 1. Try to find by Email
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, "", err
	}

	if user == nil {
		// Create new user
		user = &models.User{
			Email:      email,
			FullName:   fullName,
			Avatar:     avatar,
			Role:       models.RoleCustomer,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
			IsActive:   true,
			IsVerified: true, // OAuth is usually verified
		}
		// Note: Password is empty for OAuth users
		if err := s.userRepo.Create(user); err != nil {
			return nil, "", err
		}
	} else {
		// Update avatar if missing or changed
		if user.Avatar == "" || user.Avatar != avatar {
			user.Avatar = avatar
			_ = s.userRepo.Update(user)
		}
	}

	token, err := s.generateToken(user)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (s *authService) GetUserByID(id uint) (*models.User, error) {
	return s.userRepo.FindByID(id)
}

func (s *authService) generateToken(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // 24 hours
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.JWTSecret))
}
