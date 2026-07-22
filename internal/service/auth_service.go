package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	dto "github.com/yithsopheaktra/go-auth-api/internal/dto/auth"
	"github.com/yithsopheaktra/go-auth-api/internal/model"
	"github.com/yithsopheaktra/go-auth-api/internal/repository"
	"github.com/yithsopheaktra/go-auth-api/internal/utils"
	"gorm.io/gorm"
)

type AuthService struct {
	UserRepository         *repository.UserRepository
	RefreshTokenRepository *repository.RefreshRepository
}

func NewAuthService(
	userRepository *repository.UserRepository,
	refreshTokenRepository *repository.RefreshRepository,
) *AuthService {

	return &AuthService{
		UserRepository:         userRepository,
		RefreshTokenRepository: refreshTokenRepository,
	}

}

func (s *AuthService) Register(
	request dto.RegisterRequest,
) error {

	// Check existing email
	existingUser, err := s.UserRepository.FindByEmail(
		request.Email,
	)
	fmt.Printf("FOUND USER: %+v\n", existingUser)
	fmt.Println("FIND ERROR:", err)

	if err == nil && existingUser.ID != uuid.Nil {
		return errors.New("email already exists")
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	hashPassword, err := utils.HashPassword(request.Password)

	if err != nil {
		return err
	}

	user := model.User{

		Name: request.Name,

		Email: request.Email,

		Password: hashPassword,
	}

	return s.UserRepository.Create(&user)
}

func (s *AuthService) Login(
	request dto.LoginRequest,
) (*dto.LoginResponse, error) {

	user, err := s.UserRepository.FindByEmail(request.Email)
	if err != nil {
		return nil, err
	}

	if !utils.CheckPassword(request.Password, user.Password) {
		return nil, errors.New("invalid password")
	}

	accessToken, err := utils.GenerateAccessToken(
		user.ID.String(),
		user.Email,
	)
	if err != nil {
		return nil, err
	}

	refreshToken, jti, err := utils.GenerateRefreshToken(
		user.ID.String(),
	)
	if err != nil {
		return nil, err
	}

	refresh := model.RefreshToken{
		UserID:    user.ID,
		JTI:       jti,
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
	}

	err = s.RefreshTokenRepository.Create(&refresh)
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) Refresh(
	refreshToken string,
) (string, string, error) {

	// 1. Validate refresh JWT

	claims, err := utils.ValidateToken(refreshToken)

	if err != nil {
		return "", "", err
	}

	// 2. Check token type

	tokenType, ok := claims["type"].(string)

	if !ok || tokenType != "refresh" {

		return "", "", errors.New(
			"invalid token type",
		)
	}

	// 3. Get JTI

	jti, ok := claims["jti"].(string)

	if !ok {

		return "", "", errors.New(
			"missing token id",
		)

	}

	// 4. Find refresh token in database

	oldToken, err := s.RefreshTokenRepository.FindByJTI(jti)

	if err != nil {

		return "", "", errors.New(
			"refresh token not found",
		)

	}

	// 5. Check revoked

	if oldToken.Revoked {

		return "", "", errors.New(
			"refresh token already used",
		)

	}

	// 6. Check expiry

	if oldToken.ExpiresAt.Before(time.Now()) {

		return "", "", errors.New(
			"refresh token expired",
		)

	}

	// 7. Revoke old refresh token

	err = s.RefreshTokenRepository.Revoke(jti)

	if err != nil {
		return "", "", err
	}

	// 8. Get user id

	userID, ok := claims["user_id"].(string)

	if !ok {
		return "", "", errors.New(
			"invalid user id",
		)
	}

	// 9. Find user

	user, err := s.UserRepository.FindByID(
		userID,
	)

	if err != nil {

		return "", "", err

	}

	// 10. Create new access token

	newAccessToken, err :=
		utils.GenerateAccessToken(
			user.ID.String(),
			user.Email,
		)

	if err != nil {
		return "", "", err
	}

	// 11. Create new refresh token

	newRefreshToken, newJTI, err :=
		utils.GenerateRefreshToken(
			user.ID.String(),
		)

	if err != nil {
		return "", "", err
	}

	// 12. Save new refresh token

	newRefresh := model.RefreshToken{

		UserID: user.ID,

		JTI: newJTI,

		ExpiresAt: time.Now().
			Add(30 * 24 * time.Hour),
	}

	err = s.RefreshTokenRepository.Create(
		&newRefresh,
	)

	if err != nil {
		return "", "", err
	}

	return newAccessToken,
		newRefreshToken,
		nil
}

func (s *AuthService) Logout(
	refreshToken string,
) error {

	claims, err := utils.ValidateToken(
		refreshToken,
	)

	if err != nil {
		return err
	}

	tokenType, ok := claims["type"].(string)

	if !ok || tokenType != "refresh" {

		return errors.New(
			"invalid token type",
		)

	}

	jti, ok := claims["jti"].(string)

	if !ok {

		return errors.New(
			"missing token id",
		)

	}

	err = s.RefreshTokenRepository.Revoke(
		jti,
	)

	return err

}
