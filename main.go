package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/skibasu/auto-flow-api/internal/appMiddleware"
	"github.com/skibasu/auto-flow-api/internal/config"
	"github.com/skibasu/auto-flow-api/internal/db"
	"github.com/skibasu/auto-flow-api/internal/handlers"
	"github.com/skibasu/auto-flow-api/internal/repository"
	"github.com/skibasu/auto-flow-api/internal/services"
)

func main() {
	cfg := config.Load()
	database, err := db.New(cfg.DBUrl)
	if err != nil {
		panic(err)
	}
	userRepo := repository.NewUserRepository(database)
	authService := services.New(userRepo)
	userService := services.NewUserService(userRepo)

	defer database.Close()

	router := chi.NewRouter()
	router.Use(middleware.Recoverer)
	router.Use(middleware.Logger)
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:3000",
		},
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

	server := &http.Server{
		Addr:              ":3000",
		Handler:           router,
		ReadTimeout:       15 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	error := server.ListenAndServe()

	if error != nil {
		fmt.Println("Failed to listen the server. ", error)
	}

}
