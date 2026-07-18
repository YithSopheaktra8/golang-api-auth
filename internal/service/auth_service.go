package service

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	dto "github.com/yithsopheaktra/go-auth-api/internal/dto/auth"
	"github.com/yithsopheaktra/go-auth-api/internal/model"
	"github.com/yithsopheaktra/go-auth-api/internal/repository"
	"github.com/yithsopheaktra/go-auth-api/internal/utils"
	"gorm.io/gorm"
)

type AuthService struct {
	UserRepository *repository.UserRepository
}

func NewAuthService(
	repo *repository.UserRepository,
) *AuthService {

	return &AuthService{

		UserRepository: repo,
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
) (string, error) {

	user, err := s.UserRepository.FindByEmail(
		request.Email,
	)

	if err != nil {
		return "", err
	}

	if !utils.CheckPassword(
		request.Password,
		user.Password,
	) {
		return "", errors.New("invalid password")
	}

	token, err := utils.GenerateToken(
		user.ID.String(),
		user.Email,
	)

	if err != nil {
		return "", err
	}

	return token, nil
}
