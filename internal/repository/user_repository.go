package repository

import (
	"github.com/yithsopheaktra/go-auth-api/internal/model"

	"github.com/yithsopheaktra/go-auth-api/internal/config"
)

type UserRepository struct{}

func NewUserRepository() *UserRepository {

	return &UserRepository{}

}

func (r *UserRepository) Create(
	user *model.User,
) error {

	return config.DB.Create(user).Error

}

func (r *UserRepository) FindByEmail(
	email string,
) (*model.User, error) {

	var user model.User

	err := config.DB.
		Where("email = ?", email).
		First(&user).
		Error

	return &user, err

}
