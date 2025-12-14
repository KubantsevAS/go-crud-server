package main

import (
	"demo/go-server/configs"
	"demo/go-server/internal/auth"
	"demo/go-server/internal/link"
	"demo/go-server/internal/stat"
	"demo/go-server/internal/user"
	"demo/go-server/pkg/db"
	"demo/go-server/pkg/event"
	"demo/go-server/pkg/middleware"
	"fmt"
	"net/http"
)

func App() http.Handler {
	conf := configs.LoadConfig()
	database := db.NewDb(conf)
	router := http.NewServeMux()
	eventBus := event.NewEventBus()

	// Repositories
	linkRepo := link.NewLinkRepository(database)
	userRepo := user.NewUserRepository(database)
	statRepo := stat.NewStatRepository(database)

	// Services
	authService := auth.NewAuthService(userRepo)
	statService := stat.NewStatService(&stat.StatServiceDeps{
		StatRepository: statRepo,
		EventBus:       eventBus,
	})

	go statService.AddClick()

	// Handlers
	auth.NewAuthHandler(router, auth.AuthHandlerDeps{
		Config:      conf,
		AuthService: authService,
	})
	link.NewLinkHandler(router, link.LinkHandlerDeps{
		LinkRepository: linkRepo,
		EventBus:       eventBus,
		Config:         conf,
	})
	stat.NewStatHandler(router, stat.StatHandlerDeps{
		StatRepository: statRepo,
		Config:         conf,
	})

	// Middlewares
	stack := middleware.Chain(
		middleware.CORS,
		middleware.Logging,
	)

	return stack(router)
}

func main() {
	const PORT = "8081"
	app := App()
	address := fmt.Sprintf(":%s", PORT)

	server := http.Server{
		Addr:    address,
		Handler: app,
	}

	fmt.Printf("Server is listening on port %s\n", PORT)
	server.ListenAndServe()
}
