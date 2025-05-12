package test

import (
	"clean-architecture-gochat/internal/domain/entities"
	"clean-architecture-gochat/internal/domain/repositories"
	"clean-architecture-gochat/usecases/chat"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:123456@tcp(127.0.0.1:3308)/newgochat?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}
	return db
}

func TestGroupChat(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()

	// 初始化所需的倉儲和服務
	groupRepo := repositories.NewGroupRepository(db)
	messageRepo := repositories.NewMessageRepository(db)
	groupChatService := chat.NewGroupChatService(groupRepo, messageRepo)

	// 1. 測試創建群組
	group, err := groupChatService.CreateGroup(ctx, "測試群組", 1, 1, "這是一個測試群組", 0, 0)
	assert.NoError(t, err)
	assert.NotNil(t, group)
	assert.Equal(t, "測試群組", group.Name)
	assert.Equal(t, uint(1), group.OwnerId)

	// 2. 測試添加群組成員
	err = groupChatService.AddMember(ctx, group.ID, 2)
	assert.NoError(t, err)

	// 3. 測試發送群組消息
	message := &entities.Message{
		UserId:    1,
		TargetId:  group.ID,
		RoomID:    group.ID,
		Type:      entities.MessageTypeGroup,
		Media:     entities.MediaTypeText,
		Content:   "大家好！",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err = groupChatService.SendGroupMessage(ctx, message)
	assert.NoError(t, err)

	// 4. 測試獲取群組成員列表
	members, err := groupChatService.GetGroupMembers(ctx, group.ID)
	assert.NoError(t, err)
	assert.Contains(t, members, uint(1)) // 包含群主
	assert.Contains(t, members, uint(2)) // 包含新加入的成員

	// 5. 測試獲取群組消息歷史
	messages, err := groupChatService.GetGroupHistory(ctx, group.ID)
	assert.NoError(t, err)
	assert.NotEmpty(t, messages)
	assert.Equal(t, "大家好！", messages[0].Content)

	// 6. 測試移除群組成員
	err = groupChatService.RemoveMember(ctx, group.ID, 2)
	assert.NoError(t, err)

	// 驗證成員已被移除
	members, err = groupChatService.GetGroupMembers(ctx, group.ID)
	assert.NoError(t, err)
	assert.NotContains(t, members, uint(2))

	// 7. 測試刪除群組
	err = groupChatService.DeleteGroup(ctx, group.ID)
	assert.NoError(t, err)

	// 驗證群組已被刪除
	_, err = groupChatService.GetGroup(ctx, group.ID)
	assert.Error(t, err)
}
