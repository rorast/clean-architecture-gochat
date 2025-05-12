package redis

import (
	"clean-architecture-gochat/internal/config"
	"github.com/go-redis/redis/v8"
	"log"
	"sync"
)

var (
	client *redis.Client
	once   sync.Once
)

func Connect() *redis.Client {
	once.Do(func() {
		client := redis.NewClient(&redis.Options{
			Addr:         config.Config.Redis.Addr,
			Password:     config.Config.Redis.Password,
			DB:           config.Config.Redis.DB,
			PoolSize:     config.Config.Redis.PoolSize,
			MinIdleConns: config.Config.Redis.MinIdleConns,
		})

		if err := client.Ping(client.Context()).Err(); err != nil {
			log.Fatalf("Redis connection failed: %v", err)
		}
		log.Println("Redis connected successfully.")
	})

	return client
}

func GetRedis() *redis.Client {
	if client == nil {
		log.Fatal("Redis not initialized")
	}
	return client
}
