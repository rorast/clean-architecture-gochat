package repositories

import (
	"clean-architecture-gochat/internal/domain/entities"
	"context"

	"gorm.io/gorm"
)

type MessageRepository interface {
	SaveMessage(ctx context.Context, message *entities.Message) error
	Create(ctx context.Context, message *entities.Message) error
	FindByID(ctx context.Context, id uint) (*entities.Message, error)
	FindMessagesByUserID(ctx context.Context, userId uint) ([]*entities.Message, error)
	FindMessagesBetweenUsers(ctx context.Context, userID, targetID uint) ([]*entities.Message, error)
	FindMessagesByRoomID(ctx context.Context, roomID uint) ([]*entities.Message, error)
	Delete(ctx context.Context, id uint) error
}

type messageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) MessageRepository {
	return &messageRepository{db: db}
}

func (r *messageRepository) SaveMessage(ctx context.Context, message *entities.Message) error {
	return r.db.WithContext(ctx).Save(message).Error
}

func (r *messageRepository) Create(ctx context.Context, message *entities.Message) error {
	return r.db.WithContext(ctx).Create(message).Error
}

func (r *messageRepository) FindByID(ctx context.Context, id uint) (*entities.Message, error) {
	var message entities.Message
	err := r.db.WithContext(ctx).First(&message, id).Error
	return &message, err
}

func (r *messageRepository) FindMessagesByUserID(ctx context.Context, userId uint) ([]*entities.Message, error) {
	var messages []*entities.Message
	err := r.db.WithContext(ctx).Where("user_id = ?", userId).Find(&messages).Error
	return messages, err
}

func (r *messageRepository) FindMessagesBetweenUsers(ctx context.Context, userID, targetID uint) ([]*entities.Message, error) {
	var messages []*entities.Message
	err := r.db.WithContext(ctx).
		Where("(user_id = ? AND target_id = ?) OR (user_id = ? AND target_id = ?)", userID, targetID, targetID, userID).
		Order("created_at ASC").
		Find(&messages).Error
	return messages, err
}

func (r *messageRepository) FindMessagesByRoomID(ctx context.Context, roomID uint) ([]*entities.Message, error) {
	var messages []*entities.Message
	err := r.db.WithContext(ctx).Where("room_id = ?", roomID).Find(&messages).Error
	return messages, err
}

func (r *messageRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&entities.Message{}, id).Error
}
