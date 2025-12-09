package main

import (
	"demo/go-server/configs"
	"demo/go-server/internal/auth"
	"demo/go-server/internal/link"
	"demo/go-server/internal/user"
	"demo/go-server/pkg/db"
	"demo/go-server/pkg/middleware"
	"fmt"
	"net/http"
)

func main() {
	const PORT = "8081"
	address := fmt.Sprintf(":%s", PORT)

	conf := configs.LoadConfig()
	database := db.NewDb(conf)
	router := http.NewServeMux()

	// Repositories
	linkRepo := link.NewLinkRepository(database)
	userRepo := user.NewUserRepository(database)

	// Services
	authService := auth.NewAuthService(userRepo)

	// Handlers
	auth.NewAuthHandler(router, auth.AuthHandlerDeps{
		Config:      conf,
		AuthService: authService,
	})
	link.NewLinkHandler(router, link.LinkHandlerDeps{
		LinkRepository: linkRepo,
		Config:         conf,
	})

	// Middlewares
	stack := middleware.Chain(
		middleware.CORS,
		middleware.Logging,
	)

	server := http.Server{
		Addr:    address,
		Handler: stack(router),
	}

	fmt.Printf("Server is listening on port %s\n", PORT)
	server.ListenAndServe()
}
