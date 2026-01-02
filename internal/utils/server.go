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
)

// startServerWithGracefulShutdown handles graceful shutdown of the server
// It captures SIGINT and SIGTERM signals, and properly closes all resources
func startServerWithGracefulShutdown(
	app *fiber.App,
	port string,
	cfg *config.Config,
	db *database.PostgreSQL,
	redis *database.Redis,
) {
	// Create a channel to listen for interrupt signals
	quit := make(chan os.Signal, 1)

	// Register SIGINT (Ctrl+C) and SIGTERM signals
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start the server in a goroutine
	go func() {
		log.Printf("Server is running on port %s", port)
		if err := app.Listen(fmt.Sprintf(":%s", port)); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Block until a signal is received
	sig := <-quit
	log.Printf("Received signal: %v", sig)
	log.Println("Shutting down server gracefully...")

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
		log.Printf("Error during Fiber shutdown: %v", err)
	} else {
		log.Println("Fiber server shut down successfully")
	}

	// Close Redis connection
	log.Println("Closing Redis connection...")
	if err := redis.Close(); err != nil {
		log.Printf("Error closing Redis connection: %v", err)
	} else {
		log.Println("Redis connection closed successfully")
	}

	// Close PostgreSQL connection
	log.Println("Closing PostgreSQL connection...")
	if err := db.Close(); err != nil {
		log.Printf("Error closing PostgreSQL connection: %v", err)
	} else {
		log.Println("PostgreSQL connection closed successfully")
	}

	log.Println("Graceful shutdown completed")
}

// StartServerWithGracefulShutdown is the exported version of startServerWithGracefulShutdown
// This allows it to be used from other packages if needed
func StartServerWithGracefulShutdown(
	app *fiber.App,
	port string,
	cfg *config.Config,
	db *database.PostgreSQL,
	redis *database.Redis,
) {
	startServerWithGracefulShutdown(app, port, cfg, db, redis)
}
