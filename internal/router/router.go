package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/skibasu/auto-flow-api/internal/appMiddleware"
	"github.com/skibasu/auto-flow-api/internal/dto"
	"github.com/skibasu/auto-flow-api/internal/handlers"
)

type Router struct {
	*chi.Mux
}

func New() *Router {

	return &Router{chi.NewRouter()}
}

func (r *Router) InitializeMiddlewares() {
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{
			"https://localhost:3001",
		},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}))
	r.Use(appMiddleware.JSON)
}

func (r *Router) InitializePublicRoutes() {
	r.Group(func(r chi.Router) {

		r.With(appMiddleware.ValidateRequest[dto.Credentials](true)).Post("/auth", handlers.Auth(authService))
		r.Post("/refresh", handlers.RefreshToken(authService))
	})
}
