package mysql

import (
	"clean-architecture-gochat/internal/domain/entities"
	"clean-architecture-gochat/internal/domain/repositories"
	"context"

	"gorm.io/gorm"
)

type chatRoomRepository struct {
	db *gorm.DB
}

func NewChatRoomRepository(db *gorm.DB) repositories.ChatRoomRepository {
	return &chatRoomRepository{db: db}
}

func (r *chatRoomRepository) Create(ctx context.Context, room *entities.ChatRoom) error {
	return r.db.WithContext(ctx).Create(room).Error
}

func (r *chatRoomRepository) FindByID(ctx context.Context, id uint) (*entities.ChatRoom, error) {
	var room entities.ChatRoom
	err := r.db.WithContext(ctx).First(&room, id).Error
	if err != nil {
		return nil, err
	}
	return &room, nil
}

func (r *chatRoomRepository) FindByParticipant(ctx context.Context, userID uint) ([]*entities.ChatRoom, error) {
	var rooms []*entities.ChatRoom
	err := r.db.WithContext(ctx).
		Joins("JOIN chat_room_participants ON chat_rooms.id = chat_room_participants.chat_room_id").
		Where("chat_room_participants.user_id = ?", userID).
		Find(&rooms).Error
	if err != nil {
		return nil, err
	}
	return rooms, nil
}

func (r *chatRoomRepository) AddParticipant(ctx context.Context, roomID, userID uint) error {
	return r.db.WithContext(ctx).Exec(
		"INSERT INTO chat_room_participants (chat_room_id, user_id) VALUES (?, ?)",
		roomID, userID,
	).Error
}

func (r *chatRoomRepository) RemoveParticipant(ctx context.Context, roomID, userID uint) error {
	return r.db.WithContext(ctx).Exec(
		"DELETE FROM chat_room_participants WHERE chat_room_id = ? AND user_id = ?",
		roomID, userID,
	).Error
}
