package repositories

import (
	"clean-architecture-gochat/internal/domain/entities"
	"context"
)

// MessageService 定義訊息服務的介面
type MessageService interface {
	// SendPrivateMessage 發送私人訊息
	SendPrivateMessage(ctx context.Context, message *entities.Message) error

	// SendGroupMessage 發送群組訊息
	SendGroupMessage(ctx context.Context, message *entities.Message) error

	// GetPrivateHistory 獲取私人聊天歷史
	GetPrivateHistory(ctx context.Context, fromUserID, toUserID uint) ([]*entities.Message, error)

	// GetGroupHistory 獲取群組聊天歷史
	GetGroupHistory(ctx context.Context, roomID uint) ([]*entities.Message, error)

	// GetRecentMessages 獲取用戶最近的訊息列表
	GetRecentMessages(ctx context.Context, userID uint) ([]*entities.Message, error)
}
