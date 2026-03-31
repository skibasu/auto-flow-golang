package services

import (
	"github.com/skibasu/auto-flow-api/internal/config"
	"github.com/skibasu/auto-flow-api/internal/repository"
)

type Service struct {
	repo   *repository.Repository
	config config.Config
}

func NewService(r *repository.Repository, config config.Config) *Service {

	return &Service{repo: r, config: config}
}
