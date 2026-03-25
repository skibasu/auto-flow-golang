package repository

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/skibasu/auto-flow-api/internal/dto"
	"github.com/skibasu/auto-flow-api/internal/models"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetAuthDataByEmail(email string) (*models.UserAuth, error) {
	query := GET_AUTH_BY_EMAIL

	var user models.UserAuth
	var rolesJSON []byte

	err := r.db.QueryRow(context.Background(), query, email).
		Scan(&user.Id, &user.Password, &rolesJSON)

	if err != nil {

		return nil, err
	}

	err = json.Unmarshal(rolesJSON, &user.Roles)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetMe(id string) (*models.User, error) {
	query := GET_ME

	var user models.User
	var rolesJSON []byte

	err := r.db.QueryRow(context.Background(), query, id).
		Scan(&user.Id, &user.Email, &user.FirstName, &user.LastName, &user.PhoneNumber, &user.Avatar, &rolesJSON)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	err = json.Unmarshal(rolesJSON, &user.Roles)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetUsers(filter dto.UsersFilterRequest) (*[]models.User, error) {
	query := GET_USERS_PART_1
	conditions, args := getUserFilters(filter)

	// WHERE
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += GET_USERS_PART_2

	rows, err := r.db.Query(context.Background(), query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var user models.User
		var roles []string

		err := rows.Scan(
			&user.Id,
			&user.Email,
			&user.FirstName,
			&user.LastName,
			&user.PhoneNumber,
			&user.Avatar,
			&roles,
		)
		if err != nil {
			return nil, err
		}

		user.Roles = roles
		users = append(users, user)
	}

	return &users, nil
}

func (r *UserRepository) CreateUser(user dto.UserRequest) (*models.User, error) {
	var createdUser models.User
	query := CREATE_USER

	err := r.db.QueryRow(
		context.Background(),
		query,
		user.Email,
		user.FirstName,
		user.LastName,
		user.PhoneNumber,
		user.Password,
		user.Roles,
	).Scan(
		&createdUser.Id,
		&createdUser.Email,
		&createdUser.FirstName,
		&createdUser.LastName,
		&createdUser.PhoneNumber,
		&createdUser.Avatar,
	)

	if err != nil {
		return nil, err
	}

	createdUser.Roles = user.Roles

	return &createdUser, nil
}

func (r *UserRepository) DeleteUser(id string) error {
	query := DELETE_USER
	cmd, err := r.db.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return errors.New("user not found")
	}

	return nil
}
