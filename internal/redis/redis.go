package redis

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

var (
	// RedisClient is the global Redis client instance.
	RedisClient *redis.Client
	ctx         = context.Background()
)

// InitRedis initializes the Redis client from environment variables.
func InitRedis() {
	// Load redis.env file
	if err := godotenv.Load("configs/redis.env"); err != nil {
		log.Warn().Err(err).Msg("Could not load redis.env file, using default environment variables")
	}

	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")
	password := os.Getenv("REDIS_PASSWORD")
	dbStr := os.Getenv("REDIS_DB")

	if host == "" {
		host = "localhost"
	}
	if port == "" {
		port = "6379"
	}
	db, err := strconv.Atoi(dbStr)
	if err != nil {
		db = 0 // Default to DB 0 if not specified or invalid
	}

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: password,
		DB:       db,
	})

	// Ping the Redis server to ensure the connection is established
	if _, err := RedisClient.Ping(ctx).Result(); err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to Redis")
	}

	log.Info().Msg("Successfully connected to Redis.")
}