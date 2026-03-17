package config

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	redisClient *redis.Client
	_once       sync.Once
)

func LoadRedisConnection() *redis.Client {
	_once.Do(func() {
		addr := fmt.Sprintf("%s:%s", GetEnv("REDIS_URL"), GetEnv("REDIS_PORT"))

		client := redis.NewClient(&redis.Options{
			Addr:         addr,
			Password:     GetEnv("REDIS_PASSWORD"),
			DB:           0,
			PoolSize:     20,
			MinIdleConns: 5,
			DialTimeout:  5 * time.Second,
			ReadTimeout:  3 * time.Second,
			WriteTimeout: 3 * time.Second,
		})

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if _, err := client.Ping(ctx).Result(); err != nil {
			panic(fmt.Sprintf("Redis connection failed: %v", err))
		}

		redisClient = client
		log.Printf("Redis connected: %s", addr)
	})
	return redisClient
}

func GetRedisClient() *redis.Client {
	return redisClient
}

func CloseRedisConnection() {
	if redisClient == nil {
		return
	}
	if err := redisClient.Close(); err != nil {
		log.Println("Redis close error:", err)
	}
}
