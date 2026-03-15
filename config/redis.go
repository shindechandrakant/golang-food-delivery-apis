package config

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/redis/go-redis/v9"
)

var (
	ctx         = context.Background()
	redisClient *redis.Client
	_once       sync.Once
)

func LoadRedisConnection() *redis.Client {

	_once.Do(func() {
		redisURI := GetEnv("REDIS_URL")
		redisPORT := GetEnv("REDIS_PORT")
		RedisClient := redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", redisURI, redisPORT),
			Password: "", // no password set
			DB:       0,  // use default DB
		})
		defer RedisClient.Close()

		_, err := RedisClient.Ping(ctx).Result()
		if err != nil {
			panic(err)
		} else {
			fmt.Println("Redis Connected successfully..")
		}
	})
	return redisClient
}

func CloseRedisConnection() {
	if err := redisClient.Close(); err != nil {
		log.Println("Redis close error:", err)
	}
}
