package auth

import (
	"demo/go-server/internal/user"
	"errors"
)

type AuthService struct {
	UserRepository *user.UserRepository
}

func NewAuthService(userRepository *user.UserRepository) *AuthService {
	return &AuthService{UserRepository: userRepository}
}

func (service *AuthService) Register(email, password, name string) (string, error) {
	existedUser, _ := service.UserRepository.GetByEmail(email)

	if existedUser != nil {
		return "", errors.New(ErrUserExists)
	}

	user := &user.User{
		Email:    email,
		Password: password,
		Name:     name,
	}

	_, err := service.UserRepository.Create(user)

	if err != nil {
		return "", err
	}

	return user.Email, nil
}
