package repositories

import (
	"clean-architecture-gochat/internal/domain/entities"
	"context"

	"gorm.io/gorm"
)

type GroupRepository interface {
	Create(ctx context.Context, group *entities.Group) error
	FindByID(ctx context.Context, id uint) (*entities.Group, error)
	FindByOwnerID(ctx context.Context, ownerId uint) ([]*entities.Group, error)
	FindByType(ctx context.Context, groupType int) ([]*entities.Group, error)
	FindBySize(ctx context.Context, size int) ([]*entities.Group, error)
	Update(ctx context.Context, group *entities.Group) error
	Delete(ctx context.Context, id uint) error
	AddMember(ctx context.Context, groupID, userID uint) error
	RemoveMember(ctx context.Context, groupID, userID uint) error
	GetMembers(ctx context.Context, groupID uint) ([]uint, error)
	IsMember(ctx context.Context, groupID, userID uint) (bool, error)
	FindJoinedGroups(ctx context.Context, userID uint) ([]*entities.Group, error)
}

func NewGroupRepository(db *gorm.DB) GroupRepository {
	return &groupRepository{db: db}
}

type groupRepository struct {
	db *gorm.DB
}

func (r *groupRepository) Create(ctx context.Context, group *entities.Group) error {
	return r.db.WithContext(ctx).Create(group).Error
}

func (r *groupRepository) FindByID(ctx context.Context, id uint) (*entities.Group, error) {
	var group entities.Group
	err := r.db.WithContext(ctx).First(&group, id).Error
	if err != nil {
		return nil, err
	}
	return &group, nil
}

func (r *groupRepository) FindByOwnerID(ctx context.Context, ownerId uint) ([]*entities.Group, error) {
	var groups []*entities.Group
	err := r.db.WithContext(ctx).Where("owner_id = ?", ownerId).Find(&groups).Error
	if err != nil {
		return nil, err
	}
	return groups, nil
}

func (r *groupRepository) FindByType(ctx context.Context, groupType int) ([]*entities.Group, error) {
	var groups []*entities.Group
	err := r.db.WithContext(ctx).Where("type = ?", groupType).Find(&groups).Error
	if err != nil {
		return nil, err
	}
	return groups, nil
}

func (r *groupRepository) FindBySize(ctx context.Context, size int) ([]*entities.Group, error) {
	var groups []*entities.Group
	err := r.db.WithContext(ctx).Where("size = ?", size).Find(&groups).Error
	if err != nil {
		return nil, err
	}
	return groups, nil
}

func (r *groupRepository) Update(ctx context.Context, group *entities.Group) error {
	return r.db.WithContext(ctx).Save(group).Error
}

func (r *groupRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&entities.Group{}, id).Error
}

func (r *groupRepository) AddMember(ctx context.Context, groupID, userID uint) error {
	return r.db.WithContext(ctx).Exec(
		"INSERT INTO group_members (group_id, user_id) VALUES (?, ?)",
		groupID, userID,
	).Error
}

func (r *groupRepository) RemoveMember(ctx context.Context, groupID, userID uint) error {
	return r.db.WithContext(ctx).Exec(
		"DELETE FROM group_members WHERE group_id = ? AND user_id = ?",
		groupID, userID,
	).Error
}

func (r *groupRepository) GetMembers(ctx context.Context, groupID uint) ([]uint, error) {
	var members []uint
	err := r.db.WithContext(ctx).Table("group_members").
		Where("group_id = ?", groupID).
		Pluck("user_id", &members).Error
	if err != nil {
		return nil, err
	}
	return members, nil
}

func (r *groupRepository) IsMember(ctx context.Context, groupID, userID uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Table("group_members").
		Where("group_id = ? AND user_id = ?", groupID, userID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *groupRepository) FindJoinedGroups(ctx context.Context, userID uint) ([]*entities.Group, error) {
	var groups []*entities.Group
	err := r.db.WithContext(ctx).
		Joins("JOIN group_members ON groups.id = group_members.group_id").
		Where("group_members.user_id = ?", userID).
		Find(&groups).Error
	if err != nil {
		return nil, err
	}
	return groups, nil
}
