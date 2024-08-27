package redis

import (
	"Magaz/backend/internal/config"
	"context"
	"github.com/redis/go-redis/v9"
)

// Ctx is a context background for Redis operations.
var Ctx = context.Background()

// TODO: use zap logger insted of fmt.Println
// InitRedisClient initializes a Redis client and returns it.
func InitRedisClient(cfg *config.RedisConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Host + ":" + cfg.Port,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// Perform a Ping to check if the connection is successful
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return client, nil
}
