package services

import (
	"github.com/skibasu/auto-flow-api/internal/dto"
	"github.com/skibasu/auto-flow-api/internal/models"
	"github.com/skibasu/auto-flow-api/internal/password"
)

func (s *Service) GetMe(id string) (*models.User, error) {
	return s.repo.GetMe(id)
}

func (s *Service) GetUsers(filters dto.UsersFilterRequest) (*[]models.User, error) {
	return s.repo.GetUsers(filters)
}

func (s *Service) CreateUser(user dto.UserRequest) (*models.User, error) {

	if user.Password != "" {
		hashed, err := password.HashPassword(&user.Password)
		if err != nil {
			return nil, err
		}
		user.Password = hashed
	} else {
		initPass := "Admin0Auto@"
		hashed, err := password.HashPassword(&initPass)
		if err != nil {
			return nil, err
		}
		user.Password = hashed
	}
	return s.repo.CreateUser(user)
}

func (s *Service) DeleteUser(id string) error {
	return s.repo.DeleteUser(id)
}
func (s *Service) UpdateUser(id string, user dto.UpdateUserRequest) (*models.User, error) {

	if user.Password != nil && *user.Password != "" {
		hashed, err := password.HashPassword(user.Password)
		if err != nil {
			return nil, err
		}
		user.Password = &hashed
	}

	return s.repo.UpdateUser(id, user)
}
