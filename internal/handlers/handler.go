package handlers

import (
	"github.com/skibasu/auto-flow-api/internal/appMiddleware"
	"github.com/skibasu/auto-flow-api/internal/services"
)

type Handler struct {
	services   *services.Service
	middleware *appMiddleware.AppMiddleware
}

func NewHandler(s *services.Service, mw *appMiddleware.AppMiddleware) *Handler {
	return &Handler{services: s, middleware: mw}
}
