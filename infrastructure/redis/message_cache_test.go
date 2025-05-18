package redis

import (
	"context"
	"testing"
	"time"

	"clean-architecture-gochat/internal/domain/entities"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func setupTestRedis(t *testing.T) (*redis.Client, func()) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub redis connection", err)
	}

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	return client, func() {
		client.Close()
		mr.Close()
	}
}

func TestMessageCache_StoreAndGetPrivateMessages(t *testing.T) {
	client, cleanup := setupTestRedis(t)
	defer cleanup()

	cache := NewMessageCacheRepository(client)
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

	// 測試儲存訊息
	err := cache.StorePrivateMessage(ctx, message)
	assert.NoError(t, err)

	// 測試獲取私人訊息
	messages, err := cache.GetPrivateMessages(ctx, 1, 2)
	assert.NoError(t, err)
	if assert.NotNil(t, messages) && assert.Len(t, messages, 1) {
		assert.Equal(t, message.Content, messages[0].Content)
	}
}

func TestMessageCache_GetGroupMessages(t *testing.T) {
	client, cleanup := setupTestRedis(t)
	defer cleanup()

	cache := NewMessageCacheRepository(client)
	ctx := context.Background()

	// 準備測試數據
	messages := []*entities.Message{
		{
			ID:        1,
			UserId:    1,
			RoomID:    1,
			Content:   "Message 1",
			Type:      entities.MessageTypeGroup,
			Media:     entities.MediaTypeText,
			CreatedAt: time.Now().Add(-2 * time.Hour),
			UpdatedAt: time.Now().Add(-2 * time.Hour),
		},
		{
			ID:        2,
			UserId:    2,
			RoomID:    1,
			Content:   "Message 2",
			Type:      entities.MessageTypeGroup,
			Media:     entities.MediaTypeText,
			CreatedAt: time.Now().Add(-1 * time.Hour),
			UpdatedAt: time.Now().Add(-1 * time.Hour),
		},
	}

	// 儲存測試訊息
	for _, msg := range messages {
		err := cache.StoreGroupMessage(ctx, msg)
		assert.NoError(t, err)
	}

	// 測試獲取群組訊息
	history, err := cache.GetGroupMessages(ctx, 1)
	assert.NoError(t, err)
	if assert.NotNil(t, history) && assert.Len(t, history, 2) {
		assert.Equal(t, messages[0].Content, history[0].Content)
		assert.Equal(t, messages[1].Content, history[1].Content)
	}
}

func TestMessageCache_GetUserMessageList(t *testing.T) {
	client, cleanup := setupTestRedis(t)
	defer cleanup()

	cache := NewMessageCacheRepository(client)
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

	// 儲存訊息
	err := cache.StorePrivateMessage(ctx, message)
	assert.NoError(t, err)

	// 測試獲取用戶訊息列表
	messages, err := cache.GetUserMessageList(ctx, 1)
	assert.NoError(t, err)
	if assert.NotNil(t, messages) && assert.Len(t, messages, 1) {
		assert.Equal(t, message.Content, messages[0].Content)
	}
}
