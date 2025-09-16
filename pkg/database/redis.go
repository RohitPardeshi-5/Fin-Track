package database

import (
	"context"
	"github.com/redis/go-redis/v9"
)

func NewRedisClient(redisURL string) (*redis.Client, error) {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(opt)
	
	// Test connection
	_, err = client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return client, nil
}