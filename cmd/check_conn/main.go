package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, creating test using default/env vars")
	}

	// --- Check Postgres ---
	fmt.Println("-------------------------------------------")
	fmt.Println("Checking PostgreSQL Connectivity...")
	
	// Override host for local check
	dbHost := "localhost"
	dbPort := "5432"
	dbUser := os.Getenv("DB_USER")
	if dbUser == "" { dbUser = "karima_store" } // fallback based on .env view
	dbPass := os.Getenv("DB_PASSWORD")
	if dbPass == "" { dbPass = "lokal" }
	dbName := os.Getenv("DB_NAME")
	if dbName == "" { dbName = "karima_db" }

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		dbHost, dbUser, dbPass, dbName, dbPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Printf("❌ Failed to connect to Database: %v\n", err)
	} else {
		sqlDB, err := db.DB()
		if err != nil {
			fmt.Printf("❌ Failed to get generic database object: %v\n", err)
		} else {
			if err := sqlDB.Ping(); err != nil {
				fmt.Printf("❌ Failed to ping Database: %v\n", err)
			} else {
				fmt.Println("✅ Database Connection: OK")
				// Optional: Check a table
				var tableCount int64
				db.Raw("SELECT count(*) FROM information_schema.tables WHERE table_schema = 'public'").Scan(&tableCount)
				fmt.Printf("   Found %d tables in public schema.\n", tableCount)
			}
		}
	}

	// --- Check Redis ---
	fmt.Println("-------------------------------------------")
	fmt.Println("Checking Redis Connectivity...")

	// Override host for local check
	redisHost := "localhost"
	redisPort := "6380" // Mapped port
	redisPass := os.Getenv("REDIS_PASSWORD")

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort),
		Password: redisPass, // no password set
		DB:       0,  // use default DB
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		fmt.Printf("❌ Failed to connect to Redis: %v\n", err)
	} else {
		fmt.Printf("✅ Redis Connection: OK (Response: %s)\n", pong)
	}
	fmt.Println("-------------------------------------------")
}
