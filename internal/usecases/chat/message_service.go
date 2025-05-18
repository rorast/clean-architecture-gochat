package chat

import (
	"clean-architecture-gochat/internal/domain/entities"
	"clean-architecture-gochat/internal/domain/repositories"
	"context"
)

type messageService struct {
	messageRepo repositories.MessageRepository
}

func NewMessageService(messageRepo repositories.MessageRepository) repositories.MessageService {
	return &messageService{
		messageRepo: messageRepo,
	}
}

func (s *messageService) SendPrivateMessage(ctx context.Context, message *entities.Message) error {
	return s.messageRepo.Create(ctx, message)
}

func (s *messageService) SendGroupMessage(ctx context.Context, message *entities.Message) error {
	return s.messageRepo.Create(ctx, message)
}

func (s *messageService) GetPrivateHistory(ctx context.Context, fromUserID, toUserID uint) ([]*entities.Message, error) {
	return s.messageRepo.FindMessagesBetweenUsers(ctx, fromUserID, toUserID)
}

func (s *messageService) GetGroupHistory(ctx context.Context, roomID uint) ([]*entities.Message, error) {
	return s.messageRepo.FindMessagesByRoomID(ctx, roomID)
}

func (s *messageService) GetRecentMessages(ctx context.Context, userID uint) ([]*entities.Message, error) {
	return s.messageRepo.FindMessagesByUserID(ctx, userID)
}
