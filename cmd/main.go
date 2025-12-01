package main

import (
	"demo/go-server/internal/auth"
	"demo/go-server/internal/hello"
	"fmt"
	"net/http"
)

func main() {
	const PORT = "8081"
	address := fmt.Sprintf(":%s", PORT)

	// conf := configs.LoadConfig()
	router := http.NewServeMux()
	hello.NewHelloHandler(router)
	auth.NewAuthHandler(router)

	server := http.Server{
		Addr:    address,
		Handler: router,
	}

	fmt.Printf("Server is listening on port %s\n", PORT)
	server.ListenAndServe()
}
