package repositories

import (
	"clean-architecture-gochat/internal/domain/entities"
	"context"

	"gorm.io/gorm"
)

type ContactRepository interface {
	FindFriendContacts(ctx context.Context, userId uint) ([]entities.Contact, error)
	CreateContact(ctx context.Context, userId, targetId uint, contactType int) error
	DeleteContact(ctx context.Context, userId, targetId uint) error
}

type contactRepository struct {
	db *gorm.DB
}

func NewContactRepository(db *gorm.DB) ContactRepository {
	return &contactRepository{db: db}
}

func (r *contactRepository) FindFriendContacts(ctx context.Context, userId uint) ([]entities.Contact, error) {
	var contacts []entities.Contact
	if err := r.db.WithContext(ctx).Where("owner_id = ? AND type = 1", userId).Find(&contacts).Error; err != nil {
		return nil, err
	}
	return contacts, nil
}

func (r *contactRepository) CreateContact(ctx context.Context, userId, targetId uint, contactType int) error {
	return r.db.WithContext(ctx).Create(&entities.Contact{
		OwnerId:  userId,
		TargetId: targetId,
		Type:     contactType,
	}).Error
}

func (r *contactRepository) DeleteContact(ctx context.Context, userId, targetId uint) error {
	return r.db.WithContext(ctx).Where("owner_id = ? AND target_id = ?", userId, targetId).Delete(&entities.Contact{}).Error
}
