package utils

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/karima-store/internal/config"
	"github.com/karima-store/internal/database"
	"github.com/karima-store/internal/logger"
)

// startServerWithGracefulShutdown handles graceful shutdown of the server
// It captures SIGINT and SIGTERM signals, and properly closes all resources
func startServerWithGracefulShutdown(
	app *fiber.App,
	port string,
	cfg *config.Config,
	db *database.PostgreSQL,
	redis database.RedisClient,
) {
	// Create a channel to listen for interrupt signals
	quit := make(chan os.Signal, 1)

	// Register SIGINT (Ctrl+C) and SIGTERM signals
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start the server in a goroutine
	go func() {
		if logger.Log != nil {
			logger.Log.Infow("Server is running", "port", port)
		} else {
			log.Printf("Server is running on port %s", port)
		}
		if err := app.Listen(fmt.Sprintf(":%s", port)); err != nil {
			if logger.Log != nil {
				logger.Log.Fatalw("Failed to start server", "error", err)
			} else {
				log.Fatalf("Failed to start server: %v", err)
			}
		}
	}()

	// Block until a signal is received
	sig := <-quit
	if logger.Log != nil {
		logger.Log.Infow("Received shutdown signal", "signal", sig)
		logger.Log.Info("Shutting down server gracefully...")
	} else {
		log.Printf("Received signal: %v", sig)
		log.Println("Shutting down server gracefully...")
	}

	// Create a context with timeout for shutdown
	shutdownTimeout := 30 * time.Second
	if cfg.AppEnv == "production" {
		shutdownTimeout = 60 * time.Second // Longer timeout for production
	}
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	// Shutdown Fiber app gracefully
	// This will stop accepting new connections and wait for ongoing requests to complete
	if err := app.ShutdownWithContext(ctx); err != nil {
		if logger.Log != nil {
			logger.Log.Errorw("Error during Fiber shutdown", "error", err)
		} else {
			log.Printf("Error during Fiber shutdown: %v", err)
		}
	} else {
		if logger.Log != nil {
			logger.Log.Info("Fiber server shut down successfully")
		} else {
			log.Println("Fiber server shut down successfully")
		}
	}

	// Close Redis connection
	if logger.Log != nil {
		logger.Log.Info("Closing Redis connection...")
	} else {
		log.Println("Closing Redis connection...")
	}
	if err := redis.Close(); err != nil {
		if logger.Log != nil {
			logger.Log.Errorw("Error closing Redis connection", "error", err)
		} else {
			log.Printf("Error closing Redis connection: %v", err)
		}
	} else {
		if logger.Log != nil {
			logger.Log.Info("Redis connection closed successfully")
		} else {
			log.Println("Redis connection closed successfully")
		}
	}

	// Close PostgreSQL connection
	if logger.Log != nil {
		logger.Log.Info("Closing PostgreSQL connection...")
	} else {
		log.Println("Closing PostgreSQL connection...")
	}
	if err := db.Close(); err != nil {
		if logger.Log != nil {
			logger.Log.Errorw("Error closing PostgreSQL connection", "error", err)
		} else {
			log.Printf("Error closing PostgreSQL connection: %v", err)
		}
	} else {
		if logger.Log != nil {
			logger.Log.Info("PostgreSQL connection closed successfully")
		} else {
			log.Println("PostgreSQL connection closed successfully")
		}
	}

	if logger.Log != nil {
		logger.Log.Info("Graceful shutdown completed")
	} else {
		log.Println("Graceful shutdown completed")
	}
}

// StartServerWithGracefulShutdown is the exported version of startServerWithGracefulShutdown
// This allows it to be used from other packages if needed
func StartServerWithGracefulShutdown(
	app *fiber.App,
	port string,
	cfg *config.Config,
	db *database.PostgreSQL,
	redis database.RedisClient,
) {
	startServerWithGracefulShutdown(app, port, cfg, db, redis)
}
