package chat

import (
	"clean-architecture-gochat/internal/domain/entities"
	"clean-architecture-gochat/internal/domain/repositories"
	"context"
	"errors"
	"time"
)

type GroupChatService interface {
	CreateGroup(ctx context.Context, name string, ownerID uint, groupType int, desc string, size int, joinType int) (*entities.Group, error)
	UpdateGroup(ctx context.Context, group *entities.Group) error
	DeleteGroup(ctx context.Context, groupID uint) error
	GetGroup(ctx context.Context, groupID uint) (*entities.Group, error)
	GetUserGroups(ctx context.Context, userID uint) ([]*entities.Group, error)
	GetGroupsByType(ctx context.Context, groupType int) ([]*entities.Group, error)
	GetGroupsBySize(ctx context.Context, size int) ([]*entities.Group, error)
	AddMember(ctx context.Context, groupID, userID uint) error
	RemoveMember(ctx context.Context, groupID, userID uint) error
	GetGroupMembers(ctx context.Context, groupID uint) ([]uint, error)
	IsGroupMember(ctx context.Context, groupID, userID uint) (bool, error)
	SendGroupMessage(ctx context.Context, message *entities.Message) error
	GetGroupHistory(ctx context.Context, groupID uint) ([]*entities.Message, error)
}

type groupChatService struct {
	groupRepo   repositories.GroupRepository
	messageRepo repositories.MessageRepository
}

func NewGroupChatService(groupRepo repositories.GroupRepository, messageRepo repositories.MessageRepository) GroupChatService {
	return &groupChatService{
		groupRepo:   groupRepo,
		messageRepo: messageRepo,
	}
}

func (s *groupChatService) CreateGroup(ctx context.Context, name string, ownerID uint, groupType int, desc string, size int, joinType int) (*entities.Group, error) {
	group := &entities.Group{
		Name:      name,
		OwnerId:   ownerID,
		Type:      groupType,
		Desc:      desc,
		Size:      size,
		JoinType:  joinType,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.groupRepo.Create(ctx, group); err != nil {
		return nil, err
	}

	// 添加創建者為群組成員
	if err := s.groupRepo.AddMember(ctx, group.ID, ownerID); err != nil {
		return nil, err
	}

	return group, nil
}

func (s *groupChatService) UpdateGroup(ctx context.Context, group *entities.Group) error {
	group.UpdatedAt = time.Now()
	return s.groupRepo.Update(ctx, group)
}

func (s *groupChatService) DeleteGroup(ctx context.Context, groupID uint) error {
	return s.groupRepo.Delete(ctx, groupID)
}

func (s *groupChatService) GetGroup(ctx context.Context, groupID uint) (*entities.Group, error) {
	return s.groupRepo.FindByID(ctx, groupID)
}

func (s *groupChatService) GetUserGroups(ctx context.Context, userID uint) ([]*entities.Group, error) {
	// 獲取用戶創建的群組
	createdGroups, err := s.groupRepo.FindByOwnerID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 獲取用戶加入的群組
	joinedGroups, err := s.groupRepo.FindJoinedGroups(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 合併兩個群組列表，並去重
	groupMap := make(map[uint]*entities.Group)
	for _, group := range createdGroups {
		groupMap[group.ID] = group
	}
	for _, group := range joinedGroups {
		groupMap[group.ID] = group
	}

	// 轉換為切片
	var allGroups []*entities.Group
	for _, group := range groupMap {
		allGroups = append(allGroups, group)
	}

	return allGroups, nil
}

func (s *groupChatService) GetGroupsByType(ctx context.Context, groupType int) ([]*entities.Group, error) {
	return s.groupRepo.FindByType(ctx, groupType)
}

func (s *groupChatService) GetGroupsBySize(ctx context.Context, size int) ([]*entities.Group, error) {
	return s.groupRepo.FindBySize(ctx, size)
}

func (s *groupChatService) AddMember(ctx context.Context, groupID, userID uint) error {
	group, err := s.groupRepo.FindByID(ctx, groupID)
	if err != nil {
		return err
	}

	// 檢查群組是否已滿
	members, err := s.groupRepo.GetMembers(ctx, groupID)
	if err != nil {
		return err
	}

	var maxSize int
	switch group.Size {
	case 0:
		maxSize = 50
	case 1:
		maxSize = 200
	case 2:
		maxSize = 500
	default:
		maxSize = 50
	}

	if len(members) >= maxSize {
		return errors.New("group is full")
	}

	return s.groupRepo.AddMember(ctx, groupID, userID)
}

func (s *groupChatService) RemoveMember(ctx context.Context, groupID, userID uint) error {
	return s.groupRepo.RemoveMember(ctx, groupID, userID)
}

func (s *groupChatService) GetGroupMembers(ctx context.Context, groupID uint) ([]uint, error) {
	return s.groupRepo.GetMembers(ctx, groupID)
}

func (s *groupChatService) IsGroupMember(ctx context.Context, groupID, userID uint) (bool, error) {
	return s.groupRepo.IsMember(ctx, groupID, userID)
}

func (s *groupChatService) SendGroupMessage(ctx context.Context, message *entities.Message) error {
	// 檢查發送者是否為群組成員
	isMember, err := s.groupRepo.IsMember(ctx, uint(message.TargetId), uint(message.UserId))
	if err != nil {
		return err
	}
	if !isMember {
		return errors.New("sender is not a member of the group")
	}

	message.CreatedAt = time.Now()
	return s.messageRepo.SaveMessage(ctx, message)
}

func (s *groupChatService) GetGroupHistory(ctx context.Context, groupID uint) ([]*entities.Message, error) {
	return s.messageRepo.FindMessagesByRoomID(ctx, groupID)
}
