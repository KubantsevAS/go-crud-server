package auth_test

import (
	"bytes"
	"demo/go-server/configs"
	"demo/go-server/internal/auth"
	"demo/go-server/internal/user"
	"demo/go-server/pkg/db"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func bootstrap() (*auth.AuthHandler, sqlmock.Sqlmock, error) {
	database, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}
	gormDb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: database,
	}))
	if err != nil {
		return nil, nil, err
	}
	userRepo := user.NewUserRepository(&db.Db{
		DB: gormDb,
	})
	handler := auth.AuthHandler{
		Config: &configs.Config{
			Auth: configs.AuthConfig{
				Secret: "secret",
			},
		},
		AuthService: auth.NewAuthService(userRepo),
	}
	return &handler, mock, nil
}

func TestLoginHandlerSuccess(t *testing.T) {
	handler, mock, err := bootstrap()
	if err != nil {
		t.Fatal(err)
		return
	}

	rows := sqlmock.NewRows([]string{"email", "password"}).
		AddRow("a@a.com", "$2a$10$.DuLxeEK7oFAWYt6pXmdzucWnNUDl5I2h1qP0QavAHz4Ur/bmiLZ.")
	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	data, _ := json.Marshal(&auth.LoginRequest{
		Email:    "a@a.com",
		Password: "1",
	})
	reader := bytes.NewReader(data)
	wr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/auth/login", reader)

	handler.Login()(wr, req)
	if wr.Code != http.StatusCreated {
		t.Errorf("Got %d expected %d", wr.Code, http.StatusCreated)
	}
}

func TestRegisterHandlerSuccess(t *testing.T) {
	handler, mock, err := bootstrap()
	if err != nil {
		t.Fatal(err)
		return
	}

	rows := sqlmock.NewRows([]string{"email", "password", "name"})
	mock.ExpectQuery("SELECT").WillReturnRows(rows)
	mock.ExpectBegin()
	mock.ExpectQuery("INSERT").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	data, _ := json.Marshal(&auth.RegisterRequest{
		Email:    "a@a.com",
		Password: "1",
		Name:     "John Doe",
	})
	reader := bytes.NewReader(data)
	wr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/auth/register", reader)
	handler.Register()(wr, req)
	if wr.Code != http.StatusCreated {
		t.Errorf("Got %d expected %d", wr.Code, http.StatusCreated)
	}
}
