package database

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/karima-store/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DBPoolConfig holds configuration for database connection pool
type DBPoolConfig struct {
	MaxOpenConns    int           // Maximum number of open connections to the database
	MaxIdleConns    int           // Maximum number of idle connections in the pool
	ConnMaxLifetime time.Duration // Maximum amount of time a connection may be reused
	ConnMaxIdleTime time.Duration // Maximum amount of time a connection may be idle
	ConnMinIdleTime time.Duration // Minimum amount of time a connection may be idle before being closed
	HealthCheck     bool          // Enable periodic health checks
	HealthInterval  time.Duration // Interval between health checks
	MetricsEnabled  bool          // Enable metrics collection
}

// DefaultDBPoolConfig returns default database pool configuration
func DefaultDBPoolConfig(env string) DBPoolConfig {
	if env == "production" {
		return DBPoolConfig{
			MaxOpenConns:    100,
			MaxIdleConns:    25,
			ConnMaxLifetime: 30 * time.Minute,
			ConnMaxIdleTime: 15 * time.Minute,
			ConnMinIdleTime: 5 * time.Minute,
			HealthCheck:     true,
			HealthInterval:  1 * time.Minute,
			MetricsEnabled:  true,
		}
	}

	return DBPoolConfig{
		MaxOpenConns:    25,
		MaxIdleConns:    10,
		ConnMaxLifetime: 1 * time.Hour,
		ConnMaxIdleTime: 30 * time.Minute,
		ConnMinIdleTime: 5 * time.Minute,
		HealthCheck:     true,
		HealthInterval:  5 * time.Minute,
		MetricsEnabled:  true,
	}
}

// DBPoolStats holds database pool statistics
type DBPoolStats struct {
	MaxOpenConnections int           `json:"max_open_connections"`
	OpenConnections   int           `json:"open_connections"`
	InUse             int           `json:"in_use"`
	Idle              int           `json:"idle"`
	WaitCount         int64         `json:"wait_count"`
	WaitDuration      time.Duration `json:"wait_duration"`
	MaxIdleClosed     int64         `json:"max_idle_closed"`
	MaxIdleTimeClosed int64         `json:"max_idle_time_closed"`
	MaxLifetimeClosed int64         `json:"max_lifetime_closed"`
	HealthStatus      string        `json:"health_status"`
	LastHealthCheck   time.Time     `json:"last_health_check"`
}

type PostgreSQL struct {
	db           *gorm.DB
	cfg          *config.Config
	poolConfig   DBPoolConfig
	stats        DBPoolStats
	statsMutex   sync.RWMutex
	healthTicker *time.Ticker
	stopHealth   chan struct{}
}

func NewPostgreSQL(cfg *config.Config) (*PostgreSQL, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBSSLMode,
	)

	// Configure GORM logger
	var gormLogger logger.Interface
	if cfg.AppEnv == "production" {
		gormLogger = logger.Default.LogMode(logger.Silent)
	} else {
		gormLogger = logger.Default.LogMode(logger.Info)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get underlying SQL DB
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	// Set connection pool settings
	// MaxIdleConns: number of idle connections in the pool
	// MaxOpenConns: maximum number of open connections to the database
	// ConnMaxLifetime: maximum amount of time a connection may be reused
	if cfg.AppEnv == "production" {
		sqlDB.SetMaxIdleConns(25)                  // 25% of max open for production
		sqlDB.SetMaxOpenConns(100)                 // Higher limit for production load
		sqlDB.SetConnMaxLifetime(30 * time.Minute) // Rotate connections every 30 min
	} else {
		sqlDB.SetMaxIdleConns(10)           // Lower for development
		sqlDB.SetMaxOpenConns(25)           // Sufficient for dev/test
		sqlDB.SetConnMaxLifetime(time.Hour) // Longer lifetime in dev
	}
	sqlDB.SetConnMaxIdleTime(15 * time.Minute) // Close idle connections after 15 min

	// Test connection
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Successfully connected to PostgreSQL database")

	return &PostgreSQL{db: db, cfg: cfg}, nil
}

func (p *PostgreSQL) DB() *gorm.DB {
	return p.db
}

func (p *PostgreSQL) Close() error {
	// Stop health monitoring
	if p.healthTicker != nil {
		p.healthTicker.Stop()
		close(p.stopHealth)
	}

	sqlDB, err := p.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// AutoMigrate runs auto migration for given models
func (p *PostgreSQL) AutoMigrate(models ...interface{}) error {
	return p.db.AutoMigrate(models...)
}

// HealthCheck pings the database to verify connection health
func (p *PostgreSQL) HealthCheck(ctx context.Context) error {
	sqlDB, err := p.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}
	return sqlDB.PingContext(ctx)
}

// Stats returns database connection pool statistics
func (p *PostgreSQL) Stats() map[string]interface{} {
	sqlDB, err := p.db.DB()
	if err != nil {
		return map[string]interface{}{"error": err.Error()}
	}

	stats := sqlDB.Stats()

	p.statsMutex.RLock()
	defer p.statsMutex.RUnlock()

	return map[string]interface{}{
		"max_open_connections": stats.MaxOpenConnections,
		"open_connections":     stats.OpenConnections,
		"in_use":               stats.InUse,
		"idle":                 stats.Idle,
		"wait_count":           stats.WaitCount,
		"wait_duration":        stats.WaitDuration.String(),
		"max_idle_closed":      stats.MaxIdleClosed,
		"max_idle_time_closed": stats.MaxIdleTimeClosed,
		"max_lifetime_closed":  stats.MaxLifetimeClosed,
		"health_status":        p.stats.HealthStatus,
		"last_health_check":     p.stats.LastHealthCheck.Format(time.RFC3339),
	}
}

// GetDetailedStats returns detailed database pool statistics
func (p *PostgreSQL) GetDetailedStats() DBPoolStats {
	p.statsMutex.RLock()
	defer p.statsMutex.RUnlock()

	return p.stats
}

// UpdatePoolConfig updates the connection pool configuration
func (p *PostgreSQL) UpdatePoolConfig(config DBPoolConfig) error {
	sqlDB, err := p.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}

	sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(config.ConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(config.ConnMaxIdleTime)

	p.poolConfig = config
	return nil
}

// GetPoolConfig returns the current pool configuration
func (p *PostgreSQL) GetPoolConfig() DBPoolConfig {
	return p.poolConfig
}

// startHealthMonitoring starts periodic health checks
func (p *PostgreSQL) startHealthMonitoring() {
	p.healthTicker = time.NewTicker(p.poolConfig.HealthInterval)

	go func() {
		for {
			select {
			case <-p.healthTicker.C:
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				if err := p.HealthCheck(ctx); err != nil {
					p.statsMutex.Lock()
					p.stats.HealthStatus = "unhealthy"
					p.statsMutex.Unlock()
					log.Printf("Database health check failed: %v", err)
				} else {
					p.statsMutex.Lock()
					p.stats.HealthStatus = "healthy"
					p.stats.LastHealthCheck = time.Now()
					p.statsMutex.Unlock()
				}
				cancel()

				// Update stats
				p.updateStats()
			case <-p.stopHealth:
				return
			}
		}
	}()
}

// updateStats updates database pool statistics
func (p *PostgreSQL) updateStats() {
	sqlDB, err := p.db.DB()
	if err != nil {
		return
	}

	stats := sqlDB.Stats()

	p.statsMutex.Lock()
	p.stats.MaxOpenConnections = stats.MaxOpenConnections
	p.stats.OpenConnections = stats.OpenConnections
	p.stats.InUse = stats.InUse
	p.stats.Idle = stats.Idle
	p.stats.WaitCount = stats.WaitCount
	p.stats.WaitDuration = stats.WaitDuration
	p.stats.MaxIdleClosed = stats.MaxIdleClosed
	p.stats.MaxIdleTimeClosed = stats.MaxIdleTimeClosed
	p.stats.MaxLifetimeClosed = stats.MaxLifetimeClosed
	p.statsMutex.Unlock()
}

