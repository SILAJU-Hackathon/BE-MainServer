package repositories

import (
	"errors"

	entity "dinacom-11.0-backend/models/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user *entity.User) error
	FindUserByEmail(email string) (*entity.User, error)
	FindUserByID(id uuid.UUID) (*entity.User, error)
	UpdateUserVerified(email string, verified bool) error
	GetAllUsers() ([]entity.User, error)
	GetUsersByRole(role string) ([]entity.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(user *entity.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindUserByEmail(email string) (*entity.User, error) {
	var user entity.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindUserByID(id uuid.UUID) (*entity.User, error) {
	var user entity.User
	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) UpdateUserVerified(email string, verified bool) error {
	return r.db.Model(&entity.User{}).Where("email = ?", email).Update("verified", verified).Error
}

func (r *userRepository) GetAllUsers() ([]entity.User, error) {
	var users []entity.User
	err := r.db.Find(&users).Error
	return users, err
}

func (r *userRepository) GetUsersByRole(role string) ([]entity.User, error) {
	var users []entity.User
	err := r.db.Where("role = ?", role).Find(&users).Error
	return users, err
}
