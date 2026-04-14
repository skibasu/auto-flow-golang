//go:build integration
// +build integration

package repository

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/skibasu/auto-flow-api/internal/config"
)

type testUserSeed struct {
	ID          string
	Email       string
	FirstName   string
	LastName    string
	PhoneNumber string
	Password    string
	Roles       []string
}

func setupTestRepository(t *testing.T) *Repository {
	t.Helper()

	cfg := config.NewConfig()
	if cfg.DBUrl == "" {
		t.Skip("DATABASE_URL is not set; skipping repository integration tests")
	}

	db, err := pgxpool.New(context.Background(), cfg.DBUrl)
	if err != nil {
		t.Fatalf("failed to create test database pool: %v", err)
	}

	if err := db.Ping(context.Background()); err != nil {
		db.Close()
		t.Fatalf("failed to connect to test database: %v", err)
	}

	t.Cleanup(db.Close)

	return NewRepository(db)
}

func seedTestUser(t *testing.T, repo *Repository, roles []string) testUserSeed {
	t.Helper()

	seed := testUserSeed{
		Email:       fmt.Sprintf("repo-user-%d@test.local", time.Now().UnixNano()),
		FirstName:   "Jan",
		LastName:    "Kowalski",
		PhoneNumber: "+48123456789",
		Password:    "hashed-password",
		Roles:       roles,
	}

	return seedTestUserWithData(t, repo, seed)
}

func seedTestUserWithData(t *testing.T, repo *Repository, seed testUserSeed) testUserSeed {
	t.Helper()

	if seed.Email == "" {
		seed.Email = fmt.Sprintf("repo-user-%d@test.local", time.Now().UnixNano())
	}
	if seed.FirstName == "" {
		seed.FirstName = "Jan"
	}
	if seed.LastName == "" {
		seed.LastName = "Kowalski"
	}
	if seed.PhoneNumber == "" {
		seed.PhoneNumber = "+48123456789"
	}
	if seed.Password == "" {
		seed.Password = "hashed-password"
	}

	err := repo.db.QueryRow(
		context.Background(),
		`INSERT INTO users (email, first_name, last_name, phone_number, password)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id`,
		seed.Email,
		seed.FirstName,
		seed.LastName,
		seed.PhoneNumber,
		seed.Password,
	).Scan(&seed.ID)
	if err != nil {
		t.Fatalf("failed to seed user: %v", err)
	}

	if len(seed.Roles) > 0 {
		_, err = repo.db.Exec(
			context.Background(),
			`INSERT INTO user_roles (user_id, role)
			 SELECT $1, unnest($2::user_role[])`,
			seed.ID,
			seed.Roles,
		)
		if err != nil {
			t.Fatalf("failed to seed user roles: %v", err)
		}
	}

	t.Cleanup(func() {
		_, _ = repo.db.Exec(context.Background(), `DELETE FROM user_roles WHERE user_id = $1`, seed.ID)
		_, _ = repo.db.Exec(context.Background(), `DELETE FROM users WHERE id = $1`, seed.ID)
	})

	return seed
}
