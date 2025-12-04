package main

import (
	"demo/go-server/configs"
	"demo/go-server/internal/auth"
	"demo/go-server/internal/hello"
	"demo/go-server/pkg/db"
	"fmt"
	"net/http"
)

func main() {
	const PORT = "8081"
	address := fmt.Sprintf(":%s", PORT)

	conf := configs.LoadConfig()
	_ = db.NewDb(conf)
	router := http.NewServeMux()
	hello.NewHelloHandler(router)
	auth.NewAuthHandler(router, auth.AuthHandlerDeps{
		Config: conf,
	})

	server := http.Server{
		Addr:    address,
		Handler: router,
	}

	fmt.Printf("Server is listening on port %s\n", PORT)
	server.ListenAndServe()
}
