package chat

import (
	"clean-architecture-gochat/internal/domain/cache"
	"clean-architecture-gochat/internal/domain/entities"
	"clean-architecture-gochat/internal/domain/repositories"
	"context"
	"fmt"
)

type MessageUseCase interface {
	// 發送私聊訊息
	SendPrivateMessage(ctx context.Context, message *entities.Message) error
	// 發送群聊訊息
	SendGroupMessage(ctx context.Context, message *entities.Message) error
	// 獲取私聊訊息歷史
	GetPrivateMessageHistory(ctx context.Context, fromUserID, toUserID uint) ([]*entities.Message, error)
	// 獲取群聊訊息歷史
	GetGroupMessageHistory(ctx context.Context, roomID uint) ([]*entities.Message, error)
	// 獲取用戶的最近訊息
	GetUserRecentMessages(ctx context.Context, userID uint) ([]*entities.Message, error)
}

type messageUseCase struct {
	messageRepo      repositories.MessageRepository
	messageCacheRepo cache.MessageCacheRepository
	groupRepo        repositories.GroupRepository
}

// NewMessageUseCase 創建新的訊息用例
func NewMessageUseCase(
	messageRepo repositories.MessageRepository,
	messageCacheRepo cache.MessageCacheRepository,
	groupRepo repositories.GroupRepository,
) MessageUseCase {
	return &messageUseCase{
		messageRepo:      messageRepo,
		messageCacheRepo: messageCacheRepo,
		groupRepo:        groupRepo,
	}
}

func (uc *messageUseCase) SendPrivateMessage(ctx context.Context, message *entities.Message) error {
	// 1. 儲存訊息到資料庫
	if err := uc.messageRepo.Create(ctx, message); err != nil {
		return fmt.Errorf("儲存訊息到資料庫失敗: %v", err)
	}

	// 2. 儲存訊息到快取
	if err := uc.messageCacheRepo.StorePrivateMessage(ctx, message); err != nil {
		// 快取失敗不影響主要功能，只記錄錯誤
		fmt.Printf("儲存訊息到快取失敗: %v\n", err)
	}

	return nil
}

func (uc *messageUseCase) SendGroupMessage(ctx context.Context, message *entities.Message) error {
	// 1. 檢查群組是否存在
	group, err := uc.groupRepo.FindByID(ctx, message.RoomID)
	if err != nil {
		return fmt.Errorf("群組不存在: %v", err)
	}

	// 2. 檢查發送者是否為群組成員
	isMember, err := uc.groupRepo.IsMember(ctx, group.ID, message.UserId)
	if err != nil {
		return fmt.Errorf("檢查群組成員失敗: %v", err)
	}
	if !isMember {
		return fmt.Errorf("用戶不是群組成員")
	}

	// 3. 儲存訊息到資料庫
	if err := uc.messageRepo.Create(ctx, message); err != nil {
		return fmt.Errorf("儲存訊息到資料庫失敗: %v", err)
	}

	// 4. 儲存訊息到快取
	if err := uc.messageCacheRepo.StoreGroupMessage(ctx, message); err != nil {
		// 快取失敗不影響主要功能，只記錄錯誤
		fmt.Printf("儲存訊息到快取失敗: %v\n", err)
	}

	return nil
}

func (uc *messageUseCase) GetPrivateMessageHistory(ctx context.Context, fromUserID, toUserID uint) ([]*entities.Message, error) {
	// 1. 先嘗試從快取獲取
	messages, err := uc.messageCacheRepo.GetPrivateMessages(ctx, fromUserID, toUserID)
	if err == nil && len(messages) > 0 {
		return messages, nil
	}

	// 2. 快取未命中，從資料庫獲取
	messages, err = uc.messageRepo.FindMessagesBetweenUsers(ctx, fromUserID, toUserID)
	if err != nil {
		return nil, fmt.Errorf("從資料庫獲取訊息失敗: %v", err)
	}

	// 3. 更新快取（非阻塞）
	go func() {
		for _, msg := range messages {
			if err := uc.messageCacheRepo.StorePrivateMessage(context.Background(), msg); err != nil {
				fmt.Printf("更新訊息快取失敗: %v\n", err)
			}
		}
	}()

	return messages, nil
}

func (uc *messageUseCase) GetGroupMessageHistory(ctx context.Context, roomID uint) ([]*entities.Message, error) {
	// 1. 先嘗試從快取獲取
	messages, err := uc.messageCacheRepo.GetGroupMessages(ctx, roomID)
	if err == nil && len(messages) > 0 {
		return messages, nil
	}

	// 2. 快取未命中，從資料庫獲取
	messages, err = uc.messageRepo.FindMessagesByRoomID(ctx, roomID)
	if err != nil {
		return nil, fmt.Errorf("從資料庫獲取訊息失敗: %v", err)
	}

	// 3. 更新快取（非阻塞）
	go func() {
		for _, msg := range messages {
			if err := uc.messageCacheRepo.StoreGroupMessage(context.Background(), msg); err != nil {
				fmt.Printf("更新訊息快取失敗: %v\n", err)
			}
		}
	}()

	return messages, nil
}

func (uc *messageUseCase) GetUserRecentMessages(ctx context.Context, userID uint) ([]*entities.Message, error) {
	// 1. 先嘗試從快取獲取
	messages, err := uc.messageCacheRepo.GetUserMessageList(ctx, userID)
	if err == nil && len(messages) > 0 {
		return messages, nil
	}

	// 2. 快取未命中，從資料庫獲取
	messages, err = uc.messageRepo.FindMessagesByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("從資料庫獲取訊息失敗: %v", err)
	}

	// 3. 更新快取（非阻塞）
	go func() {
		for _, msg := range messages {
			if err := uc.messageCacheRepo.StorePrivateMessage(context.Background(), msg); err != nil {
				fmt.Printf("更新訊息快取失敗: %v\n", err)
			}
		}
	}()

	return messages, nil
}
