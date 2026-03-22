package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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
	query := `
	SELECT 
		u.id,
		u.password,
		COALESCE(json_agg(ur.role) FILTER (WHERE ur.role IS NOT NULL), '[]') AS roles
	FROM users u
	LEFT JOIN user_roles ur ON ur.user_id = u.id
	WHERE u.email = $1
	GROUP BY u.id, u.password
	`

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
	query := `
	SELECT 
		u.id,
		u.email,
		u.first_name,
		u.last_name,
		u.phone_number,
		u.avatar,
		COALESCE(json_agg(ur.role) FILTER (WHERE ur.role IS NOT NULL), '[]') AS roles
	FROM users u
	LEFT JOIN user_roles ur ON ur.user_id = u.id
	WHERE u.id = $1
	GROUP BY 
	u.id,
	u.email,
	u.first_name,
	u.last_name,
	u.phone_number,
	u.avatar;
	`

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
	query := `
	SELECT 
		u.id,
		u.email,
		u.first_name,
		u.last_name,
		u.phone_number,
		u.avatar,
		COALESCE(array_agg(ur.role) FILTER (WHERE ur.role IS NOT NULL), '{}') AS roles
	FROM users u
	LEFT JOIN user_roles ur ON ur.user_id = u.id
	`

	var conditions []string
	var args []any
	argID := 1

	if filter.Email != "" {
		conditions = append(conditions, fmt.Sprintf("u.email ILIKE $%d", argID))
		args = append(args, "%"+filter.Email+"%")
		argID++
	}

	if filter.FirstName != "" {
		conditions = append(conditions, fmt.Sprintf("u.first_name ILIKE $%d", argID))
		args = append(args, "%"+filter.FirstName+"%")
		argID++
	}
	if filter.PhoneNumber != "" {
		conditions = append(conditions, fmt.Sprintf("u.phone_number ILIKE $%d", argID))
		args = append(args, "%"+filter.PhoneNumber+"%")
		argID++
	}
	if filter.LastName != "" {
		conditions = append(conditions, fmt.Sprintf("u.last_name ILIKE $%d", argID))
		args = append(args, "%"+filter.LastName+"%")
		argID++
	}

	// 🔥 roles filter
	if len(filter.Roles) > 0 {
		conditions = append(conditions, fmt.Sprintf("ur.role = ANY($%d)", argID))
		args = append(args, filter.Roles)
		argID++
	}

	// WHERE
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += `
	GROUP BY 
		u.id,
		u.email,
		u.first_name,
		u.last_name,
		u.phone_number,
		u.avatar
	`

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
