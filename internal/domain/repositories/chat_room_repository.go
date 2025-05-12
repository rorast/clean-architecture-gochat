package repositories

import (
	"clean-architecture-gochat/internal/domain/entities"
	"context"
)

type ChatRoomRepository interface {
	Create(ctx context.Context, room *entities.ChatRoom) error
	FindByID(ctx context.Context, id uint) (*entities.ChatRoom, error)
	FindByParticipant(ctx context.Context, userID uint) ([]*entities.ChatRoom, error)
	AddParticipant(ctx context.Context, roomID, userID uint) error
	RemoveParticipant(ctx context.Context, roomID, userID uint) error
}
