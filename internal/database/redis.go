package database

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/karima-store/internal/config"
	"github.com/redis/go-redis/v9"
)

type Redis struct {
	client *redis.Client
}

func NewRedis(cfg *config.Config) (*Redis, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
		Password: cfg.RedisPassword,
		DB:       0,
		// Connection pool settings
		PoolSize:        20,               // Maximum number of socket connections
		MinIdleConns:    5,                // Minimum number of idle connections
		MaxRetries:      3,                // Maximum number of retries before giving up
		DialTimeout:     5 * time.Second,  // Dial timeout
		ReadTimeout:     3 * time.Second,  // Read timeout
		WriteTimeout:    3 * time.Second,  // Write timeout
		PoolTimeout:     4 * time.Second,  // Time client waits for connection if all busy
		ConnMaxIdleTime: 15 * time.Minute, // Close idle connections after 15 min
		ConnMaxLifetime: 30 * time.Minute, // Rotate connections every 30 min
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Test connection
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	log.Println("Successfully connected to Redis")

	return &Redis{client: client}, nil
}

func (r *Redis) Client() *redis.Client {
	return r.client
}

func (r *Redis) Close() error {
	return r.client.Close()
}

// Get retrieves a value from Redis
func (r *Redis) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

// Set stores a value in Redis with an expiration time
func (r *Redis) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.client.Set(ctx, key, value, expiration).Err()
}

// GetJSON retrieves a JSON value from Redis and unmarshals it into dest
func (r *Redis) GetJSON(ctx context.Context, key string, dest interface{}) error {
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(val), dest)
}

// SetJSON marshals a value to JSON and stores it in Redis with an expiration time
func (r *Redis) SetJSON(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	bytes, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, key, bytes, expiration).Err()
}

// Delete removes a key from Redis
func (r *Redis) Delete(ctx context.Context, keys ...string) error {
	return r.client.Del(ctx, keys...).Err()
}

// Exists checks if a key exists in Redis
func (r *Redis) Exists(ctx context.Context, keys ...string) (int64, error) {
	return r.client.Exists(ctx, keys...).Result()
}

// FlushDB clears the current database
func (r *Redis) FlushDB(ctx context.Context) error {
	return r.client.FlushDB(ctx).Err()
}

// DeleteByPattern deletes all keys matching a pattern
func (r *Redis) DeleteByPattern(ctx context.Context, pattern string) error {
	var cursor uint64
	var keys []string
	var err error

	for {
		keys, cursor, err = r.client.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			return err
		}

		if len(keys) > 0 {
			if err := r.client.Del(ctx, keys...).Err(); err != nil {
				return err
			}
		}

		if cursor == 0 {
			break
		}
	}
	return nil
}

// HealthCheck pings Redis to verify connection health
func (r *Redis) HealthCheck(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}

// PoolStats returns Redis connection pool statistics
func (r *Redis) PoolStats() map[string]interface{} {
	stats := r.client.PoolStats()
	return map[string]interface{}{
		"hits":        stats.Hits,
		"misses":      stats.Misses,
		"timeouts":    stats.Timeouts,
		"total_conns": stats.TotalConns,
		"idle_conns":  stats.IdleConns,
		"stale_conns": stats.StaleConns,
	}
}
