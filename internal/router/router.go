package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/skibasu/auto-flow-api/internal/appMiddleware"
	"github.com/skibasu/auto-flow-api/internal/dto"
)

type Router struct {
	*chi.Mux
}

type PublicHandler interface {
	Auth() http.HandlerFunc
	RefreshToken() http.HandlerFunc
}

type PrivateHandler interface {
	GetMe() http.HandlerFunc
	GetRepair() http.HandlerFunc
}

type AdminHandler interface {
	CreateUser() http.HandlerFunc
	UpdateUser() http.HandlerFunc
	DeleteUser() http.HandlerFunc
	GetUsers() http.HandlerFunc
}

type Handler interface {
	PublicHandler
	PrivateHandler
	AdminHandler
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

func (r *Router) InitializeRoutes(h Handler, mw *appMiddleware.AppMiddleware) {
	r.InitializeAdminRoutes(h, mw)
	r.InitializePrivateRoutes(h, mw)
	r.InitializePublicRoutes(h, mw)
}
func (r *Router) InitializePublicRoutes(h PublicHandler, mw *appMiddleware.AppMiddleware) {
	r.Group(func(r chi.Router) {

		r.With(appMiddleware.ValidateRequest[dto.Credentials](mw, true)).Post("/auth", h.Auth())
		r.Post("/refresh", h.RefreshToken())
	})
}

func (r *Router) InitializePrivateRoutes(h PrivateHandler, mw *appMiddleware.AppMiddleware) {

	r.Group(func(r chi.Router) {
		r.Use(mw.AuthMiddleware)
		r.Route("/me", func(r chi.Router) {

			r.Get("/", h.GetMe())
		})
	})

	r.Group(func(r chi.Router) {
		r.Use(mw.AuthMiddleware)
		r.Use(mw.RequireRole([]string{"MANAGER", "ADMIN"}))
		r.Route("/repairs", func(r chi.Router) {

			r.Get("/", h.GetRepair())
		})
	})
}

func (r *Router) InitializeAdminRoutes(h AdminHandler, mw *appMiddleware.AppMiddleware) {
	r.Group(func(r chi.Router) {
		r.Use(mw.AuthMiddleware)
		r.Use(mw.RequireRole([]string{"ADMIN"}))
		r.Route("/users", func(r chi.Router) {
			r.With(appMiddleware.ValidateRequest[dto.UserRequest](mw, false)).Post("/", h.CreateUser())
			r.With(appMiddleware.ValidateRequest[dto.UpdateUserRequest](mw, false)).Patch("/{id}", h.UpdateUser())
			r.Delete("/{id}", h.DeleteUser())
			r.Get("/", h.GetUsers())
		})
	})
}
