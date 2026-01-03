package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/karima-store/internal/config"
	"github.com/karima-store/internal/database"
	"github.com/karima-store/internal/logger"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	// Parse command line flags
	action := flag.String("action", "up", "Migration action: up, down, steps, version, force")
	steps := flag.Int("steps", 1, "Number of migration steps (for 'steps' action)")
	version := flag.Int("version", 0, "Target version (for 'force' action)")
	migrationsPath := flag.String("path", "migrations", "Path to migrations directory")
	flag.Parse()

	// Load configuration
	cfg := config.Load()

	// Initialize structured logger
	if err := cfg.InitLogger(); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	// Initialize database connection
	db, err := database.NewPostgreSQL(cfg)
	if err != nil {
		if logger.Log != nil {
			logger.Log.Fatalw("Failed to connect to database", "error", err)
		} else {
			log.Fatalf("Failed to connect to database: %v", err)
		}
	}
	defer db.Close()

	// Execute the requested action
	switch *action {
	case "up":
		if err := db.Migrate(*migrationsPath); err != nil {
			if logger.Log != nil {
				logger.Log.Fatalw("Migration failed", "error", err)
			} else {
				log.Fatalf("Migration failed: %v", err)
			}
		}
		fmt.Println("✅ Migrations applied successfully")

	case "down":
		if err := db.MigrateDown(*migrationsPath); err != nil {
			if logger.Log != nil {
				logger.Log.Fatalw("Migration rollback failed", "error", err)
			} else {
				log.Fatalf("Migration rollback failed: %v", err)
			}
		}
		fmt.Println("✅ Migration rolled back successfully")

	case "steps":
		if err := db.MigrateSteps(*migrationsPath, *steps); err != nil {
			if logger.Log != nil {
				logger.Log.Fatalw("Migration steps failed", "error", err, "steps", *steps)
			} else {
				log.Fatalf("Migration steps failed: %v", err)
			}
		}
		fmt.Printf("✅ %d migration step(s) completed successfully\n", *steps)

	case "version":
		version, dirty, err := db.GetMigrationVersion(*migrationsPath)
		if err != nil {
			if logger.Log != nil {
				logger.Log.Fatalw("Failed to get migration version", "error", err)
			} else {
				log.Fatalf("Failed to get migration version: %v", err)
			}
		}
		if dirty {
			fmt.Printf("⚠️  Current migration version: %d (DIRTY)\n", version)
			fmt.Println("The database is in a dirty state. Please fix it manually or use force.")
			os.Exit(1)
		} else {
			fmt.Printf("✅ Current migration version: %d\n", version)
		}

	case "force":
		if *version < 0 {
			if logger.Log != nil {
				logger.Log.Fatal("Version must be a non-negative integer")
			} else {
				log.Fatal("Version must be a non-negative integer")
			}
		}
		if err := db.ForceVersion(*migrationsPath, *version); err != nil {
			if logger.Log != nil {
				logger.Log.Fatalw("Failed to force migration version", "error", err, "version", *version)
			} else {
				log.Fatalf("Failed to force migration version: %v", err)
			}
		}
		fmt.Printf("✅ Migration version forced to %d\n", *version)

	default:
		fmt.Printf("Unknown action: %s\n", *action)
		fmt.Println("Available actions: up, down, steps, version, force")
		flag.Usage()
		os.Exit(1)
	}
}
