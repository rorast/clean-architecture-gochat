package redis

import (
	"clean-architecture-gochat/internal/domain/cache"
	"clean-architecture-gochat/internal/domain/entities"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	// 訊息快取的 key 格式
	messageKeyFormat  = "chat:msg:%d:%d"  // chat:msg:fromUserId:toUserId
	roomKeyFormat     = "chat:room:%d"    // chat:room:roomId
	userMsgListFormat = "chat:msglist:%d" // chat:msglist:userId
	messagePattern    = "chat:msg:*:*"    // 用於搜尋所有私人訊息
	roomPattern       = "chat:room:*"     // 用於搜尋所有群組訊息
	userListPattern   = "chat:msglist:*"  // 用於搜尋所有用戶列表

	// 快取過期時間
	messageTTL     = 24 * time.Hour // 單條訊息快取 24 小時
	roomTTL        = 48 * time.Hour // 群組訊息快取 48 小時
	messageListTTL = 72 * time.Hour // 訊息列表快取 72 小時

	// 訊息列表的最大長度
	maxMessageListLength = 100
)

type redisCacheRepository struct {
	client *redis.Client
}

// NewMessageCacheRepository 創建新的訊息快取儲存庫
func NewMessageCacheRepository(client *redis.Client) cache.MessageCacheRepository {
	return &redisCacheRepository{
		client: client,
	}
}

func (r *redisCacheRepository) StorePrivateMessage(ctx context.Context, message *entities.Message) error {
	key := fmt.Sprintf(messageKeyFormat, message.UserId, message.TargetId)

	// 序列化消息
	data, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("序列化訊息失敗: %v", err)
	}

	// 使用 LPUSH 將消息添加到列表頭部
	if err := r.client.LPush(ctx, key, data).Err(); err != nil {
		return fmt.Errorf("儲存訊息到 Redis 失敗: %v", err)
	}

	// 設置過期時間
	if err := r.client.Expire(ctx, key, messageTTL).Err(); err != nil {
		return fmt.Errorf("設置訊息過期時間失敗: %v", err)
	}

	// 修剪列表以保持最大長度
	if err := r.client.LTrim(ctx, key, 0, maxMessageListLength-1).Err(); err != nil {
		return fmt.Errorf("修剪訊息列表失敗: %v", err)
	}

	return nil
}

func (r *redisCacheRepository) StoreGroupMessage(ctx context.Context, message *entities.Message) error {
	key := fmt.Sprintf(roomKeyFormat, message.RoomID)

	// 序列化消息
	data, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("序列化群組訊息失敗: %v", err)
	}

	// 使用 LPUSH 將消息添加到列表頭部
	if err := r.client.LPush(ctx, key, data).Err(); err != nil {
		return fmt.Errorf("儲存群組訊息到 Redis 失敗: %v", err)
	}

	// 設置過期時間
	if err := r.client.Expire(ctx, key, roomTTL).Err(); err != nil {
		return fmt.Errorf("設置群組訊息過期時間失敗: %v", err)
	}

	// 修剪列表以保持最大長度
	if err := r.client.LTrim(ctx, key, 0, maxMessageListLength-1).Err(); err != nil {
		return fmt.Errorf("修剪群組訊息列表失敗: %v", err)
	}

	return nil
}

func (r *redisCacheRepository) GetPrivateMessages(ctx context.Context, fromUserID, toUserID uint) ([]*entities.Message, error) {
	key := fmt.Sprintf(messageKeyFormat, fromUserID, toUserID)

	// 獲取所有消息
	data, err := r.client.LRange(ctx, key, 0, -1).Result()
	if err == redis.Nil {
		return []*entities.Message{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("從 Redis 獲取私人訊息失敗: %v", err)
	}

	messages := make([]*entities.Message, 0, len(data))
	for _, msgData := range data {
		var msg entities.Message
		if err := json.Unmarshal([]byte(msgData), &msg); err != nil {
			continue // 跳過無法解析的消息
		}
		messages = append(messages, &msg)
	}

	return messages, nil
}

func (r *redisCacheRepository) GetGroupMessages(ctx context.Context, roomID uint) ([]*entities.Message, error) {
	key := fmt.Sprintf(roomKeyFormat, roomID)

	// 獲取所有消息
	data, err := r.client.LRange(ctx, key, 0, -1).Result()
	if err == redis.Nil {
		return []*entities.Message{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("從 Redis 獲取群組訊息失敗: %v", err)
	}

	messages := make([]*entities.Message, 0, len(data))
	for _, msgData := range data {
		var msg entities.Message
		if err := json.Unmarshal([]byte(msgData), &msg); err != nil {
			continue // 跳過無法解析的消息
		}
		messages = append(messages, &msg)
	}

	return messages, nil
}

func (r *redisCacheRepository) GetUserMessageList(ctx context.Context, userID uint) ([]*entities.Message, error) {
	pattern := fmt.Sprintf("chat:msg:%d:*", userID)
	keys, err := r.client.Keys(ctx, pattern).Result()
	if err != nil {
		return nil, fmt.Errorf("獲取用戶訊息列表失敗: %v", err)
	}

	var allMessages []*entities.Message
	for _, key := range keys {
		data, err := r.client.LRange(ctx, key, 0, -1).Result()
		if err != nil {
			continue
		}

		for _, msgData := range data {
			var msg entities.Message
			if err := json.Unmarshal([]byte(msgData), &msg); err != nil {
				continue
			}
			allMessages = append(allMessages, &msg)
		}
	}

	return allMessages, nil
}

func (r *redisCacheRepository) CleanExpiredMessages(ctx context.Context) error {
	// 清理私人訊息
	if err := r.cleanExpiredKeys(ctx, messagePattern); err != nil {
		return err
	}

	// 清理群組訊息
	if err := r.cleanExpiredKeys(ctx, roomPattern); err != nil {
		return err
	}

	return nil
}

func (r *redisCacheRepository) cleanExpiredKeys(ctx context.Context, pattern string) error {
	keys, err := r.client.Keys(ctx, pattern).Result()
	if err != nil {
		return fmt.Errorf("獲取過期訊息失敗: %v", err)
	}

	for _, key := range keys {
		if err := r.client.Del(ctx, key).Err(); err != nil {
			return fmt.Errorf("刪除過期訊息失敗: %v", err)
		}
	}

	return nil
}
