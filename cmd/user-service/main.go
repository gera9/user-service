package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gera9/user-service/internal/middleware"
	httpUser "github.com/gera9/user-service/internal/users/delivery/http"
	"github.com/gera9/user-service/internal/users/repository"
	"github.com/gera9/user-service/internal/users/service"
	"github.com/gera9/user-service/pkg/postgres"
	"github.com/go-chi/chi/v5"
)

func main() {
	ctx := context.Background()
	dbConn, err := postgres.NewPostgresConn(ctx, "postgres://user:password@postgres:5432/user_service")
	if err != nil {
		log.Fatal(err)
	}
	defer dbConn.Close()

	mm := &middleware.MiddlewareManager{}

	usersRepo := repository.NewUsersRepository(dbConn)
	usersService := service.NewUsersService(usersRepo)
	usersHandler := httpUser.NewUsersHandler(usersService, mm)

	r := chi.NewRouter()

	r.Route("/api/v1", func(r chi.Router) {
		r.Mount("/users", usersHandler.Routes())
	})

	err = http.ListenAndServe(":3001", r)
	if err != nil {
		log.Println(err)
	}
}
