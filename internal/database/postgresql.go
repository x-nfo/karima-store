package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/karima-store/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type PostgreSQL struct {
	db  *gorm.DB
	cfg *config.Config
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
	}
}
