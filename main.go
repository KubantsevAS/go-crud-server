package main

import (
	"fmt"
	"net/http"
)

func main() {
	const PORT = "8081"
	address := fmt.Sprintf(":%s", PORT)

	router := http.NewServeMux()
	NewHelloHandler(router)

	server := http.Server{
		Addr:    address,
		Handler: router,
	}

	fmt.Printf("Server is listening on port %s\n", PORT)
	server.ListenAndServe()
}
