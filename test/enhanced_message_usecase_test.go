package test

import (
	"context"
	"testing"
	"time"

	"clean-architecture-gochat/internal/domain/entities"
	"clean-architecture-gochat/internal/usecases/message"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// 模擬 MessageRepository
type MockMessageRepository struct {
	mock.Mock
}

func (m *MockMessageRepository) SaveMessage(ctx context.Context, message *entities.Message) error {
	args := m.Called(ctx, message)
	return args.Error(0)
}

func (m *MockMessageRepository) Create(ctx context.Context, message *entities.Message) error {
	args := m.Called(ctx, message)
	return args.Error(0)
}

func (m *MockMessageRepository) FindByID(ctx context.Context, id uint) (*entities.Message, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Message), args.Error(1)
}

func (m *MockMessageRepository) FindMessagesByUserID(ctx context.Context, userId uint) ([]*entities.Message, error) {
	args := m.Called(ctx, userId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.Message), args.Error(1)
}

func (m *MockMessageRepository) FindMessagesBetweenUsers(ctx context.Context, userID, targetID uint) ([]*entities.Message, error) {
	args := m.Called(ctx, userID, targetID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.Message), args.Error(1)
}

func (m *MockMessageRepository) FindMessagesByRoomID(ctx context.Context, roomID uint) ([]*entities.Message, error) {
	args := m.Called(ctx, roomID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.Message), args.Error(1)
}

func (m *MockMessageRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// 模擬 MessageCacheRepository
type MockMessageCacheRepository struct {
	mock.Mock
}

func (m *MockMessageCacheRepository) StorePrivateMessage(ctx context.Context, message *entities.Message) error {
	args := m.Called(ctx, message)
	return args.Error(0)
}

func (m *MockMessageCacheRepository) StoreGroupMessage(ctx context.Context, message *entities.Message) error {
	args := m.Called(ctx, message)
	return args.Error(0)
}

func (m *MockMessageCacheRepository) GetPrivateMessages(ctx context.Context, fromUserID, toUserID uint) ([]*entities.Message, error) {
	args := m.Called(ctx, fromUserID, toUserID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.Message), args.Error(1)
}

func (m *MockMessageCacheRepository) GetGroupMessages(ctx context.Context, roomID uint) ([]*entities.Message, error) {
	args := m.Called(ctx, roomID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.Message), args.Error(1)
}

// CleanExpiredMessages 清理過期訊息
func (m *MockMessageCacheRepository) CleanExpiredMessages(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// GetUserMessageList 獲取用戶訊息列表
func (m *MockMessageCacheRepository) GetUserMessageList(ctx context.Context, userID uint) ([]*entities.Message, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.Message), args.Error(1)
}

// 測試發送私人訊息
func TestEnhancedMessageUseCase_SendPrivateMessage(t *testing.T) {
	mockRepo := new(MockMessageRepository)
	mockCache := new(MockMessageCacheRepository)
	useCase := message.NewEnhancedMessageUseCase(mockRepo, mockCache)

	ctx := context.Background()
	msg := &entities.Message{
		UserId:   1,
		TargetId: 2,
		Content:  "Hello, how are you?",
	}

	// 設置模擬的行為
	mockRepo.On("Create", ctx, mock.AnythingOfType("*entities.Message")).Return(nil)
	mockCache.On("StorePrivateMessage", ctx, mock.AnythingOfType("*entities.Message")).Return(nil)

	// 執行測試
	err := useCase.SendPrivateMessage(ctx, msg)

	// 驗證結果
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)

	// 驗證訊息類型和媒體類型被正確設置
	assert.Equal(t, entities.MessageTypePrivate, msg.Type)
	assert.Equal(t, entities.MediaTypeText, msg.Media)
}

// 測試發送群組訊息
func TestEnhancedMessageUseCase_SendGroupMessage(t *testing.T) {
	mockRepo := new(MockMessageRepository)
	mockCache := new(MockMessageCacheRepository)
	useCase := message.NewEnhancedMessageUseCase(mockRepo, mockCache)

	ctx := context.Background()
	msg := &entities.Message{
		UserId:  1,
		RoomID:  101,
		Content: "Hello everyone!",
	}

	// 設置模擬的行為
	mockRepo.On("Create", ctx, mock.AnythingOfType("*entities.Message")).Return(nil)
	mockCache.On("StoreGroupMessage", ctx, mock.AnythingOfType("*entities.Message")).Return(nil)

	// 執行測試
	err := useCase.SendGroupMessage(ctx, msg)

	// 驗證結果
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)

	// 驗證訊息類型和媒體類型被正確設置
	assert.Equal(t, entities.MessageTypeGroup, msg.Type)
	assert.Equal(t, entities.MediaTypeText, msg.Media)
}

// 測試從緩存獲取私人訊息
func TestEnhancedMessageUseCase_GetPrivateMessages_FromCache(t *testing.T) {
	mockRepo := new(MockMessageRepository)
	mockCache := new(MockMessageCacheRepository)
	useCase := message.NewEnhancedMessageUseCase(mockRepo, mockCache)

	ctx := context.Background()
	userID, targetID := uint(1), uint(2)

	// 模擬緩存中有數據
	cachedMessages := []*entities.Message{
		{
			ID:        1,
			UserId:    userID,
			TargetId:  targetID,
			Content:   "Hello from cache",
			Type:      entities.MessageTypePrivate,
			CreatedAt: time.Now(),
		},
	}

	mockCache.On("GetPrivateMessages", ctx, userID, targetID).Return(cachedMessages, nil)

	// 執行測試
	messages, err := useCase.GetPrivateMessages(ctx, userID, targetID)

	// 驗證結果
	assert.NoError(t, err)
	assert.Len(t, messages, 1)
	assert.Equal(t, "Hello from cache", messages[0].Content)

	// 確認沒有調用數據庫查詢
	mockRepo.AssertNotCalled(t, "FindMessagesBetweenUsers")
}

// 測試從數據庫獲取私人訊息 (當緩存中沒有數據時)
func TestEnhancedMessageUseCase_GetPrivateMessages_FromDB(t *testing.T) {
	mockRepo := new(MockMessageRepository)
	mockCache := new(MockMessageCacheRepository)
	useCase := message.NewEnhancedMessageUseCase(mockRepo, mockCache)

	ctx := context.Background()
	userID, targetID := uint(1), uint(2)

	// 模擬緩存中沒有數據
	mockCache.On("GetPrivateMessages", ctx, userID, targetID).Return([]*entities.Message{}, nil)

	// 模擬數據庫中有數據
	dbMessages := []*entities.Message{
		{
			ID:        1,
			UserId:    userID,
			TargetId:  targetID,
			Content:   "Hello from database",
			Type:      entities.MessageTypePrivate,
			CreatedAt: time.Now(),
		},
	}
	mockRepo.On("FindMessagesBetweenUsers", ctx, userID, targetID).Return(dbMessages, nil)

	// 模擬異步緩存存儲 (我們不能直接測試 Go 的協程，但可以確保方法被調用)
	mockCache.On("StorePrivateMessage", mock.Anything, mock.AnythingOfType("*entities.Message")).Return(nil)

	// 執行測試
	messages, err := useCase.GetPrivateMessages(ctx, userID, targetID)

	// 驗證結果
	assert.NoError(t, err)
	assert.Len(t, messages, 1)
	assert.Equal(t, "Hello from database", messages[0].Content)

	// 確認調用了數據庫查詢
	mockRepo.AssertCalled(t, "FindMessagesBetweenUsers", ctx, userID, targetID)

	// 注意：我們不能直接驗證異步調用，因為它在另一個 goroutine 中執行
}

// 測試獲取群組訊息的邏輯
func TestEnhancedMessageUseCase_GetGroupMessages(t *testing.T) {
	mockRepo := new(MockMessageRepository)
	mockCache := new(MockMessageCacheRepository)
	useCase := message.NewEnhancedMessageUseCase(mockRepo, mockCache)

	ctx := context.Background()
	roomID := uint(101)

	// 模擬緩存中有數據
	cachedMessages := []*entities.Message{
		{
			ID:      1,
			UserId:  1,
			RoomID:  roomID,
			Content: "Group message from cache",
			Type:    entities.MessageTypeGroup,
		},
	}

	mockCache.On("GetGroupMessages", ctx, roomID).Return(cachedMessages, nil)

	// 執行測試
	messages, err := useCase.GetGroupMessages(ctx, roomID)

	// 驗證結果
	assert.NoError(t, err)
	assert.Len(t, messages, 1)
	assert.Equal(t, "Group message from cache", messages[0].Content)

	// 確認沒有調用數據庫查詢
	mockRepo.AssertNotCalled(t, "FindMessagesByRoomID")
}

// 測試發送訊息時緩存寫入失敗的情況
func TestEnhancedMessageUseCase_SendPrivateMessage_CacheError(t *testing.T) {
	mockRepo := new(MockMessageRepository)
	mockCache := new(MockMessageCacheRepository)
	useCase := message.NewEnhancedMessageUseCase(mockRepo, mockCache)

	ctx := context.Background()
	msg := &entities.Message{
		UserId:   1,
		TargetId: 2,
		Content:  "Test cache error handling",
	}

	// 設置模擬的行為：資料庫寫入成功，但緩存寫入失敗
	mockRepo.On("Create", ctx, mock.AnythingOfType("*entities.Message")).Return(nil)
	mockCache.On("StorePrivateMessage", ctx, mock.AnythingOfType("*entities.Message")).Return(assert.AnError)

	// 執行測試
	err := useCase.SendPrivateMessage(ctx, msg)

	// 驗證結果：即使緩存失敗，整體操作應該仍然成功
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)

	// 驗證訊息類型和媒體類型被正確設置
	assert.Equal(t, entities.MessageTypePrivate, msg.Type)
	assert.Equal(t, entities.MediaTypeText, msg.Media)
}

// 測試當緩存讀取失敗時，從資料庫獲取私人訊息的情況
func TestEnhancedMessageUseCase_GetPrivateMessages_CacheError(t *testing.T) {
	mockRepo := new(MockMessageRepository)
	mockCache := new(MockMessageCacheRepository)
	useCase := message.NewEnhancedMessageUseCase(mockRepo, mockCache)

	ctx := context.Background()
	userID, targetID := uint(1), uint(2)

	// 模擬緩存讀取失敗
	mockCache.On("GetPrivateMessages", ctx, userID, targetID).Return(nil, assert.AnError)

	// 模擬數據庫中有數據
	dbMessages := []*entities.Message{
		{
			ID:        1,
			UserId:    userID,
			TargetId:  targetID,
			Content:   "Hello from database when cache fails",
			Type:      entities.MessageTypePrivate,
			CreatedAt: time.Now(),
		},
	}
	mockRepo.On("FindMessagesBetweenUsers", ctx, userID, targetID).Return(dbMessages, nil)

	// 模擬異步緩存存儲
	mockCache.On("StorePrivateMessage", mock.Anything, mock.AnythingOfType("*entities.Message")).Return(nil)

	// 執行測試
	messages, err := useCase.GetPrivateMessages(ctx, userID, targetID)

	// 驗證結果：即使緩存失敗，應該仍能從資料庫獲取數據
	assert.NoError(t, err)
	assert.Len(t, messages, 1)
	assert.Equal(t, "Hello from database when cache fails", messages[0].Content)

	// 確認調用了數據庫查詢
	mockRepo.AssertCalled(t, "FindMessagesBetweenUsers", ctx, userID, targetID)
}
