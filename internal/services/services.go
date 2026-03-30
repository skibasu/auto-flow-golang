package services

import "github.com/skibasu/auto-flow-api/internal/repository"

type Service struct {
	repo *repository.UserRepository
}

func NewService(r *repository.UserRepository) *Service {

	return &Service{repo: r}
}
