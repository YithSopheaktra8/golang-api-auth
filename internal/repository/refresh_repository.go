package repository

import (
	"time"

	"github.com/yithsopheaktra/go-auth-api/internal/config"
	"github.com/yithsopheaktra/go-auth-api/internal/model"
)

type RefreshRepository struct{}

func NewRefreshRepository() *RefreshRepository {

	return &RefreshRepository{}

}

func (r *RefreshRepository) Create(
	refreshToken *model.RefreshToken,
) error {

	return config.DB.Create(refreshToken).Error

}

func (r *RefreshRepository) FindByJTI(
	jti string,
) (*model.RefreshToken, error) {

	var refreshToken model.RefreshToken

	err := config.DB.
		Where("jti = ?", jti).
		First(&refreshToken).
		Error

	return &refreshToken, err

}

func (r *RefreshRepository) Revoke(
	jti string,
) error {

	return config.DB.
		Model(&model.RefreshToken{}).
		Where("jti = ?", jti).
		Update("revoked", true).
		Error

}

func (r *RefreshRepository) DeleteExpired() error {

	return config.DB.
		Where("expires_at < ?", time.Now()).
		Delete(&model.RefreshToken{}).
		Error

}
