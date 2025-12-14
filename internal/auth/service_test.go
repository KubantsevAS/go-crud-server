package auth_test

import (
	"demo/go-server/internal/auth"
	"demo/go-server/internal/user"
	"testing"
)

type MockUserRepository struct {
}

func (repo *MockUserRepository) Create(u *user.User) (*user.User, error) {
	return &user.User{Email: "a@a.com"}, nil
}

func (repo *MockUserRepository) GetByEmail(email string) (*user.User, error) {
	return nil, nil
}

func TestRegisterSuccess(t *testing.T) {
	const initEmail = "a@a.com"
	authService := auth.NewAuthService(&MockUserRepository{})
	email, err := authService.Register(initEmail, "1", "John Doe")
	if err != nil {
		t.Fatal(err)
	}
	if email != initEmail {
		t.Fatalf("Email %s not equal %s", email, initEmail)
	}
}
