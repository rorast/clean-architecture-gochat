package chat

import (
	websocketInfra "clean-architecture-gochat/infrastructure/websocket"
	"clean-architecture-gochat/internal/domain/entities"
	"clean-architecture-gochat/internal/domain/repositories"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	ws "github.com/gorilla/websocket"
)

type Service interface {
	UpgradeConnection(ctx context.Context, conn *ws.Conn, userID int64) error
	SendMessage(ctx context.Context, msg *entities.Message) error
	BroadcastMessage(ctx context.Context, msg *entities.Message) error
	GetChatHistory(ctx context.Context, userID, targetID int64) ([]*entities.Message, error)
}

type service struct {
	repo repositories.MessageRepository
}

func NewPrivateChatService(repo repositories.MessageRepository) Service {
	return &service{repo: repo}
}

// 升級 WebSocket 連線，並加入 WebSocket Hub
func (s *service) UpgradeConnection(ctx context.Context, conn *ws.Conn, userID int64) error {
	client := websocketInfra.UpgradeConnection(conn, fmt.Sprintf("%d", userID))
	if client == nil {
		return errors.New("failed to upgrade WebSocket connection")
	}
	return nil
}

// 發送訊息到單一 WebSocket 客戶端
func (s *service) SendMessage(ctx context.Context, msg *entities.Message) error {
	client, exists := websocketInfra.GetClient(fmt.Sprintf("%d", msg.TargetId))
	if !exists {
		return errors.New("target user is not online")
	}

	//msg.CreateTime = uint64(time.Now().Unix())
	msg.CreatedAt = time.Now()
	if err := s.repo.SaveMessage(ctx, msg); err != nil {
		return err
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	client.SendMessage(data)
	return nil
}

// 廣播訊息到所有 WebSocket 連線
func (s *service) BroadcastMessage(ctx context.Context, msg *entities.Message) error {
	msg.CreatedAt = time.Now()
	if err := s.repo.SaveMessage(ctx, msg); err != nil {
		return err
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	websocketInfra.Broadcast(data)
	return nil
}

// 獲取聊天歷史紀錄
func (s *service) GetChatHistory(ctx context.Context, userID, targetID int64) ([]*entities.Message, error) {
	//return s.repo.GetMessagesBetweenUsers(ctx, uint(userID), uint(targetID))
	return s.repo.FindMessagesBetweenUsers(ctx, uint(userID), uint(targetID))
}
