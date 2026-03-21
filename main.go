package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/skibasu/auto-flow-api/internal/appMiddleware"
	"github.com/skibasu/auto-flow-api/internal/db"
	"github.com/skibasu/auto-flow-api/internal/handlers"
	"github.com/skibasu/auto-flow-api/internal/repository"
	"github.com/skibasu/auto-flow-api/internal/services"
)

func main() {
	database, err := db.New()
	userRepo := repository.NewUserRepository(database)
	authService := services.New(userRepo)
	userService := services.NewUserService(userRepo)

	if err != nil {
		panic(err)
	}
	defer database.Close()

	router := chi.NewRouter()
	router.Use(middleware.Recoverer)
	router.Use(middleware.Logger)
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}))
	router.Use(appMiddleware.JSON)

	//Public
	router.Group(func(r chi.Router) {

		r.Post("/auth", handlers.Auth(authService))
	})
	//Privet

	router.Group(func(r chi.Router) {
		r.Use(appMiddleware.AuthMiddleware)
		r.Route("/me", func(r chi.Router) {

			r.Get("/", handlers.GetMe(userService))
		})
	})

	router.Group(func(r chi.Router) {
		r.Use(appMiddleware.AuthMiddleware)
		r.Use(appMiddleware.RequireRole([]string{"MANAGER", "ADMIN"}))
		r.Route("/repairs", func(r chi.Router) {

			r.Get("/", handlers.GetRepairs)
		})
	})
	//Admin
	router.Group(func(r chi.Router) {
		r.Use(appMiddleware.AuthMiddleware)
		r.Use(appMiddleware.RequireRole([]string{"ADMIN"}))
		r.Route("/users", func(r chi.Router) {

			r.Get("/", handlers.GetUsers(userService))
		})
	})

	server := &http.Server{Addr: ":3000", Handler: router}
	error := server.ListenAndServe()

	if error != nil {
		fmt.Println("Failed to listen the server. ", error)
	}

}
