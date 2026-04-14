package handlers

import (
	"github.com/skibasu/auto-flow-api/internal/services"
)

type Handler struct {
	services *services.Service
}

func NewHandler(s *services.Service) *Handler {
	return &Handler{services: s}
}
