package repositories

import (
	"clean-architecture-gochat/internal/domain/entities"
	"context"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, user *entities.User) error
	GetAll(ctx context.Context) ([]*entities.User, error)
	FindByID(ctx context.Context, id uint) (*entities.User, error)
	FindByName(ctx context.Context, name string) (*entities.User, error)
	Update(ctx context.Context, user *entities.User) error
	Delete(ctx context.Context, id uint) error
	FindUserByIds(ctx context.Context, ids []uint) ([]entities.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *entities.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) GetAll(ctx context.Context) ([]*entities.User, error) {
	var users []*entities.User
	err := r.db.WithContext(ctx).Find(&users).Error
	return users, err
}

func (r *userRepository) FindByID(ctx context.Context, id uint) (*entities.User, error) {
	var user entities.User
	err := r.db.WithContext(ctx).First(&user, id).Error
	return &user, err
}

func (r *userRepository) FindByName(ctx context.Context, name string) (*entities.User, error) {
	var user entities.User
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&user).Error
	return &user, err
}

func (r *userRepository) Update(ctx context.Context, user *entities.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *userRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&entities.User{}, id).Error
}

func (r *userRepository) FindUserByIds(ctx context.Context, ids []uint) ([]entities.User, error) {
	var users []entities.User
	if err := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
