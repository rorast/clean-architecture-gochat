package websocket

import (
	"context"
	"fmt"

	websocketInfra "clean-architecture-gochat/infrastructure/websocket"

	ws "github.com/gorilla/websocket"
)

type ConnectionService interface {
	Connect(ctx context.Context, conn *ws.Conn, userID uint) error
	Disconnect(ctx context.Context, userID uint) error
	SendToUser(ctx context.Context, userID uint, message []byte) error
	Broadcast(ctx context.Context, message []byte) error
	IsUserOnline(ctx context.Context, userID uint) bool
}

type connectionService struct {
	hub *websocketInfra.Hub
}

func NewConnectionService() ConnectionService {
	return &connectionService{
		hub: websocketInfra.GetHub(),
	}
}

func (s *connectionService) Connect(ctx context.Context, conn *ws.Conn, userID uint) error {
	client := websocketInfra.UpgradeConnection(conn, fmt.Sprintf("%d", userID))
	if client == nil {
		return fmt.Errorf("failed to upgrade WebSocket connection")
	}
	return nil
}

func (s *connectionService) Disconnect(ctx context.Context, userID uint) error {
	client, exists := websocketInfra.GetClient(fmt.Sprintf("%d", userID))
	if exists {
		client.Hub.Unregister <- client
	}
	return nil
}

func (s *connectionService) SendToUser(ctx context.Context, userID uint, message []byte) error {
	client, exists := websocketInfra.GetClient(fmt.Sprintf("%d", userID))
	if !exists {
		return fmt.Errorf("user %d is not online", userID)
	}
	client.Send <- message
	return nil
}

func (s *connectionService) Broadcast(ctx context.Context, message []byte) error {
	websocketInfra.Broadcast(message)
	return nil
}

func (s *connectionService) IsUserOnline(ctx context.Context, userID uint) bool {
	_, exists := websocketInfra.GetClient(fmt.Sprintf("%d", userID))
	return exists
}
