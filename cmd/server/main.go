package main

import (
	"log"
	"net/http"
	"github.com/johnlin/user-service/internal/user"
)


func main() {
	repo := user.NewMemoryRepository()

	service := user.NewService(repo)

	handler := user.NewHandler(service)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /users/{id}", handler.GetUser)
	mux.HandleFunc("POST /users", handler.CreateUser)

	log.Println("server running on :8080")

	err := http.ListenAndServe(":8080", mux)

	if err != nil {
		log.Fatal(err)
	}
}