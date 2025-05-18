package redis

import (
	"clean-architecture-gochat/internal/config"
	"context"
	"log"
	"sync"

	"github.com/redis/go-redis/v9"
)

var (
	client *redis.Client
	once   sync.Once
)

func Connect() *redis.Client {
	once.Do(func() {
		client = redis.NewClient(&redis.Options{
			Addr:         config.Config.Redis.Addr,
			Password:     config.Config.Redis.Password,
			DB:           config.Config.Redis.DB,
			PoolSize:     config.Config.Redis.PoolSize,
			MinIdleConns: config.Config.Redis.MinIdleConns,
		})

		if err := client.Ping(context.Background()).Err(); err != nil {
			log.Fatalf("Redis 連接失敗: %v", err)
		}
		log.Println("Redis 連接成功")
	})

	return client
}

func GetRedis() *redis.Client {
	if client == nil {
		log.Fatal("Redis 尚未初始化")
	}
	return client
}
