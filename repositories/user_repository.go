package repositories

import (
	"errors"

	entity "dinacom-11.0-backend/models/entity"

	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user *entity.User) error
	FindUserByEmail(email string) (*entity.User, error)
	UpdateUserVerified(email string, verified bool) error
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
			return nil, nil // Not found is not an error in this context, just nil user
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) UpdateUserVerified(email string, verified bool) error {
	return r.db.Model(&entity.User{}).Where("email = ?", email).Update("verified", verified).Error
}
