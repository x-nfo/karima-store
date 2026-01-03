package services

import (
	"errors"
	"testing"
	"time"

	"github.com/karima-store/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Helper function to create test user
func createTestUser(id uint, email string, role models.UserRole, isActive bool) *models.User {
	return &models.User{
		ID:        id,
		Email:     email,
		FullName:  "Test User",
		KratosID:  "kratos-123",
		Phone:     "08123456789",
		Role:      role,
		IsActive:  isActive,
		IsVerified: true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// ============================================================================
// SERVICE INITIALIZATION TESTS
// ============================================================================

func TestNewUserService(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	assert.NotNil(t, service, "Service should not be nil")
}

func TestUserService_ImplementsInterface(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	// Verify service implements interface
	var _ UserService = service
}

// ============================================================================
// GET USERS TESTS (FR-061: Profile Management)
// ============================================================================

func TestUserService_GetUsers_DefaultPagination(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	users, total, err := service.GetUsers(0, 0, nil)

	assert.NoError(t, err, "Should not return error")
	assert.NotNil(t, users, "Users should not be nil")
	assert.Equal(t, int64(0), total, "Total should be 0 for placeholder implementation")
	assert.Empty(t, users, "Users should be empty for placeholder implementation")
}

func TestUserService_GetUsers_ValidPagination(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	testCases := []struct {
		name          string
		limit         int
		offset        int
		expectedLimit int
	}{
		{"Limit 10", 10, 0, 10},
		{"Limit 20", 20, 0, 20},
		{"Limit 50", 50, 0, 50},
		{"Limit 100", 100, 0, 100},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			users, total, err := service.GetUsers(tc.limit, tc.offset, nil)

			assert.NoError(t, err, "Should not return error")
			assert.NotNil(t, users, "Users should not be nil")
			assert.Equal(t, int64(0), total)
		})
	}
}

func TestUserService_GetUsers_LimitExceedsMaximum(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	// Limit > 100 should be clamped to 20 (default)
	users, total, err := service.GetUsers(150, 0, nil)

	assert.NoError(t, err, "Should not return error")
	assert.NotNil(t, users, "Users should not be nil")
	assert.Equal(t, int64(0), total)
}

func TestUserService_GetUsers_NegativeLimit(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	// Negative limit should be clamped to 20 (default)
	users, total, err := service.GetUsers(-10, 0, nil)

	assert.NoError(t, err, "Should not return error")
	assert.NotNil(t, users, "Users should not be nil")
	assert.Equal(t, int64(0), total)
}

func TestUserService_GetUsers_NegativeOffset(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	// Negative offset should be clamped to 0
	users, total, err := service.GetUsers(20, -5, nil)

	assert.NoError(t, err, "Should not return error")
	assert.NotNil(t, users, "Users should not be nil")
	assert.Equal(t, int64(0), total)
}

func TestUserService_GetUsers_WithFilters(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	filters := map[string]interface{}{
		"role":     "admin",
		"is_active": true,
	}

	users, total, err := service.GetUsers(20, 0, filters)

	assert.NoError(t, err, "Should not return error")
	assert.NotNil(t, users, "Users should not be nil")
	assert.Equal(t, int64(0), total)
}

// ============================================================================
// GET USER BY ID TESTS (FR-061: Profile Management)
// ============================================================================

func TestUserService_GetUserByID_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	expectedUser := createTestUser(1, "test@example.com", models.RoleCustomer, true)
	mockRepo.On("FindByID", uint(1)).Return(expectedUser, nil).Once()

	user, err := service.GetUserByID(1)

	require.NoError(t, err, "Should not return error")
	assert.NotNil(t, user, "User should not be nil")
	assert.Equal(t, uint(1), user.ID, "User ID should match")
	assert.Equal(t, "test@example.com", user.Email, "User email should match")
	assert.Equal(t, models.RoleCustomer, user.Role, "User role should match")
	assert.True(t, user.IsActive, "User should be active")
	mockRepo.AssertExpectations(t)
}

func TestUserService_GetUserByID_NotFound(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	mockRepo.On("FindByID", uint(999)).Return(nil, nil).Once()

	user, err := service.GetUserByID(999)

	assert.Error(t, err, "Should return error for non-existent user")
	assert.Nil(t, user, "User should be nil")
	assert.Contains(t, err.Error(), "user not found", "Error message should contain 'user not found'")
	mockRepo.AssertExpectations(t)
}

func TestUserService_GetUserByID_RepositoryError(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	expectedErr := errors.New("database connection failed")
	mockRepo.On("FindByID", uint(1)).Return(nil, expectedErr).Once()

	user, err := service.GetUserByID(1)

	assert.Error(t, err, "Should return error")
	assert.Nil(t, user, "User should be nil")
	assert.Contains(t, err.Error(), "failed to get user", "Error should contain 'failed to get user'")
	assert.Contains(t, err.Error(), "database connection failed", "Error should contain repository error")
	mockRepo.AssertExpectations(t)
}

func TestUserService_GetUserByID_ZeroID(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	mockRepo.On("FindByID", uint(0)).Return(nil, nil).Once()

	user, err := service.GetUserByID(0)

	assert.Error(t, err, "Should return error for zero ID")
	assert.Nil(t, user, "User should be nil")
	mockRepo.AssertExpectations(t)
}

func TestUserService_GetUserByID_AdminUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	adminUser := createTestUser(1, "admin@example.com", models.RoleAdmin, true)
	mockRepo.On("FindByID", uint(1)).Return(adminUser, nil).Once()

	user, err := service.GetUserByID(1)

	require.NoError(t, err)
	assert.Equal(t, models.RoleAdmin, user.Role)
	assert.True(t, user.IsActive)
	mockRepo.AssertExpectations(t)
}

func TestUserService_GetUserByID_InactiveUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	inactiveUser := createTestUser(1, "inactive@example.com", models.RoleCustomer, false)
	mockRepo.On("FindByID", uint(1)).Return(inactiveUser, nil).Once()

	user, err := service.GetUserByID(1)

	require.NoError(t, err)
	assert.False(t, user.IsActive, "User should be inactive")
	mockRepo.AssertExpectations(t)
}

// ============================================================================
// UPDATE USER ROLE TESTS (FR-062: Role-Based Access Control)
// ============================================================================

func TestUserService_UpdateUserRole_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	user := createTestUser(1, "user@example.com", models.RoleCustomer, true)
	mockRepo.On("FindByID", uint(1)).Return(user, nil).Once()
	mockRepo.On("Update", user).Return(nil).Once()

	err := service.UpdateUserRole(1, models.RoleAdmin)

	require.NoError(t, err, "Should not return error")
	assert.Equal(t, models.RoleAdmin, user.Role, "User role should be updated to admin")
	mockRepo.AssertExpectations(t)
}

func TestUserService_UpdateUserRole_CustomerToAdmin(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	user := createTestUser(1, "customer@example.com", models.RoleCustomer, true)
	mockRepo.On("FindByID", uint(1)).Return(user, nil).Once()
	mockRepo.On("Update", user).Return(nil).Once()

	err := service.UpdateUserRole(1, models.RoleAdmin)

	require.NoError(t, err)
	assert.Equal(t, models.RoleAdmin, user.Role)
	mockRepo.AssertExpectations(t)
}

func TestUserService_UpdateUserRole_AdminToCustomer(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	user := createTestUser(1, "admin@example.com", models.RoleAdmin, true)
	mockRepo.On("FindByID", uint(1)).Return(user, nil).Once()
	mockRepo.On("Update", user).Return(nil).Once()

	err := service.UpdateUserRole(1, models.RoleCustomer)

	require.NoError(t, err)
	assert.Equal(t, models.RoleCustomer, user.Role)
	mockRepo.AssertExpectations(t)
}

func TestUserService_UpdateUserRole_InvalidRole(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	invalidRole := models.UserRole("invalid_role")

	err := service.UpdateUserRole(1, invalidRole)

	assert.Error(t, err, "Should return error for invalid role")
	assert.Contains(t, err.Error(), "invalid role", "Error should contain 'invalid role'")
	assert.Contains(t, err.Error(), "invalid_role", "Error should contain invalid role name")
}

func TestUserService_UpdateUserRole_UserNotFound(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	mockRepo.On("FindByID", uint(999)).Return(nil, nil).Once()

	err := service.UpdateUserRole(999, models.RoleAdmin)

	assert.Error(t, err, "Should return error for non-existent user")
	assert.Contains(t, err.Error(), "user not found", "Error should contain 'user not found'")
	mockRepo.AssertExpectations(t)
}

func TestUserService_UpdateUserRole_RepositoryErrorOnFind(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	expectedErr := errors.New("database connection failed")
	mockRepo.On("FindByID", uint(1)).Return(nil, expectedErr).Once()

	err := service.UpdateUserRole(1, models.RoleAdmin)

	assert.Error(t, err, "Should return error")
	assert.Contains(t, err.Error(), "failed to get user", "Error should contain 'failed to get user'")
	mockRepo.AssertExpectations(t)
}

func TestUserService_UpdateUserRole_RepositoryErrorOnUpdate(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	user := createTestUser(1, "user@example.com", models.RoleCustomer, true)
	expectedErr := errors.New("update failed")
	mockRepo.On("FindByID", uint(1)).Return(user, nil).Once()
	mockRepo.On("Update", user).Return(expectedErr).Once()

	err := service.UpdateUserRole(1, models.RoleAdmin)

	assert.Error(t, err, "Should return error")
	assert.Contains(t, err.Error(), "failed to update user role", "Error should contain 'failed to update user role'")
	mockRepo.AssertExpectations(t)
}

func TestUserService_UpdateUserRole_EmptyRole(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	emptyRole := models.UserRole("")

	err := service.UpdateUserRole(1, emptyRole)

	assert.Error(t, err, "Should return error for empty role")
	assert.Contains(t, err.Error(), "invalid role", "Error should contain 'invalid role'")
}

// ============================================================================
// DEACTIVATE USER TESTS (FR-061: Profile Management)
// ============================================================================

func TestUserService_DeactivateUser_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	user := createTestUser(1, "user@example.com", models.RoleCustomer, true)
	mockRepo.On("FindByID", uint(1)).Return(user, nil).Once()
	mockRepo.On("Update", user).Return(nil).Once()

	err := service.DeactivateUser(1)

	require.NoError(t, err, "Should not return error")
	assert.False(t, user.IsActive, "User should be deactivated")
	mockRepo.AssertExpectations(t)
}

func TestUserService_DeactivateUser_AlreadyInactive(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	user := createTestUser(1, "inactive@example.com", models.RoleCustomer, false)
	mockRepo.On("FindByID", uint(1)).Return(user, nil).Once()
	mockRepo.On("Update", user).Return(nil).Once()

	err := service.DeactivateUser(1)

	require.NoError(t, err, "Should not return error even if already inactive")
	assert.False(t, user.IsActive, "User should remain inactive")
	mockRepo.AssertExpectations(t)
}

func TestUserService_DeactivateUser_UserNotFound(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	mockRepo.On("FindByID", uint(999)).Return(nil, nil).Once()

	err := service.DeactivateUser(999)

	assert.Error(t, err, "Should return error for non-existent user")
	assert.Contains(t, err.Error(), "user not found", "Error should contain 'user not found'")
	mockRepo.AssertExpectations(t)
}

func TestUserService_DeactivateUser_RepositoryErrorOnFind(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	expectedErr := errors.New("database connection failed")
	mockRepo.On("FindByID", uint(1)).Return(nil, expectedErr).Once()

	err := service.DeactivateUser(1)

	assert.Error(t, err, "Should return error")
	assert.Contains(t, err.Error(), "failed to get user", "Error should contain 'failed to get user'")
	mockRepo.AssertExpectations(t)
}

func TestUserService_DeactivateUser_RepositoryErrorOnUpdate(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	user := createTestUser(1, "user@example.com", models.RoleCustomer, true)
	expectedErr := errors.New("update failed")
	mockRepo.On("FindByID", uint(1)).Return(user, nil).Once()
	mockRepo.On("Update", user).Return(expectedErr).Once()

	err := service.DeactivateUser(1)

	assert.Error(t, err, "Should return error")
	assert.Contains(t, err.Error(), "failed to deactivate user", "Error should contain 'failed to deactivate user'")
	mockRepo.AssertExpectations(t)
}

func TestUserService_DeactivateUser_AdminUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	adminUser := createTestUser(1, "admin@example.com", models.RoleAdmin, true)
	mockRepo.On("FindByID", uint(1)).Return(adminUser, nil).Once()
	mockRepo.On("Update", adminUser).Return(nil).Once()

	err := service.DeactivateUser(1)

	require.NoError(t, err, "Should allow deactivating admin user")
	assert.False(t, adminUser.IsActive)
	mockRepo.AssertExpectations(t)
}

// ============================================================================
// ACTIVATE USER TESTS (FR-061: Profile Management)
// ============================================================================

func TestUserService_ActivateUser_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	user := createTestUser(1, "user@example.com", models.RoleCustomer, false)
	mockRepo.On("FindByID", uint(1)).Return(user, nil).Once()
	mockRepo.On("Update", user).Return(nil).Once()

	err := service.ActivateUser(1)

	require.NoError(t, err, "Should not return error")
	assert.True(t, user.IsActive, "User should be activated")
	mockRepo.AssertExpectations(t)
}

func TestUserService_ActivateUser_AlreadyActive(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	user := createTestUser(1, "active@example.com", models.RoleCustomer, true)
	mockRepo.On("FindByID", uint(1)).Return(user, nil).Once()
	mockRepo.On("Update", user).Return(nil).Once()

	err := service.ActivateUser(1)

	require.NoError(t, err, "Should not return error even if already active")
	assert.True(t, user.IsActive, "User should remain active")
	mockRepo.AssertExpectations(t)
}

func TestUserService_ActivateUser_UserNotFound(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	mockRepo.On("FindByID", uint(999)).Return(nil, nil).Once()

	err := service.ActivateUser(999)

	assert.Error(t, err, "Should return error for non-existent user")
	assert.Contains(t, err.Error(), "user not found", "Error should contain 'user not found'")
	mockRepo.AssertExpectations(t)
}

func TestUserService_ActivateUser_RepositoryErrorOnFind(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	expectedErr := errors.New("database connection failed")
	mockRepo.On("FindByID", uint(1)).Return(nil, expectedErr).Once()

	err := service.ActivateUser(1)

	assert.Error(t, err, "Should return error")
	assert.Contains(t, err.Error(), "failed to get user", "Error should contain 'failed to get user'")
	mockRepo.AssertExpectations(t)
}

func TestUserService_ActivateUser_RepositoryErrorOnUpdate(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	user := createTestUser(1, "user@example.com", models.RoleCustomer, false)
	expectedErr := errors.New("update failed")
	mockRepo.On("FindByID", uint(1)).Return(user, nil).Once()
	mockRepo.On("Update", user).Return(expectedErr).Once()

	err := service.ActivateUser(1)

	assert.Error(t, err, "Should return error")
	assert.Contains(t, err.Error(), "failed to activate user", "Error should contain 'failed to activate user'")
	mockRepo.AssertExpectations(t)
}

func TestUserService_ActivateUser_AdminUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	adminUser := createTestUser(1, "admin@example.com", models.RoleAdmin, false)
	mockRepo.On("FindByID", uint(1)).Return(adminUser, nil).Once()
	mockRepo.On("Update", adminUser).Return(nil).Once()

	err := service.ActivateUser(1)

	require.NoError(t, err, "Should allow activating admin user")
	assert.True(t, adminUser.IsActive)
	mockRepo.AssertExpectations(t)
}

// ============================================================================
// GET USER STATS TESTS (FR-061: Profile Management)
// ============================================================================

func TestUserService_GetUserStats_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	stats, err := service.GetUserStats()

	require.NoError(t, err, "Should not return error")
	assert.NotNil(t, stats, "Stats should not be nil")

	// Verify stats structure
	assert.Contains(t, stats, "total_users", "Stats should contain total_users")
	assert.Contains(t, stats, "active_users", "Stats should contain active_users")
	assert.Contains(t, stats, "inactive_users", "Stats should contain inactive_users")
	assert.Contains(t, stats, "admin_users", "Stats should contain admin_users")
	assert.Contains(t, stats, "customer_users", "Stats should contain customer_users")

	// Verify stats values (placeholder implementation returns 0)
	assert.Equal(t, 0, stats["total_users"], "Total users should be 0")
	assert.Equal(t, 0, stats["active_users"], "Active users should be 0")
	assert.Equal(t, 0, stats["inactive_users"], "Inactive users should be 0")
	assert.Equal(t, 0, stats["admin_users"], "Admin users should be 0")
	assert.Equal(t, 0, stats["customer_users"], "Customer users should be 0")
}

func TestUserService_GetUserStats_MultipleCalls(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	// Call GetUserStats multiple times to ensure consistency
	for i := 0; i < 5; i++ {
		stats, err := service.GetUserStats()

		require.NoError(t, err)
		assert.NotNil(t, stats)
		assert.Equal(t, 0, stats["total_users"])
	}
}

// ============================================================================
// INTEGRATION TESTS
// ============================================================================

func TestUserService_UserLifecycle(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	// Create a user
	user := createTestUser(1, "lifecycle@example.com", models.RoleCustomer, true)

	// 1. Get user by ID
	mockRepo.On("FindByID", uint(1)).Return(user, nil).Once()
	foundUser, err := service.GetUserByID(1)
	require.NoError(t, err)
	assert.Equal(t, user.Email, foundUser.Email)

	// 2. Update user role to admin
	mockRepo.On("FindByID", uint(1)).Return(user, nil).Once()
	mockRepo.On("Update", user).Return(nil).Once()
	err = service.UpdateUserRole(1, models.RoleAdmin)
	require.NoError(t, err)
	assert.Equal(t, models.RoleAdmin, user.Role)

	// 3. Deactivate user
	mockRepo.On("FindByID", uint(1)).Return(user, nil).Once()
	mockRepo.On("Update", user).Return(nil).Once()
	err = service.DeactivateUser(1)
	require.NoError(t, err)
	assert.False(t, user.IsActive)

	// 4. Reactivate user
	mockRepo.On("FindByID", uint(1)).Return(user, nil).Once()
	mockRepo.On("Update", user).Return(nil).Once()
	err = service.ActivateUser(1)
	require.NoError(t, err)
	assert.True(t, user.IsActive)

	// 5. Update role back to customer
	mockRepo.On("FindByID", uint(1)).Return(user, nil).Once()
	mockRepo.On("Update", user).Return(nil).Once()
	err = service.UpdateUserRole(1, models.RoleCustomer)
	require.NoError(t, err)
	assert.Equal(t, models.RoleCustomer, user.Role)

	mockRepo.AssertExpectations(t)
}

func TestUserService_MultipleUsersManagement(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	users := []*models.User{
		createTestUser(1, "user1@example.com", models.RoleCustomer, true),
		createTestUser(2, "user2@example.com", models.RoleCustomer, true),
		createTestUser(3, "admin@example.com", models.RoleAdmin, true),
	}

	// Update each user's role
	for _, user := range users {
		newRole := models.RoleAdmin
		if user.Role == models.RoleAdmin {
			newRole = models.RoleCustomer
		}

		mockRepo.On("FindByID", user.ID).Return(user, nil).Once()
		mockRepo.On("Update", user).Return(nil).Once()

		err := service.UpdateUserRole(user.ID, newRole)
		require.NoError(t, err)
		assert.Equal(t, newRole, user.Role)
	}

	mockRepo.AssertExpectations(t)
}

// ============================================================================
// EDGE CASES AND BOUNDARY TESTS
// ============================================================================

func TestUserService_GetUserByID_VeryLargeID(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	veryLargeID := uint(999999999)
	mockRepo.On("FindByID", veryLargeID).Return(nil, nil).Once()

	user, err := service.GetUserByID(veryLargeID)

	assert.Error(t, err, "Should return error for non-existent large ID")
	assert.Nil(t, user)
	mockRepo.AssertExpectations(t)
}

func TestUserService_UpdateUserRole_SameRole(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	user := createTestUser(1, "user@example.com", models.RoleCustomer, true)
	mockRepo.On("FindByID", uint(1)).Return(user, nil).Once()
	mockRepo.On("Update", user).Return(nil).Once()

	// Update to same role
	err := service.UpdateUserRole(1, models.RoleCustomer)

	require.NoError(t, err, "Should allow updating to same role")
	assert.Equal(t, models.RoleCustomer, user.Role)
	mockRepo.AssertExpectations(t)
}

func TestUserService_GetUsers_MaximumPagination(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	users, total, err := service.GetUsers(100, 999999, nil)

	assert.NoError(t, err)
	assert.NotNil(t, users)
	assert.Equal(t, int64(0), total)
}

func TestUserService_GetUsers_ZeroPagination(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	users, total, err := service.GetUsers(0, 0, nil)

	assert.NoError(t, err)
	assert.NotNil(t, users)
	assert.Equal(t, int64(0), total)
}

// ============================================================================
// ROLE VALIDATION TESTS (FR-062: Role-Based Access Control)
// ============================================================================

func TestUserService_RoleValidation_AllValidRoles(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	validRoles := []models.UserRole{
		models.RoleAdmin,
		models.RoleCustomer,
	}

	for _, role := range validRoles {
		t.Run(string(role), func(t *testing.T) {
			user := createTestUser(1, "user@example.com", models.RoleCustomer, true)
			mockRepo.On("FindByID", uint(1)).Return(user, nil).Once()
			mockRepo.On("Update", user).Return(nil).Once()

			err := service.UpdateUserRole(1, role)

			assert.NoError(t, err, "Should accept valid role: %s", role)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserService_RoleValidation_InvalidRoles(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	invalidRoles := []models.UserRole{
		"superadmin",
		"moderator",
		"guest",
		"user",
		"",
		"UNKNOWN",
	}

	for _, role := range invalidRoles {
		t.Run(string(role), func(t *testing.T) {
			err := service.UpdateUserRole(1, role)

			assert.Error(t, err, "Should reject invalid role: %s", role)
			assert.Contains(t, err.Error(), "invalid role")
		})
	}
}

// ============================================================================
// CONCURRENT OPERATIONS TESTS
// ============================================================================

func TestUserService_ConcurrentRoleUpdates(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	user := createTestUser(1, "user@example.com", models.RoleCustomer, true)
	mockRepo.On("FindByID", uint(1)).Return(user, nil).Times(10)
	mockRepo.On("Update", user).Return(nil).Times(10)

	// Simulate concurrent role updates
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func() {
			err := service.UpdateUserRole(1, models.RoleAdmin)
			assert.NoError(t, err)
			done <- true
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}

	mockRepo.AssertExpectations(t)
}

// ============================================================================
// MOCK VERIFICATION TESTS
// ============================================================================

func TestUserService_RepositoryCalls(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	// Test GetUserByID calls repository
	user := createTestUser(1, "test@example.com", models.RoleCustomer, true)
	mockRepo.On("FindByID", uint(1)).Return(user, nil).Once()
	_, err := service.GetUserByID(1)
	require.NoError(t, err)
	mockRepo.AssertExpectations(t)

	// Test UpdateUserRole calls repository
	mockRepo.On("FindByID", uint(1)).Return(user, nil).Once()
	mockRepo.On("Update", user).Return(nil).Once()
	err = service.UpdateUserRole(1, models.RoleAdmin)
	require.NoError(t, err)
	mockRepo.AssertExpectations(t)

	// Test DeactivateUser calls repository
	mockRepo.On("FindByID", uint(1)).Return(user, nil).Once()
	mockRepo.On("Update", user).Return(nil).Once()
	err = service.DeactivateUser(1)
	require.NoError(t, err)
	mockRepo.AssertExpectations(t)

	// Test ActivateUser calls repository
	mockRepo.On("FindByID", uint(1)).Return(user, nil).Once()
	mockRepo.On("Update", user).Return(nil).Once()
	err = service.ActivateUser(1)
	require.NoError(t, err)
	mockRepo.AssertExpectations(t)

	// Test GetUserStats doesn't call repository (placeholder implementation)
	stats, err := service.GetUserStats()
	require.NoError(t, err)
	assert.NotNil(t, stats)
	mockRepo.AssertNotCalled(t, "FindByID")
	mockRepo.AssertNotCalled(t, "Update")
}
