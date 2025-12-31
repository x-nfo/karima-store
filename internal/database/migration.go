package database

import (
	"fmt"
	"log"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4"
)

// Migrate runs database migrations using golang-migrate
func (p *PostgreSQL) Migrate(migrationsPath string) error {
	// Build the database URL for migrate
	// We need to construct a connection string that includes the database name
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		p.cfg.DBUser,
		p.cfg.DBPassword,
		p.cfg.DBHost,
		p.cfg.DBPort,
		p.cfg.DBName,
		p.cfg.DBSSLMode,
	)

	// Create a new migrate instance
	m, err := migrate.New(
		"file://"+migrationsPath,
		dsn,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	// Ensure the migrate instance is closed
	defer func() {
		srcErr, dbErr := m.Close()
		if srcErr != nil {
			log.Printf("Error closing migrate source: %v", srcErr)
		}
		if dbErr != nil {
			log.Printf("Error closing migrate database: %v", dbErr)
		}
	}()

	// Run migrations
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("Database migrations completed successfully")
	return nil
}

// MigrateDown rolls back the last migration
func (p *PostgreSQL) MigrateDown(migrationsPath string) error {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		p.cfg.DBUser,
		p.cfg.DBPassword,
		p.cfg.DBHost,
		p.cfg.DBPort,
		p.cfg.DBName,
		p.cfg.DBSSLMode,
	)

	m, err := migrate.New(
		"file://"+migrationsPath,
		dsn,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	defer func() {
		srcErr, dbErr := m.Close()
		if srcErr != nil {
			log.Printf("Error closing migrate source: %v", srcErr)
		}
		if dbErr != nil {
			log.Printf("Error closing migrate database: %v", dbErr)
		}
	}()

	if err := m.Down(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to rollback migration: %w", err)
	}

	log.Println("Database rollback completed successfully")
	return nil
}

// MigrateSteps runs or rolls back a specific number of migration steps
func (p *PostgreSQL) MigrateSteps(migrationsPath string, steps int) error {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		p.cfg.DBUser,
		p.cfg.DBPassword,
		p.cfg.DBHost,
		p.cfg.DBPort,
		p.cfg.DBName,
		p.cfg.DBSSLMode,
	)

	m, err := migrate.New(
		"file://"+migrationsPath,
		dsn,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	defer func() {
		srcErr, dbErr := m.Close()
		if srcErr != nil {
			log.Printf("Error closing migrate source: %v", srcErr)
		}
		if dbErr != nil {
			log.Printf("Error closing migrate database: %v", dbErr)
		}
	}()

	if err := m.Steps(steps); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migration steps: %w", err)
	}

	log.Printf("Database migration steps (%d) completed successfully", steps)
	return nil
}

// GetMigrationVersion returns the current migration version
func (p *PostgreSQL) GetMigrationVersion(migrationsPath string) (uint, bool, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		p.cfg.DBUser,
		p.cfg.DBPassword,
		p.cfg.DBHost,
		p.cfg.DBPort,
		p.cfg.DBName,
		p.cfg.DBSSLMode,
	)

	m, err := migrate.New(
		"file://"+migrationsPath,
		dsn,
	)
	if err != nil {
		return 0, false, fmt.Errorf("failed to create migrate instance: %w", err)
	}

	defer func() {
		srcErr, dbErr := m.Close()
		if srcErr != nil {
			log.Printf("Error closing migrate source: %v", srcErr)
		}
		if dbErr != nil {
			log.Printf("Error closing migrate database: %v", dbErr)
		}
	}()

	version, dirty, err := m.Version()
	if err != nil {
		if err == migrate.ErrNilVersion {
			return 0, false, nil
		}
		return 0, false, fmt.Errorf("failed to get migration version: %w", err)
	}

	return version, dirty, nil
}

// ForceVersion sets the migration version to a specific value (use with caution)
func (p *PostgreSQL) ForceVersion(migrationsPath string, version int) error {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		p.cfg.DBUser,
		p.cfg.DBPassword,
		p.cfg.DBHost,
		p.cfg.DBPort,
		p.cfg.DBName,
		p.cfg.DBSSLMode,
	)

	m, err := migrate.New(
		"file://"+migrationsPath,
		dsn,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	defer func() {
		srcErr, dbErr := m.Close()
		if srcErr != nil {
			log.Printf("Error closing migrate source: %v", srcErr)
		}
		if dbErr != nil {
			log.Printf("Error closing migrate database: %v", dbErr)
		}
	}()

	if err := m.Force(version); err != nil {
		return fmt.Errorf("failed to force migration version: %w", err)
	}

	log.Printf("Migration version forced to %d", version)
	return nil
}
