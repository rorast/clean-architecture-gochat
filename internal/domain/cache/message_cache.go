package cache

import (
	"clean-architecture-gochat/internal/domain/entities"
	"context"
	"time"
)

// MessageCacheRepository 定義訊息快取儲存庫的介面
type MessageCacheRepository interface {
	// StorePrivateMessage 儲存私人訊息到快取
	StorePrivateMessage(ctx context.Context, message *entities.Message) error

	// StoreGroupMessage 儲存群組訊息到快取
	StoreGroupMessage(ctx context.Context, message *entities.Message) error

	// GetPrivateMessages 獲取私人訊息
	GetPrivateMessages(ctx context.Context, fromUserID, toUserID uint) ([]*entities.Message, error)

	// GetGroupMessages 獲取群組訊息
	GetGroupMessages(ctx context.Context, roomID uint) ([]*entities.Message, error)

	// GetUserMessageList 獲取用戶的訊息列表
	GetUserMessageList(ctx context.Context, userID uint) ([]*entities.Message, error)

	// CleanExpiredMessages 清理過期的訊息
	CleanExpiredMessages(ctx context.Context) error
}

// MessageCache 定義訊息快取的介面
type MessageCache interface {
	// StoreMessage 儲存訊息到快取
	StoreMessage(ctx context.Context, userID string, message *Message) error

	// GetUnreadMessages 獲取用戶的未讀訊息
	GetUnreadMessages(ctx context.Context, userID string) ([]*Message, error)

	// GetHistoryMessages 獲取歷史訊息
	GetHistoryMessages(ctx context.Context, userID string, limit, offset int) ([]*Message, error)

	// MarkMessageAsRead 標記訊息為已讀
	MarkMessageAsRead(ctx context.Context, userID string, messageID string) error
}

// Message 快取中的訊息結構
type Message struct {
	ID        string    `json:"id"`
	FromID    string    `json:"from_id"`
	ToID      string    `json:"to_id"`
	Content   string    `json:"content"`
	Type      string    `json:"type"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}
