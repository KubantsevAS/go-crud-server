package auth

import (
	"demo/go-server/configs"
	"demo/go-server/pkg/request"
	"demo/go-server/pkg/response"
	"fmt"
	"net/http"
)

type AuthHandlerDeps struct {
	*configs.Config
}

type AuthHandler struct {
	*configs.Config
}

func NewAuthHandler(router *http.ServeMux, deps AuthHandlerDeps) {
	handler := &AuthHandler{
		Config: deps.Config,
	}
	router.HandleFunc("POST /auth/login", handler.Login())
	router.HandleFunc("POST /auth/register", handler.Register())
}

func (handler *AuthHandler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		body, err := request.HandleBody[LoginRequest](&w, req)

		if err != nil {
			fmt.Println(err.Error())
			return
		}

		fmt.Println(body)
		data := LoginResponse{
			Token: handler.Config.Auth.Secret,
		}
		response.WriteResponse(w, data, 201)
	}
}

func (handler *AuthHandler) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		body, err := request.HandleBody[RegisterRequest](&w, req)

		if err != nil {
			fmt.Println(err.Error())
			return
		}

		fmt.Println(body)
		data := RegisterResponse{
			Token: handler.Config.Auth.Secret,
		}
		response.WriteResponse(w, data, 201)
	}
}
