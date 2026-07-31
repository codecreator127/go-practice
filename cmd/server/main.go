package main

import (
	"context"
	"log"
	"net/http"

	"github.com/johnlin/user-service/database"
	"github.com/johnlin/user-service/internal/user"
)

func main() {
	ctx := context.Background()

	dbURL := "postgres://postgres:password@localhost:5433/users?sslmode=disable"

	err := database.RunMigrations(dbURL)

	if err != nil {
		log.Fatal(err)
	}

	pool, err := database.NewPostgresPool(
		ctx, dbURL,
	)

	if err != nil {
		log.Fatal(err)
	}

	repo := user.NewPostgresRepository(pool)

	service := user.NewService(repo)

	handler := user.NewHandler(service)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /users/{id}", handler.GetUser)
	mux.HandleFunc("POST /users", handler.CreateUser)

	log.Println("server running on :8080")

	err = http.ListenAndServe(":8080", mux)

	if err != nil {
		log.Fatal(err)
	}
}
