package services

import (
	"github.com/skibasu/auto-flow-api/internal/dto"
	"github.com/skibasu/auto-flow-api/internal/models"
	"github.com/skibasu/auto-flow-api/internal/repository"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetMe(id string) (*models.User, error) {
	return s.repo.GetMe(id)
}

func (s *UserService) GetUsers(filters dto.UsersFilterRequest) (*[]models.User, error) {
	return s.repo.GetUsers(filters)
}

func (s *UserService) CreateUser(user dto.UserRequest) (*models.User, error) {
	return s.repo.CreateUser(user)
}

func (s *UserService) DeleteUser(id string) error {
	return s.repo.DeleteUser(id)
}
