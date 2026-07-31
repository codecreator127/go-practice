package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/johnlin/user-service/database"
	"github.com/johnlin/user-service/internal/user"

	"os"
	"os/signal"
	"syscall"
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

	defer pool.Close()

	repo := user.NewPostgresRepository(pool)

	service := user.NewService(repo)

	handler := user.NewHandler(service)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /users/{id}", handler.GetUser)
	mux.HandleFunc("POST /users", handler.CreateUser)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	go func() {
		log.Println("server running on :8080")

		err := server.ListenAndServe()

		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(
		quit,
		os.Interrupt,
		syscall.SIGTERM,
	)

	<-quit
	log.Println("shutting down server")

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)

	defer cancel()

	err = server.Shutdown(ctx)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("server stopped")
}
