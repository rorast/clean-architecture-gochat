package user

import (
	"clean-architecture-gochat/internal/domain/entities"
	"clean-architecture-gochat/internal/domain/repositories"
	"clean-architecture-gochat/pkg/utils"
	"context"
	"errors"
	"fmt"
	"log"
)

type Service interface {
	CreateUser(ctx context.Context, user *entities.User) error
	GetUserList(ctx context.Context) ([]*entities.User, error)
	DeleteUser(ctx context.Context, id uint) error
	UpdateUser(ctx context.Context, user *entities.User) error
	FindUserByNameAndPwd(ctx context.Context, name, password string) (*entities.User, error)
	SearchFriend(ctx context.Context, userId uint) ([]entities.User, error)
	AddFriend(ctx context.Context, userId uint, targetId uint) error
	FindUserByID(ctx context.Context, id uint) (*entities.User, error)
	UpdateAvatar(ctx context.Context, userId uint, avatarPath string) error
}

type service struct {
	userRepo    repositories.UserRepository
	contactRepo repositories.ContactRepository
}

func NewService(userRepo repositories.UserRepository, contactRepo repositories.ContactRepository) Service {
	return &service{userRepo: userRepo, contactRepo: contactRepo}
}

func (s *service) CreateUser(ctx context.Context, user *entities.User) error {
	return s.userRepo.Create(ctx, user)
}

func (s *service) GetUserList(ctx context.Context) ([]*entities.User, error) {
	return s.userRepo.GetAll(ctx)
}

func (s *service) DeleteUser(ctx context.Context, id uint) error {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return errors.New("user not found")
	}
	return s.userRepo.Delete(ctx, user.ID)
}

func (s *service) UpdateUser(ctx context.Context, user *entities.User) error {
	_, err := s.userRepo.FindByID(ctx, user.ID)
	if err != nil {
		return errors.New("user not found")
	}
	return s.userRepo.Update(ctx, user)
}

func (s *service) FindUserByNameAndPwd(ctx context.Context, name, password string) (*entities.User, error) {
	user, err := s.userRepo.FindByName(ctx, name)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// 從 user 結構中獲取鹽值
	salt := user.Salt
	// 使用相同的加密方法
	encryptedPassword := utils.MakePassword(password, salt)

	// 比較加密後的密碼
	if encryptedPassword != user.Password {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

func (s *service) SearchFriend(ctx context.Context, userId uint) ([]entities.User, error) {
	contacts, err := s.contactRepo.FindFriendContacts(ctx, userId)
	if err != nil {
		return nil, err
	}

	targetIds := make([]uint, 0, len(contacts))
	for _, c := range contacts {
		targetIds = append(targetIds, c.TargetId)
	}

	return s.userRepo.FindUserByIds(ctx, targetIds)
}

// 加入好友
func (s *service) AddFriend(ctx context.Context, userId uint, targetId uint) error {
	log.Printf("AddFriend 被呼叫: userId=%d, targetId=%d", userId, targetId)

	// 1. 檢查是否為自己加自己
	if targetId == userId {
		log.Println("❌ 不能加自己為好友")
		return errors.New("cannot add yourself as a friend")
	}

	// 2. 查找目標使用者
	targetUser, err := s.userRepo.FindByID(ctx, targetId)
	if err != nil {
		log.Printf("❌ 查找 target user 發生錯誤: %v", err)
		return fmt.Errorf("failed to find target user: %w", err)
	}
	if targetUser == nil {
		log.Printf("❌ 找不到 target user，ID: %d", targetId)
		return errors.New("target user not found")
	}

	// 3. 建立雙向好友關係
	// 3.1 建立當前用戶的好友關係
	err = s.contactRepo.CreateContact(ctx, userId, targetUser.ID, 1)
	if err != nil {
		log.Printf("❌ 建立好友關係失敗: %v", err)
		return errors.New("failed to add friend")
	}

	// 3.2 建立對方的好友關係
	err = s.contactRepo.CreateContact(ctx, targetUser.ID, userId, 1)
	if err != nil {
		log.Printf("❌ 建立對方好友關係失敗: %v", err)
		// 如果建立對方關係失敗，刪除已建立的關係
		s.contactRepo.DeleteContact(ctx, userId, targetUser.ID)
		return errors.New("failed to add friend")
	}

	log.Printf("✅ 成功新增雙向好友關係: %d ↔ %d", userId, targetUser.ID)
	return nil
}

func (s *service) FindUserByID(ctx context.Context, id uint) (*entities.User, error) {
	return s.userRepo.FindByID(ctx, id)
}

// UpdateAvatar 更新用戶頭像
func (s *service) UpdateAvatar(ctx context.Context, userId uint, avatarPath string) error {
	user, err := s.userRepo.FindByID(ctx, userId)
	if err != nil {
		return err
	}

	user.Avatar = avatarPath
	return s.userRepo.Update(ctx, user)
}
