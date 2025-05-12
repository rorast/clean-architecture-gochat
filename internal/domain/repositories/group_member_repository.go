package repositories

import (
	"clean-architecture-gochat/internal/domain/entities"
	"context"

	"gorm.io/gorm"
)

type GroupMemberRepository interface {
	Create(ctx context.Context, member *entities.GroupMember) error
	Delete(ctx context.Context, groupID, userID uint) error
	GetMembers(ctx context.Context, groupID uint) ([]uint, error)
	IsMember(ctx context.Context, groupID, userID uint) (bool, error)
	GetUserGroups(ctx context.Context, userID uint) ([]uint, error)
}

func NewGroupMemberRepository(db *gorm.DB) GroupMemberRepository {
	return &groupMemberRepository{db: db}
}

type groupMemberRepository struct {
	db *gorm.DB
}

func (r *groupMemberRepository) Create(ctx context.Context, member *entities.GroupMember) error {
	return r.db.WithContext(ctx).Create(member).Error
}

func (r *groupMemberRepository) Delete(ctx context.Context, groupID, userID uint) error {
	return r.db.WithContext(ctx).Delete(&entities.GroupMember{}, "group_id = ? AND user_id = ?", groupID, userID).Error
}

func (r *groupMemberRepository) GetMembers(ctx context.Context, groupID uint) ([]uint, error) {
	var members []uint
	err := r.db.WithContext(ctx).Model(&entities.GroupMember{}).
		Where("group_id = ?", groupID).
		Pluck("user_id", &members).Error
	if err != nil {
		return nil, err
	}
	return members, nil
}

func (r *groupMemberRepository) IsMember(ctx context.Context, groupID, userID uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entities.GroupMember{}).
		Where("group_id = ? AND user_id = ?", groupID, userID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *groupMemberRepository) GetUserGroups(ctx context.Context, userID uint) ([]uint, error) {
	var groups []uint
	err := r.db.WithContext(ctx).Model(&entities.GroupMember{}).
		Where("user_id = ?", userID).
		Pluck("group_id", &groups).Error
	if err != nil {
		return nil, err
	}
	return groups, nil
}
