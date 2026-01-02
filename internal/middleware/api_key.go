package middleware

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"karima_store/internal/errors"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
)

// APIKeyConfig holds configuration for API key management
type APIKeyConfig struct {
	KeyLength       int           // Length of API key in bytes
	KeyPrefix       string        // Prefix for API keys (e.g., "kar_")
	KeyExpiration   time.Duration // Default key expiration time
	KeyHeader       string        // Header name for API key
	KeyQueryParam   string        // Query parameter name for API key
	RotationEnabled bool          // Enable automatic key rotation
	RotationPeriod  time.Duration // Rotation period for automatic rotation
	MaxKeyVersions  int           // Maximum number of key versions to keep
}

// DefaultAPIKeyConfig returns default API key configuration
func DefaultAPIKeyConfig() APIKeyConfig {
	return APIKeyConfig{
		KeyLength:       32,
		KeyPrefix:       "kar_",
		KeyExpiration:   90 * 24 * time.Hour, // 90 days
		KeyHeader:       "X-API-Key",
		KeyQueryParam:   "api_key",
		RotationEnabled: false,
		RotationPeriod:  30 * 24 * time.Hour, // 30 days
		MaxKeyVersions:  3,
	}
}

// APIKeyInfo holds information about an API key
type APIKeyInfo struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	KeyHash     string    `json:"key_hash"`
	Prefix      string    `json:"prefix"`
	Version     int       `json:"version"`
	IsActive    bool      `json:"is_active"`
	Scopes      []string  `json:"scopes"`
	CreatedAt   time.Time `json:"created_at"`
	ExpiresAt   time.Time `json:"expires_at"`
	RotatedAt   time.Time `json:"rotated_at"`
	LastUsedAt  time.Time `json:"last_used_at"`
	CreatedBy   string    `json:"created_by"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// APIKeyManager manages API keys
type APIKeyManager struct {
	config      APIKeyConfig
	keys        map[string]*APIKeyInfo // Map of key hash to key info
	keyVersions map[string][]string    // Map of key ID to list of key hashes (versions)
	mu          sync.RWMutex
	cleanupTick *time.Ticker
	done        chan struct{}
}

// NewAPIKeyManager creates a new API key manager
func NewAPIKeyManager(config APIKeyConfig) *APIKeyManager {
	if config.KeyLength == 0 {
		config = DefaultAPIKeyConfig()
	}

	manager := &APIKeyManager{
		config:      config,
		keys:        make(map[string]*APIKeyInfo),
		keyVersions: make(map[string][]string),
		done:        make(chan struct{}),
	}

	// Start cleanup goroutine
	manager.startCleanup()

	return manager
}

// startCleanup starts the cleanup goroutine to remove expired keys
func (m *APIKeyManager) startCleanup() {
	m.cleanupTick = time.NewTicker(1 * time.Hour)
	go func() {
		for {
			select {
			case <-m.cleanupTick.C:
				m.cleanupExpiredKeys()
			case <-m.done:
				return
			}
		}
	}()
}

// cleanupExpiredKeys removes expired keys from memory
func (m *APIKeyManager) cleanupExpiredKeys() {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	for keyHash, keyInfo := range m.keys {
		if now.After(keyInfo.ExpiresAt) && !keyInfo.IsActive {
			// Remove expired and inactive keys
			delete(m.keys, keyHash)
		}
	}
}

// Stop stops the cleanup goroutine
func (m *APIKeyManager) Stop() {
	close(m.done)
	if m.cleanupTick != nil {
		m.cleanupTick.Stop()
	}
}

// GenerateKey generates a new API key
func (m *APIKeyManager) GenerateKey(name string, scopes []string, createdBy string) (string, *APIKeyInfo, error) {
	// Generate random bytes
	bytes := make([]byte, m.config.KeyLength)
	if _, err := rand.Read(bytes); err != nil {
		return "", nil, fmt.Errorf("failed to generate API key: %w", err)
	}

	// Encode to base64
	key := base64.URLEncoding.EncodeToString(bytes)

	// Add prefix
	fullKey := m.config.KeyPrefix + key

	// Calculate hash of the key
	hash := sha256.Sum256([]byte(fullKey))
	keyHash := hex.EncodeToString(hash[:])

	// Generate key ID
	keyID := generateKeyID()

	// Create key info
	keyInfo := &APIKeyInfo{
		ID:         keyID,
		Name:       name,
		KeyHash:    keyHash,
		Prefix:     m.config.KeyPrefix,
		Version:    1,
		IsActive:   true,
		Scopes:     scopes,
		CreatedAt:  time.Now(),
		ExpiresAt:  time.Now().Add(m.config.KeyExpiration),
		CreatedBy:  createdBy,
		Metadata:   make(map[string]interface{}),
	}

	// Store key
	m.mu.Lock()
	m.keys[keyHash] = keyInfo
	m.keyVersions[keyID] = []string{keyHash}
	m.mu.Unlock()

	return fullKey, keyInfo, nil
}

// ValidateKey validates an API key and returns key info
func (m *APIKeyManager) ValidateKey(key string) (*APIKeyInfo, error) {
	// Remove prefix if present
	fullKey := key
	if strings.HasPrefix(key, m.config.KeyPrefix) {
		fullKey = key
	} else {
		fullKey = m.config.KeyPrefix + key
	}

	// Calculate hash
	hash := sha256.Sum256([]byte(fullKey))
	keyHash := hex.EncodeToString(hash[:])

	// Look up key
	m.mu.RLock()
	keyInfo, exists := m.keys[keyHash]
	m.mu.RUnlock()

	if !exists {
		return nil, errors.NewUnauthorizedError("Invalid API key")
	}

	// Check if key is active
	if !keyInfo.IsActive {
		return nil, errors.NewUnauthorizedError("API key is inactive")
	}

	// Check if key is expired
	if time.Now().After(keyInfo.ExpiresAt) {
		return nil, errors.NewUnauthorizedError("API key has expired")
	}

	// Update last used time
	m.mu.Lock()
	keyInfo.LastUsedAt = time.Now()
	m.mu.Unlock()

	return keyInfo, nil
}

// RotateKey rotates an existing API key
func (m *APIKeyManager) RotateKey(keyID string) (string, *APIKeyInfo, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Get existing key versions
	versions, exists := m.keyVersions[keyID]
	if !exists || len(versions) == 0 {
		return "", nil, fmt.Errorf("API key not found")
	}

	// Get current key info
	currentHash := versions[len(versions)-1]
	currentKeyInfo := m.keys[currentHash]
	if currentKeyInfo == nil {
		return "", nil, fmt.Errorf("API key not found")
	}

	// Deactivate current key
	currentKeyInfo.IsActive = false
	currentKeyInfo.RotatedAt = time.Now()

	// Generate new key
	bytes := make([]byte, m.config.KeyLength)
	if _, err := rand.Read(bytes); err != nil {
		return "", nil, fmt.Errorf("failed to generate API key: %w", err)
	}

	key := base64.URLEncoding.EncodeToString(bytes)
	fullKey := m.config.KeyPrefix + key

	// Calculate hash
	hash := sha256.Sum256([]byte(fullKey))
	newHash := hex.EncodeToString(hash[:])

	// Create new key info
	newKeyInfo := &APIKeyInfo{
		ID:         keyID,
		Name:       currentKeyInfo.Name,
		KeyHash:    newHash,
		Prefix:     m.config.KeyPrefix,
		Version:    currentKeyInfo.Version + 1,
		IsActive:   true,
		Scopes:     currentKeyInfo.Scopes,
		CreatedAt:  time.Now(),
		ExpiresAt:  time.Now().Add(m.config.KeyExpiration),
		RotatedAt:  time.Now(),
		CreatedBy:  currentKeyInfo.CreatedBy,
		Metadata:   currentKeyInfo.Metadata,
	}

	// Store new key
	m.keys[newHash] = newKeyInfo
	m.keyVersions[keyID] = append(versions, newHash)

	// Cleanup old versions if exceeding max
	if len(m.keyVersions[keyID]) > m.config.MaxKeyVersions {
		oldVersions := m.keyVersions[keyID][:len(m.keyVersions[keyID])-m.config.MaxKeyVersions]
		for _, oldHash := range oldVersions {
			delete(m.keys, oldHash)
		}
		m.keyVersions[keyID] = m.keyVersions[keyID][len(m.keyVersions[keyID])-m.config.MaxKeyVersions:]
	}

	return fullKey, newKeyInfo, nil
}

// RevokeKey revokes an API key
func (m *APIKeyManager) RevokeKey(keyID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	versions, exists := m.keyVersions[keyID]
	if !exists {
		return fmt.Errorf("API key not found")
	}

	// Deactivate all versions
	for _, hash := range versions {
		if keyInfo, ok := m.keys[hash]; ok {
			keyInfo.IsActive = false
		}
	}

	return nil
}

// GetKeyInfo retrieves information about an API key
func (m *APIKeyManager) GetKeyInfo(keyID string) (*APIKeyInfo, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	versions, exists := m.keyVersions[keyID]
	if !exists || len(versions) == 0 {
		return nil, fmt.Errorf("API key not found")
	}

	// Return the latest version
	latestHash := versions[len(versions)-1]
	keyInfo := m.keys[latestHash]
	if keyInfo == nil {
		return nil, fmt.Errorf("API key not found")
	}

	return keyInfo, nil
}

// ListKeys lists all API keys
func (m *APIKeyManager) ListKeys() []*APIKeyInfo {
	m.mu.RLock()
	defer m.mu.RUnlock()

	keys := make([]*APIKeyInfo, 0, len(m.keyVersions))
	for keyID := range m.keyVersions {
		versions := m.keyVersions[keyID]
		if len(versions) > 0 {
			latestHash := versions[len(versions)-1]
			if keyInfo := m.keys[latestHash]; keyInfo != nil {
				keys = append(keys, keyInfo)
			}
		}
	}

	return keys
}

// APIKeyAuth creates API key authentication middleware
func APIKeyAuth(manager *APIKeyManager, config APIKeyConfig) fiber.Handler {
	if config.KeyLength == 0 {
		config = DefaultAPIKeyConfig()
	}

	return func(c *fiber.Ctx) error {
		// Get API key from header or query parameter
		apiKey := c.Get(config.KeyHeader)
		if apiKey == "" {
			apiKey = c.Query(config.KeyQueryParam)
		}

		if apiKey == "" {
			return errors.NewUnauthorizedError("API key is required")
		}

		// Validate API key
		keyInfo, err := manager.ValidateKey(apiKey)
		if err != nil {
			return err
		}

		// Store key info in context
		c.Locals("api_key_info", keyInfo)
		c.Locals("api_key_id", keyInfo.ID)
		c.Locals("api_key_name", keyInfo.Name)
		c.Locals("api_key_scopes", keyInfo.Scopes)

		return c.Next()
	}
}

// RequireScope creates middleware that requires specific scopes
func RequireScope(requiredScopes ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		keyScopes := c.Locals("api_key_scopes")
		if keyScopes == nil {
			return errors.NewForbiddenError("No API key scopes found")
		}

		scopes, ok := keyScopes.([]string)
		if !ok {
			return errors.NewForbiddenError("Invalid API key scopes")
		}

		// Check if all required scopes are present
		for _, required := range requiredScopes {
			found := false
			for _, scope := range scopes {
				if scope == required {
					found = true
					break
				}
			}
			if !found {
				return errors.NewForbiddenError(fmt.Sprintf("Missing required scope: %s", required))
			}
		}

		return c.Next()
	}
}

// generateKeyID generates a unique key ID
func generateKeyID() string {
	return fmt.Sprintf("key_%d", time.Now().UnixNano())
}

// APIKeyMiddleware creates an API key middleware with default configuration
func APIKeyMiddleware(manager *APIKeyManager) fiber.Handler {
	return APIKeyAuth(manager, DefaultAPIKeyConfig())
}
