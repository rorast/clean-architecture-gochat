package redis

import (
	"clean-architecture-gochat/internal/common/enum"
	"clean-architecture-gochat/internal/config"
	appErrors "clean-architecture-gochat/internal/errors"
	"context"
	"log"
	"sync"

	"github.com/redis/go-redis/v9"
)

var (
	client *redis.Client
	once   sync.Once
)

// Connect 連接到Redis並返回客戶端
func Connect() (*redis.Client, error) {
	var err error

	once.Do(func() {
		client = redis.NewClient(&redis.Options{
			Addr:         config.Config.Redis.Addr,
			Password:     config.Config.Redis.Password,
			DB:           config.Config.Redis.DB,
			PoolSize:     config.Config.Redis.PoolSize,
			MinIdleConns: config.Config.Redis.MinIdleConns,
		})

		if pingErr := client.Ping(context.Background()).Err(); pingErr != nil {
			err = appErrors.Wrap(pingErr, enum.ErrRedisConnectionFailed, map[string]interface{}{
				"addr":     config.Config.Redis.Addr,
				"database": config.Config.Redis.DB,
			})
			log.Printf("Redis連接失敗: %v", err)
			return
		}
		log.Println("Redis連接成功")
	})

	return client, err
}

// GetRedis 獲取已初始化的Redis客戶端
func GetRedis() (*redis.Client, error) {
	if client == nil {
		return nil, appErrors.New(enum.ErrRedisConnectionFailed, "Redis尚未初始化")
	}
	return client, nil
}
