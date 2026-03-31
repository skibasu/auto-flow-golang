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

func NewRouter() *Router {
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

func (r *Router) InitializeRoutes(handlers *handlers.Handler, mw *appMiddleware.AppMiddleware) {
	r.InitializeAdminRoutes(handlers, mw)
	r.InitializePrivateRoutes(handlers, mw)
	r.InitializePublicRoutes(handlers, mw)
}
func (r *Router) InitializePublicRoutes(handlers *handlers.Handler, mw *appMiddleware.AppMiddleware) {
	r.Group(func(r chi.Router) {

		r.With(appMiddleware.ValidateRequest[dto.Credentials](mw, true)).Post("/auth", handlers.Auth())
		r.Post("/refresh", handlers.RefreshToken())
	})
}

func (r *Router) InitializePrivateRoutes(handlers *handlers.Handler, mw *appMiddleware.AppMiddleware) {

	r.Group(func(r chi.Router) {
		r.Use(mw.AuthMiddleware)
		r.Route("/me", func(r chi.Router) {

			r.Get("/", handlers.GetMe())
		})
	})

	r.Group(func(r chi.Router) {
		r.Use(mw.AuthMiddleware)
		r.Use(mw.RequireRole([]string{"MANAGER", "ADMIN"}))
		r.Route("/repairs", func(r chi.Router) {

			r.Get("/", handlers.GetRepairs)
		})
	})
}

func (r *Router) InitializeAdminRoutes(handlers *handlers.Handler, mw *appMiddleware.AppMiddleware) {
	r.Group(func(r chi.Router) {
		r.Use(mw.AuthMiddleware)
		r.Use(mw.RequireRole([]string{"ADMIN"}))
		r.Route("/users", func(r chi.Router) {
			r.With(appMiddleware.ValidateRequest[dto.UserRequest](mw, false)).Post("/", handlers.CreateUser())
			r.With(appMiddleware.ValidateRequest[dto.UpdateUserRequest](mw, false)).Patch("/{id}", handlers.UpdateUser())
			r.Delete("/{id}", handlers.DeleteUser())
			r.Get("/", handlers.GetUsers())
		})
	})
}
