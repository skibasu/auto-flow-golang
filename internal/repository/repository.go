package repository

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/skibasu/auto-flow-api/internal/dto"
	"github.com/skibasu/auto-flow-api/internal/models"
)

// RepositoryUser defines the interface for user repository operations
type RepositoryUser interface {
	GetMe(id string) (*models.User, error)
	GetUsers(filters dto.UsersFilterRequest) (*[]models.User, error)
	CreateUser(user dto.UserRequest) (*models.User, error)
	DeleteUser(id string) error
	UpdateUser(id string, user dto.UpdateUserRequest) (*models.User, error)
	GetAuthDataByEmail(email string) (*models.UserAuth, error)
}

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}
