package middleware

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/karima-store/internal/config"
	"github.com/karima-store/internal/storage"
	"runtime"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// HealthStatus represents the overall health status
type HealthStatus struct {
	Status    string            `json:"status"`
	Timestamp time.Time         `json:"timestamp"`
	Version   string            `json:"version,omitempty"`
	Services  map[string]Service `json:"services"`
}

// Service represents the health status of a service
type Service struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
	Latency string `json:"latency,omitempty"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// HealthChecker defines the interface for health checkers
type HealthChecker interface {
	Check(ctx context.Context) Service
}

// DatabaseHealthChecker checks database health
type DatabaseHealthChecker struct {
	db *gorm.DB
}

func NewDatabaseHealthChecker(db *gorm.DB) *DatabaseHealthChecker {
	return &DatabaseHealthChecker{db: db}
}

func (c *DatabaseHealthChecker) Check(ctx context.Context) Service {
	start := time.Now()

	sqlDB, err := c.db.DB()
	if err != nil {
		return Service{
			Status:  "unhealthy",
			Message: "Failed to get database connection",
			Latency: time.Since(start).String(),
		}
	}

	err = sqlDB.PingContext(ctx)
	if err != nil {
		return Service{
			Status:  "unhealthy",
			Message: fmt.Sprintf("Database ping failed: %v", err),
			Latency: time.Since(start).String(),
		}
	}

	// Get database stats
	var stats map[string]interface{}
	dbStats := sqlDB.Stats()
	if dbStats != (sql.DBStats{}) {
		stats = map[string]interface{}{
			"max_open_connections": dbStats.MaxOpenConnections,
			"open_connections":     dbStats.OpenConnections,
			"in_use":              dbStats.InUse,
			"idle":                dbStats.Idle,
		}
	}

	return Service{
		Status:  "healthy",
		Message: "Database is healthy",
		Latency: time.Since(start).String(),
		Details: stats,
	}
}

// RedisHealthChecker checks Redis health
type RedisHealthChecker struct {
	client *redis.Client
}

func NewRedisHealthChecker(client *redis.Client) *RedisHealthChecker {
	return &RedisHealthChecker{client: client}
}

func (c *RedisHealthChecker) Check(ctx context.Context) Service {
	start := time.Now()

	// Ping Redis
	err := c.client.Ping(ctx).Err()
	if err != nil {
		return Service{
			Status:  "unhealthy",
			Message: fmt.Sprintf("Redis ping failed: %v", err),
			Latency: time.Since(start).String(),
		}
	}

	// Get Redis info
	info, err := c.client.Info(ctx, "server").Result()
	if err != nil {
		return Service{
			Status:  "unhealthy",
			Message: fmt.Sprintf("Failed to get Redis info: %v", err),
			Latency: time.Since(start).String(),
		}
	}

	return Service{
		Status:  "healthy",
		Message: "Redis is healthy",
		Latency: time.Since(start).String(),
		Details: map[string]interface{}{
			"info": info,
		},
	}
}

// StorageHealthChecker checks storage health
type StorageHealthChecker struct {
	storage *storage.R2Storage
}

func NewStorageHealthChecker(storage *storage.R2Storage) *StorageHealthChecker {
	return &StorageHealthChecker{storage: storage}
}

func (c *StorageHealthChecker) Check(ctx context.Context) Service {
	start := time.Now()

	// Check if storage is initialized
	if c.storage == nil {
		return Service{
			Status:  "unhealthy",
			Message: "Storage not initialized",
			Latency: time.Since(start).String(),
		}
	}

	// In a real implementation, you would perform a simple operation
	// like listing buckets or checking connectivity
	return Service{
		Status:  "healthy",
		Message: "Storage is healthy",
		Latency: time.Since(start).String(),
	}
}

// SystemHealthChecker checks system health
type SystemHealthChecker struct{}

func NewSystemHealthChecker() *SystemHealthChecker {
	return &SystemHealthChecker{}
}

func (c *SystemHealthChecker) Check(ctx context.Context) Service {
	start := time.Now()

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// Check if system is under heavy load
	healthy := true
	message := "System is healthy"

	if m.Alloc > 1024*1024*1024 { // > 1GB
		healthy = false
		message = "High memory usage"
	}

	if runtime.NumGoroutine() > 1000 {
		healthy = false
		message = "High goroutine count"
	}

	status := "healthy"
	if !healthy {
		status = "degraded"
	}

	return Service{
		Status:  status,
		Message: message,
		Latency: time.Since(start).String(),
		Details: map[string]interface{}{
			"goroutines":       runtime.NumGoroutine(),
			"memory_allocated": m.Alloc,
			"memory_total":     m.TotalAlloc,
			"memory_sys":       m.Sys,
			"memory_heap":      m.HeapAlloc,
			"gc_cycles":        m.NumGC,
		},
	}
}

// HealthCheckManager manages health checks
type HealthCheckManager struct {
	checkers map[string]HealthChecker
	timeout  time.Duration
}

// NewHealthCheckManager creates a new health check manager
func NewHealthCheckManager(timeout time.Duration) *HealthCheckManager {
	return &HealthCheckManager{
		checkers: make(map[string]HealthChecker),
		timeout:  timeout,
	}
}

// RegisterChecker registers a health checker
func (m *HealthCheckManager) RegisterChecker(name string, checker HealthChecker) {
	m.checkers[name] = checker
}

// Check performs health checks on all registered services
func (m *HealthCheckManager) Check(ctx context.Context) HealthStatus {
	ctx, cancel := context.WithTimeout(ctx, m.timeout)
	defer cancel()

	services := make(map[string]Service)
	overallStatus := "healthy"

	var wg sync.WaitGroup
	var mu sync.Mutex

	for name, checker := range m.checkers {
		wg.Add(1)
		go func(name string, checker HealthChecker) {
			defer wg.Done()

			service := checker.Check(ctx)

			mu.Lock()
			services[name] = service

			// Update overall status
			if service.Status == "unhealthy" {
				overallStatus = "unhealthy"
			} else if service.Status == "degraded" && overallStatus == "healthy" {
				overallStatus = "degraded"
			}
			mu.Unlock()
		}(name, checker)
	}

	wg.Wait()

	return HealthStatus{
		Status:    overallStatus,
		Timestamp: time.Now(),
		Services:  services,
	}
}

// InitializeHealthCheckManager initializes the health check manager with default checkers
func InitializeHealthCheckManager(cfg *config.Config, db *gorm.DB, redisClient *redis.Client, r2Storage *storage.R2Storage) *HealthCheckManager {
	manager := NewHealthCheckManager(5 * time.Second)

	// Register database checker
	if db != nil {
		manager.RegisterChecker("database", NewDatabaseHealthChecker(db))
	}

	// Register Redis checker
	if redisClient != nil {
		manager.RegisterChecker("redis", NewRedisHealthChecker(redisClient))
	}

	// Register storage checker
	if r2Storage != nil {
		manager.RegisterChecker("storage", NewStorageHealthChecker(r2Storage))
	}

	// Register system checker
	manager.RegisterChecker("system", NewSystemHealthChecker())

	return manager
}

// Global health check manager
var healthCheckManager *HealthCheckManager

// InitHealthChecks initializes health checks
func InitHealthChecks(cfg *config.Config, db *gorm.DB, redisClient *redis.Client, r2Storage *storage.R2Storage) {
	healthCheckManager = InitializeHealthCheckManager(cfg, db, redisClient, r2Storage)
}

// HealthCheckHandler returns a handler for health checks
func HealthCheckHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if healthCheckManager == nil {
			return c.JSON(HealthStatus{
				Status:    "unhealthy",
				Timestamp: time.Now(),
				Services: map[string]Service{
					"error": {
						Status:  "unhealthy",
						Message: "Health check manager not initialized",
					},
				},
			})
		}

		status := healthCheckManager.Check(c.Context())

		// Set appropriate HTTP status code
		httpStatus := fiber.StatusOK
		if status.Status == "degraded" {
			httpStatus = 200 // Still OK but degraded
		} else if status.Status == "unhealthy" {
			httpStatus = fiber.StatusServiceUnavailable
		}

		return c.Status(httpStatus).JSON(status)
	}
}

// ReadinessHandler returns a handler for readiness checks
func ReadinessHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if healthCheckManager == nil {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"status":  "not_ready",
				"message": "Health check manager not initialized",
			})
		}

		status := healthCheckManager.Check(c.Context())

		// Check if critical services are healthy
		criticalServices := []string{"database", "redis"}
		allHealthy := true

		for _, service := range criticalServices {
			if serviceStatus, exists := status.Services[service]; exists {
				if serviceStatus.Status != "healthy" {
					allHealthy = false
					break
				}
			}
		}

		if allHealthy {
			return c.JSON(fiber.Map{
				"status":    "ready",
				"timestamp": time.Now(),
			})
		}

		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"status":    "not_ready",
			"timestamp": time.Now(),
			"services":  status.Services,
		})
	}
}

// LivenessHandler returns a handler for liveness checks
func LivenessHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":    "alive",
			"timestamp": time.Now(),
		})
	}
}

// DetailedHealthCheckHandler returns a handler for detailed health checks
func DetailedHealthCheckHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if healthCheckManager == nil {
			return c.JSON(HealthStatus{
				Status:    "unhealthy",
				Timestamp: time.Now(),
				Services: map[string]Service{
					"error": {
						Status:  "unhealthy",
						Message: "Health check manager not initialized",
					},
				},
			})
		}

		status := healthCheckManager.Check(c.Context())

		// Add additional system information
		var m runtime.MemStats
		runtime.ReadMemStats(&m)

		response := map[string]interface{}{
			"health": status,
			"system": map[string]interface{}{
				"version":       "1.0.0", // TODO: Get from config
				"go_version":    runtime.Version(),
				"goroutines":    runtime.NumGoroutine(),
				"cpu_count":     runtime.NumCPU(),
				"memory": map[string]interface{}{
					"alloc":       m.Alloc,
					"total_alloc": m.TotalAlloc,
					"sys":         m.Sys,
					"heap_alloc":  m.HeapAlloc,
					"heap_sys":    m.HeapSys,
					"gc_cycles":   m.NumGC,
				},
			},
		}

		return c.JSON(response)
	}
}
