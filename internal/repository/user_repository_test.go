package repository

import (
	"testing"

	"github.com/karima-store/internal/models"
	"github.com/karima-store/internal/test_setup"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func setupUserTest(t *testing.T) (*gorm.DB, func()) {
	db, cleanup := test_setup.SetupTestDB(t)

	// Clean up existing data
	db.Exec("DELETE FROM users")

	return db, cleanup
}

func createTestUser(email string) *models.User {
	return &models.User{
		FullName: "Test User",
		Email:    email,
		KratosID: "kratos_" + email, // Ensure unique KratosID
		Password: "hashedpassword",
		Role:     models.RoleCustomer,
		Phone:    "08123456789",
		Address:  "123 Test Street",
		City:     "Jakarta",
		Province: "DKI Jakarta",
	}
}

func TestUserRepository_NewUserRepository(t *testing.T) {
	db, cleanup := setupUserTest(t)
	defer cleanup()

	repo := NewUserRepository(db)
	assert.NotNil(t, repo)
}

func TestUserRepository_Create(t *testing.T) {
	db, cleanup := setupUserTest(t)
	defer cleanup()

	repo := NewUserRepository(db)

	user := createTestUser("test@example.com")
	err := repo.Create(user)
	require.NoError(t, err)
	assert.NotZero(t, user.ID)
}

func TestUserRepository_Create_DuplicateEmail(t *testing.T) {
	db, cleanup := setupUserTest(t)
	defer cleanup()

	repo := NewUserRepository(db)

	// Create first user
	user1 := createTestUser("duplicate@example.com")
	err := repo.Create(user1)
	require.NoError(t, err)

	// Try to create second user with same email
	user2 := createTestUser("duplicate@example.com")
	err = repo.Create(user2)
	assert.Error(t, err) // Should fail due to unique constraint
}

func TestUserRepository_FindByID(t *testing.T) {
	db, cleanup := setupUserTest(t)
	defer cleanup()

	repo := NewUserRepository(db)

	// Create user
	user := createTestUser("findbyid@example.com")
	err := repo.Create(user)
	require.NoError(t, err)

	// Find by ID
	fetched, err := repo.FindByID(user.ID)
	require.NoError(t, err)
	assert.NotNil(t, fetched)
	assert.Equal(t, user.ID, fetched.ID)
	assert.Equal(t, "findbyid@example.com", fetched.Email)
}

func TestUserRepository_FindByID_NotFound(t *testing.T) {
	db, cleanup := setupUserTest(t)
	defer cleanup()

	repo := NewUserRepository(db)

	// Find non-existent user
	fetched, err := repo.FindByID(99999)
	require.NoError(t, err) // Returns nil, not error
	assert.Nil(t, fetched)
}

func TestUserRepository_FindByEmail(t *testing.T) {
	db, cleanup := setupUserTest(t)
	defer cleanup()

	repo := NewUserRepository(db)

	// Create user
	user := createTestUser("findbyemail@example.com")
	err := repo.Create(user)
	require.NoError(t, err)

	// Find by email
	fetched, err := repo.FindByEmail("findbyemail@example.com")
	require.NoError(t, err)
	assert.NotNil(t, fetched)
	assert.Equal(t, user.ID, fetched.ID)
	assert.Equal(t, "Test User", fetched.FullName)
}

func TestUserRepository_FindByEmail_NotFound(t *testing.T) {
	db, cleanup := setupUserTest(t)
	defer cleanup()

	repo := NewUserRepository(db)

	// Find non-existent email
	fetched, err := repo.FindByEmail("nonexistent@example.com")
	require.NoError(t, err) // Returns nil, not error
	assert.Nil(t, fetched)
}

func TestUserRepository_FindByKratosID(t *testing.T) {
	db, cleanup := setupUserTest(t)
	defer cleanup()

	repo := NewUserRepository(db)

	// Create user with KratosID
	user := createTestUser("kratos@example.com")
	user.KratosID = "kratos-uuid-12345"
	err := repo.Create(user)
	require.NoError(t, err)

	// Find by KratosID
	fetched, err := repo.FindByKratosID("kratos-uuid-12345")
	require.NoError(t, err)
	assert.NotNil(t, fetched)
	assert.Equal(t, user.ID, fetched.ID)
	assert.Equal(t, "kratos-uuid-12345", fetched.KratosID)
}

func TestUserRepository_FindByKratosID_NotFound(t *testing.T) {
	db, cleanup := setupUserTest(t)
	defer cleanup()

	repo := NewUserRepository(db)

	// Find non-existent KratosID
	fetched, err := repo.FindByKratosID("non-existent-kratos-id")
	require.NoError(t, err) // Returns nil, not error
	assert.Nil(t, fetched)
}

func TestUserRepository_Update(t *testing.T) {
	db, cleanup := setupUserTest(t)
	defer cleanup()

	repo := NewUserRepository(db)

	// Create user
	user := createTestUser("update@example.com")
	err := repo.Create(user)
	require.NoError(t, err)

	// Update user
	user.FullName = "Updated Name"
	user.Phone = "08199999999"
	err = repo.Update(user)
	require.NoError(t, err)

	// Verify update
	fetched, err := repo.FindByID(user.ID)
	require.NoError(t, err)
	assert.Equal(t, "Updated Name", fetched.FullName)
	assert.Equal(t, "08199999999", fetched.Phone)
}

func TestUserRepository_Update_Role(t *testing.T) {
	db, cleanup := setupUserTest(t)
	defer cleanup()

	repo := NewUserRepository(db)

	// Create user as customer
	user := createTestUser("rolechange@example.com")
	user.Role = models.RoleCustomer
	err := repo.Create(user)
	require.NoError(t, err)

	// Update to admin
	user.Role = models.RoleAdmin
	err = repo.Update(user)
	require.NoError(t, err)

	// Verify
	fetched, err := repo.FindByID(user.ID)
	require.NoError(t, err)
	assert.Equal(t, models.RoleAdmin, fetched.Role)
}

func TestUserRepository_Update_Verification(t *testing.T) {
	db, cleanup := setupUserTest(t)
	defer cleanup()

	repo := NewUserRepository(db)

	// Create unverified user
	user := createTestUser("verify@example.com")
	user.IsVerified = false
	err := repo.Create(user)
	require.NoError(t, err)

	// Verify user
	user.IsVerified = true
	err = repo.Update(user)
	require.NoError(t, err)

	// Check
	fetched, err := repo.FindByID(user.ID)
	require.NoError(t, err)
	assert.True(t, fetched.IsVerified)
}

func TestUserRepository_MultipleUsers(t *testing.T) {
	db, cleanup := setupUserTest(t)
	defer cleanup()

	repo := NewUserRepository(db)

	// Create multiple users
	emails := []string{"user1@example.com", "user2@example.com", "user3@example.com"}
	for _, email := range emails {
		user := createTestUser(email)
		err := repo.Create(user)
		require.NoError(t, err)
	}

	// Verify each user can be found
	for _, email := range emails {
		fetched, err := repo.FindByEmail(email)
		require.NoError(t, err)
		assert.NotNil(t, fetched)
		assert.Equal(t, email, fetched.Email)
	}
}

func TestUserRepository_AddressUpdate(t *testing.T) {
	db, cleanup := setupUserTest(t)
	defer cleanup()

	repo := NewUserRepository(db)

	// Create user
	user := createTestUser("address@example.com")
	err := repo.Create(user)
	require.NoError(t, err)

	// Update address
	user.Address = "456 New Street"
	user.City = "Bandung"
	user.Province = "West Java"
	user.PostalCode = "40100"
	err = repo.Update(user)
	require.NoError(t, err)

	// Verify
	fetched, err := repo.FindByID(user.ID)
	require.NoError(t, err)
	assert.Equal(t, "456 New Street", fetched.Address)
	assert.Equal(t, "Bandung", fetched.City)
	assert.Equal(t, "West Java", fetched.Province)
	assert.Equal(t, "40100", fetched.PostalCode)
}
