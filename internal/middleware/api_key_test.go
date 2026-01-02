package middleware

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultAPIKeyConfig(t *testing.T) {
	config := DefaultAPIKeyConfig()

	assert.Equal(t, 32, config.KeyLength)
	assert.Equal(t, "kar_", config.KeyPrefix)
	assert.Equal(t, 90*24*time.Hour, config.KeyExpiration)
	assert.Equal(t, "X-API-Key", config.KeyHeader)
	assert.Equal(t, "api_key", config.KeyQueryParam)
	assert.False(t, config.RotationEnabled)
	assert.Equal(t, 30*24*time.Hour, config.RotationPeriod)
	assert.Equal(t, 3, config.MaxKeyVersions)
}

func TestNewAPIKeyManager(t *testing.T) {
	config := DefaultAPIKeyConfig()
	manager := NewAPIKeyManager(config)

	require.NotNil(t, manager)
	assert.Equal(t, config, manager.config)
	assert.NotNil(t, manager.keys)
	assert.NotNil(t, manager.keyVersions)
	assert.NotNil(t, manager.done)

	// Cleanup
	manager.Stop()
}

func TestAPIKeyManager_GenerateKey(t *testing.T) {
	config := DefaultAPIKeyConfig()
	manager := NewAPIKeyManager(config)
	defer manager.Stop()

	name := "Test Key"
	scopes := []string{"read", "write"}
	createdBy := "admin"

	key, keyInfo, err := manager.GenerateKey(name, scopes, createdBy)

	require.NoError(t, err)
	assert.NotEmpty(t, key)
	assert.True(t, len(key) > len(config.KeyPrefix))
	assert.True(t, strings.HasPrefix(key, config.KeyPrefix))

	// Verify key info
	assert.NotEmpty(t, keyInfo.ID)
	assert.Equal(t, name, keyInfo.Name)
	assert.Equal(t, config.KeyPrefix, keyInfo.Prefix)
	assert.Equal(t, 1, keyInfo.Version)
	assert.True(t, keyInfo.IsActive)
	assert.Equal(t, scopes, keyInfo.Scopes)
	assert.Equal(t, createdBy, keyInfo.CreatedBy)
	assert.True(t, keyInfo.ExpiresAt.After(time.Now()))
}

func TestAPIKeyManager_ValidateKey(t *testing.T) {
	config := DefaultAPIKeyConfig()
	manager := NewAPIKeyManager(config)
	defer manager.Stop()

	name := "Test Key"
	scopes := []string{"read"}
	createdBy := "admin"

	key, _, err := manager.GenerateKey(name, scopes, createdBy)
	require.NoError(t, err)

	// Valid key
	keyInfo, err := manager.ValidateKey(key)
	require.NoError(t, err)
	assert.NotNil(t, keyInfo)
	assert.Equal(t, name, keyInfo.Name)

	// Invalid key
	_, err = manager.ValidateKey("invalid_key")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Invalid API key")
}

func TestAPIKeyManager_RotateKey(t *testing.T) {
	config := DefaultAPIKeyConfig()
	manager := NewAPIKeyManager(config)
	defer manager.Stop()

	name := "Test Key"
	scopes := []string{"read"}
	createdBy := "admin"

	oldKey, oldKeyInfo, err := manager.GenerateKey(name, scopes, createdBy)
	require.NoError(t, err)

	// Rotate key
	newKey, newKeyInfo, err := manager.RotateKey(oldKeyInfo.ID)
	require.NoError(t, err)

	assert.NotEqual(t, oldKey, newKey)
	assert.Equal(t, oldKeyInfo.ID, newKeyInfo.ID)
	assert.Equal(t, oldKeyInfo.Name, newKeyInfo.Name)
	assert.Equal(t, oldKeyInfo.Version+1, newKeyInfo.Version)
	assert.True(t, newKeyInfo.IsActive)

	// Old key should be inactive
	_, err = manager.ValidateKey(oldKey)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "inactive")

	// New key should be valid
	keyInfo, err := manager.ValidateKey(newKey)
	require.NoError(t, err)
	assert.True(t, keyInfo.IsActive)
}

func TestAPIKeyManager_RevokeKey(t *testing.T) {
	config := DefaultAPIKeyConfig()
	manager := NewAPIKeyManager(config)
	defer manager.Stop()

	name := "Test Key"
	scopes := []string{"read"}
	createdBy := "admin"

	key, keyInfo, err := manager.GenerateKey(name, scopes, createdBy)
	require.NoError(t, err)

	// Key should be valid
	_, err = manager.ValidateKey(key)
	require.NoError(t, err)

	// Revoke key
	err = manager.RevokeKey(keyInfo.ID)
	require.NoError(t, err)

	// Key should be inactive
	_, err = manager.ValidateKey(key)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "inactive")
}

func TestAPIKeyManager_GetKeyInfo(t *testing.T) {
	config := DefaultAPIKeyConfig()
	manager := NewAPIKeyManager(config)
	defer manager.Stop()

	name := "Test Key"
	scopes := []string{"read"}
	createdBy := "admin"

	_, keyInfo, err := manager.GenerateKey(name, scopes, createdBy)
	require.NoError(t, err)

	// Get key info
	retrievedInfo, err := manager.GetKeyInfo(keyInfo.ID)
	require.NoError(t, err)
	assert.Equal(t, keyInfo.ID, retrievedInfo.ID)
	assert.Equal(t, keyInfo.Name, retrievedInfo.Name)
	assert.Equal(t, keyInfo.Version, retrievedInfo.Version)
}

func TestAPIKeyManager_ListKeys(t *testing.T) {
	config := DefaultAPIKeyConfig()
	manager := NewAPIKeyManager(config)
	defer manager.Stop()

	// Generate multiple keys
	for i := 0; i < 3; i++ {
		name := fmt.Sprintf("Test Key %d", i)
		scopes := []string{"read"}
		createdBy := "admin"
		_, _, err := manager.GenerateKey(name, scopes, createdBy)
		require.NoError(t, err)
	}

	// List keys
	keys := manager.ListKeys()
	assert.Len(t, keys, 3)
}

func TestAPIKeyManager_KeyExpiration(t *testing.T) {
	config := APIKeyConfig{
		KeyLength:     32,
		KeyPrefix:     "test_",
		KeyExpiration: 100 * time.Millisecond, // Short expiration for testing
	}
	manager := NewAPIKeyManager(config)
	defer manager.Stop()

	name := "Test Key"
	scopes := []string{"read"}
	createdBy := "admin"

	key, _, err := manager.GenerateKey(name, scopes, createdBy)
	require.NoError(t, err)

	// Key should be valid immediately
	_, err = manager.ValidateKey(key)
	require.NoError(t, err)

	// Wait for expiration
	time.Sleep(150 * time.Millisecond)

	// Key should be expired
	_, err = manager.ValidateKey(key)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "expired")
}

func TestAPIKeyManager_MaxKeyVersions(t *testing.T) {
	config := APIKeyConfig{
		KeyLength:      32,
		KeyPrefix:      "test_",
		KeyExpiration:  1 * time.Hour,
		MaxKeyVersions: 2,
	}
	manager := NewAPIKeyManager(config)
	defer manager.Stop()

	name := "Test Key"
	scopes := []string{"read"}
	createdBy := "admin"

	_, keyInfo, err := manager.GenerateKey(name, scopes, createdBy)
	require.NoError(t, err)

	// Rotate key multiple times
	for i := 0; i < 3; i++ {
		_, _, err = manager.RotateKey(keyInfo.ID)
		require.NoError(t, err)
	}

	// Should only keep max versions
	versions := manager.keyVersions[keyInfo.ID]
	assert.LessOrEqual(t, len(versions), config.MaxKeyVersions)
}

func TestAPIKeyAuthMiddleware(t *testing.T) {
	app := fiber.New()
	config := DefaultAPIKeyConfig()
	manager := NewAPIKeyManager(config)
	defer manager.Stop()

	// Generate API key
	name := "Test Key"
	scopes := []string{"read"}
	createdBy := "admin"
	key, _, err := manager.GenerateKey(name, scopes, createdBy)
	require.NoError(t, err)

	// Add middleware
	app.Use(APIKeyAuth(manager, config))

	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	// Test with valid API key in header
	req := getTestRequest("GET", "/test", nil)
	req.Header.Set(config.KeyHeader, key)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	// Test with valid API key in query parameter
	req = getTestRequest("GET", "/test?api_key="+key, nil)
	resp, err = app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	// Test without API key
	req = getTestRequest("GET", "/test", nil)
	resp, err = app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 401, resp.StatusCode)

	// Test with invalid API key
	req = getTestRequest("GET", "/test", nil)
	req.Header.Set(config.KeyHeader, "invalid_key")
	resp, err = app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 401, resp.StatusCode)
}

func TestRequireScopeMiddleware(t *testing.T) {
	app := fiber.New()
	config := DefaultAPIKeyConfig()
	manager := NewAPIKeyManager(config)
	defer manager.Stop()

	// Generate API key with scopes
	name := "Test Key"
	scopes := []string{"read", "write"}
	createdBy := "admin"
	key, _, err := manager.GenerateKey(name, scopes, createdBy)
	require.NoError(t, err)

	// Add middleware
	app.Use(APIKeyAuth(manager, config))

	app.Get("/test", RequireScope("read"), func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	app.Post("/test", RequireScope("write"), func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	app.Delete("/test", RequireScope("delete"), func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	// Test with required scope
	req := getTestRequest("GET", "/test", nil)
	req.Header.Set(config.KeyHeader, key)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	// Test with another required scope
	req = getTestRequest("POST", "/test", nil)
	req.Header.Set(config.KeyHeader, key)
	resp, err = app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	// Test with missing scope
	req = getTestRequest("DELETE", "/test", nil)
	req.Header.Set(config.KeyHeader, key)
	resp, err = app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 403, resp.StatusCode)
}

func TestAPIKeyMiddleware(t *testing.T) {
	config := DefaultAPIKeyConfig()
	manager := NewAPIKeyManager(config)
	defer manager.Stop()

	middleware := APIKeyMiddleware(manager)
	require.NotNil(t, middleware)
}

// Helper function to create test requests
func getTestRequest(method, path string, body []byte) *http.Request {
	req := httptest.NewRequest(method, path, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	return req
}
