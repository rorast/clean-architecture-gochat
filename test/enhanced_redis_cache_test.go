package test

import (
	"context"
	"testing"
	"time"

	"clean-architecture-gochat/infrastructure/redis"
	"clean-architecture-gochat/internal/domain/entities"
	baseErrors "clean-architecture-gochat/pkg/errors"

	"github.com/alicebob/miniredis/v2"
	goredis "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func setupTestEnhancedRedis(t *testing.T) (*goredis.Client, func()) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub redis connection", err)
	}

	client := goredis.NewClient(&goredis.Options{
		Addr: mr.Addr(),
	})

	return client, func() {
		client.Close()
		mr.Close()
	}
}

func TestEnhancedMessageCache_Basic(t *testing.T) {
	client, cleanup := setupTestEnhancedRedis(t)
	defer cleanup()

	cache := redis.NewMessageCacheRepository(client)
	ctx := context.Background()

	// 準備測試數據
	message := &entities.Message{
		ID:        1,
		UserId:    1,
		TargetId:  2,
		Content:   "Hello",
		Type:      entities.MessageTypePrivate,
		Media:     entities.MediaTypeText,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 測試儲存私人訊息
	err := cache.StorePrivateMessage(ctx, message)
	assert.NoError(t, err)

	// 測試讀取私人訊息
	messages, err := cache.GetPrivateMessages(ctx, 1, 2)
	assert.NoError(t, err)
	assert.Len(t, messages, 1)
	assert.Equal(t, "Hello", messages[0].Content)
}

func TestEnhancedMessageCache_GroupMessages(t *testing.T) {
	client, cleanup := setupTestEnhancedRedis(t)
	defer cleanup()

	cache := redis.NewMessageCacheRepository(client)
	ctx := context.Background()

	// 準備測試數據
	groupMessage := &entities.Message{
		ID:        2,
		UserId:    1,
		RoomID:    101,
		Content:   "Group message",
		Type:      entities.MessageTypeGroup,
		Media:     entities.MediaTypeText,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 測試儲存群組訊息
	err := cache.StoreGroupMessage(ctx, groupMessage)
	assert.NoError(t, err)

	// 測試讀取群組訊息
	messages, err := cache.GetGroupMessages(ctx, 101)
	assert.NoError(t, err)
	assert.Len(t, messages, 1)
	assert.Equal(t, "Group message", messages[0].Content)
}

func TestEnhancedMessageCache_ErrorHandling(t *testing.T) {
	ctx := context.Background()

	// 準備測試數據，使用可以觸發錯誤的數據
	// 這裡我們使用通道關閉的方式強制產生錯誤
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("無法啟動 miniredis：%v", err)
	}

	// 獲取正確連接並立即關閉連接
	brokenClient := goredis.NewClient(&goredis.Options{
		Addr: mr.Addr(),
	})
	mr.Close() // 立即關閉 miniredis 服務器，導致後續操作失敗

	brokenCache := redis.NewMessageCacheRepository(brokenClient)

	message := &entities.Message{
		ID:        1,
		UserId:    1,
		TargetId:  2,
		Content:   "Test Message",
		Type:      entities.MessageTypePrivate,
		Media:     entities.MediaTypeText,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 測試儲存消息，應該會因為連接關閉而失敗
	err = brokenCache.StorePrivateMessage(ctx, message)

	// 檢查錯誤是否發生
	if assert.Error(t, err, "應該發生錯誤") {
		// 檢查返回的錯誤是否為我們的自定義錯誤類型
		if assert.True(t, baseErrors.IsAppError(err), "應該是應用錯誤類型") {
			// 獲取應用錯誤並檢查其屬性
			appErr, ok := baseErrors.GetAppError(err)
			if assert.True(t, ok, "應該能獲取到應用錯誤") {
				// 檢查是否為 Redis 操作錯誤
				errCode := appErr.Code()
				assert.Equal(t, 8001, errCode, "應該是 Redis 操作錯誤") // ErrRedisOperationFailed 的值
				assert.Contains(t, appErr.Message(), "Redis操作失敗", "錯誤信息應包含 'Redis操作失敗'")
			}
		}
	}
}

// 測試連接錯誤處理
func TestEnhancedMessageCache_ConnectionError(t *testing.T) {
	// 建立錯誤的連接
	client := goredis.NewClient(&goredis.Options{
		Addr: "nonexistent:6379", // 不存在的地址
	})

	cache := redis.NewMessageCacheRepository(client)
	ctx := context.Background()

	message := &entities.Message{
		ID:        1,
		UserId:    1,
		TargetId:  2,
		Content:   "Hello",
		Type:      entities.MessageTypePrivate,
		Media:     entities.MediaTypeText,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 應該返回連接錯誤
	err := cache.StorePrivateMessage(ctx, message)
	assert.Error(t, err)

	// 檢查錯誤類型
	appErr, ok := baseErrors.GetAppError(err)
	assert.True(t, ok)

	// 檢查是否為 Redis 操作錯誤
	assert.Equal(t, 8001, appErr.Code()) // ErrRedisOperationFailed 的值
}
