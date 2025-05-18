package message

import (
	"clean-architecture-gochat/internal/common/enum"
	"clean-architecture-gochat/internal/domain/cache"
	"clean-architecture-gochat/internal/domain/entities"
	"clean-architecture-gochat/internal/domain/repositories"
	appErrors "clean-architecture-gochat/internal/errors"
	"context"
	"time"
)

// EnhancedMessageUseCase 加強版消息用例，整合新的錯誤處理系統
type EnhancedMessageUseCase struct {
	messageRepo  repositories.MessageRepository
	messageCache cache.MessageCacheRepository
}

// NewEnhancedMessageUseCase 創建新的加強版消息用例
func NewEnhancedMessageUseCase(
	messageRepo repositories.MessageRepository,
	messageCache cache.MessageCacheRepository,
) *EnhancedMessageUseCase {
	return &EnhancedMessageUseCase{
		messageRepo:  messageRepo,
		messageCache: messageCache,
	}
}

// SendPrivateMessage 發送私人消息，使用新的錯誤處理系統
func (uc *EnhancedMessageUseCase) SendPrivateMessage(ctx context.Context, message *entities.Message) error {
	if message == nil {
		return appErrors.New(enum.ErrInvalidInput, "消息不能為空")
	}

	if message.UserId == 0 {
		return appErrors.New(enum.ErrInvalidInput, "發送者ID不能為空")
	}

	if message.TargetId == 0 {
		return appErrors.New(enum.ErrInvalidInput, "接收者ID不能為空")
	}

	// 設置消息類型
	message.Type = entities.MessageTypePrivate
	if message.Media == 0 {
		message.Media = entities.MediaTypeText
	}

	// 設置時間戳
	now := time.Now()
	message.CreatedAt = now
	message.UpdatedAt = now

	// 保存到數據庫
	if err := uc.messageRepo.Create(ctx, message); err != nil {
		return appErrors.Wrap(err, enum.ErrMessageSendFailed, map[string]interface{}{
			"userId":   message.UserId,
			"targetId": message.TargetId,
			"content":  message.Content,
		})
	}

	// 保存到緩存
	if err := uc.messageCache.StorePrivateMessage(ctx, message); err != nil {
		// 這裡僅記錄緩存錯誤，不影響主流程
		appErrors.WithDevMessage(err, "緩存私人消息失敗").LogError()
	}

	return nil
}

// SendGroupMessage 發送群組消息，使用新的錯誤處理系統
func (uc *EnhancedMessageUseCase) SendGroupMessage(ctx context.Context, message *entities.Message) error {
	if message == nil {
		return appErrors.New(enum.ErrInvalidInput, "消息不能為空")
	}

	if message.UserId == 0 {
		return appErrors.New(enum.ErrInvalidInput, "發送者ID不能為空")
	}

	if message.RoomID == 0 {
		return appErrors.New(enum.ErrInvalidInput, "聊天室ID不能為空")
	}

	// 設置消息類型
	message.Type = entities.MessageTypeGroup
	if message.Media == 0 {
		message.Media = entities.MediaTypeText
	}

	// 設置時間戳
	now := time.Now()
	message.CreatedAt = now
	message.UpdatedAt = now

	// 保存到數據庫
	if err := uc.messageRepo.Create(ctx, message); err != nil {
		return appErrors.Wrap(err, enum.ErrMessageSendFailed, map[string]interface{}{
			"userId":  message.UserId,
			"roomId":  message.RoomID,
			"content": message.Content,
		})
	}

	// 保存到緩存
	if err := uc.messageCache.StoreGroupMessage(ctx, message); err != nil {
		// 這裡僅記錄緩存錯誤，不影響主流程
		appErrors.WithDevMessage(err, "緩存群組消息失敗").LogError()
	}

	return nil
}

// GetPrivateMessages 獲取私人消息歷史記錄
func (uc *EnhancedMessageUseCase) GetPrivateMessages(ctx context.Context, fromUserID, toUserID uint) ([]*entities.Message, error) {
	// 首先嘗試從緩存獲取
	cacheMessages, err := uc.messageCache.GetPrivateMessages(ctx, fromUserID, toUserID)
	if err == nil && len(cacheMessages) > 0 {
		return cacheMessages, nil
	}

	// 如果緩存不可用或為空，從數據庫獲取
	if err != nil {
		// 記錄緩存錯誤，但繼續從數據庫獲取
		appErrors.WithDevMessage(err, "從緩存獲取私人消息失敗").LogError()
	}

	// 從數據庫獲取
	dbMessages, err := uc.messageRepo.FindMessagesBetweenUsers(ctx, fromUserID, toUserID)
	if err != nil {
		return nil, appErrors.Wrap(err, enum.ErrHistoryLoadFailed, map[string]interface{}{
			"fromUserId": fromUserID,
			"toUserId":   toUserID,
		})
	}

	// 將數據庫結果保存到緩存（異步）
	go func() {
		cacheCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		for _, msg := range dbMessages {
			if cacheErr := uc.messageCache.StorePrivateMessage(cacheCtx, msg); cacheErr != nil {
				// 記錄緩存錯誤
				appErrors.WithDevMessage(cacheErr, "將數據庫消息保存到緩存失敗").LogError()
			}
		}
	}()

	return dbMessages, nil
}

// GetGroupMessages 獲取群組消息歷史記錄
func (uc *EnhancedMessageUseCase) GetGroupMessages(ctx context.Context, roomID uint) ([]*entities.Message, error) {
	// 首先嘗試從緩存獲取
	cacheMessages, err := uc.messageCache.GetGroupMessages(ctx, roomID)
	if err == nil && len(cacheMessages) > 0 {
		return cacheMessages, nil
	}

	// 如果緩存不可用或為空，從數據庫獲取
	if err != nil {
		// 記錄緩存錯誤，但繼續從數據庫獲取
		appErrors.WithDevMessage(err, "從緩存獲取群組消息失敗").LogError()
	}

	// 從數據庫獲取
	dbMessages, err := uc.messageRepo.FindMessagesByRoomID(ctx, roomID)
	if err != nil {
		return nil, appErrors.Wrap(err, enum.ErrHistoryLoadFailed, map[string]interface{}{
			"roomId": roomID,
		})
	}

	// 將數據庫結果保存到緩存（異步）
	go func() {
		cacheCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		for _, msg := range dbMessages {
			if cacheErr := uc.messageCache.StoreGroupMessage(cacheCtx, msg); cacheErr != nil {
				// 記錄緩存錯誤
				appErrors.WithDevMessage(cacheErr, "將數據庫群組消息保存到緩存失敗").LogError()
			}
		}
	}()

	return dbMessages, nil
}
