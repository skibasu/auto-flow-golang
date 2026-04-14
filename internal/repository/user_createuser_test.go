//go:build integration
// +build integration

package repository

import (
	"context"
	"testing"

	"github.com/skibasu/auto-flow-api/internal/dto"
)

func TestRepositoryCreateUser(t *testing.T) {
	repo := setupTestRepository(t)

	t.Run("creates user successfully with roles", func(t *testing.T) {
		req := dto.UserRequest{
			Email:       "createuser.test@test.local",
			FirstName:   "Tomasz",
			LastName:    "Nowak",
			PhoneNumber: "+48987654321",
			Password:    "hashed-password",
			Roles:       []string{"MANAGER", "CLIENT"},
		}

		user, err := repo.CreateUser(req)
		if err != nil {
			t.Fatalf("CreateUser returned error: %v", err)
		}

		t.Cleanup(func() {
			_, _ = repo.db.Exec(context.Background(), `DELETE FROM user_roles WHERE user_id = $1`, user.Id)
			_, _ = repo.db.Exec(context.Background(), `DELETE FROM users WHERE id = $1`, user.Id)
		})

		if user.Id == "" {
			t.Fatal("created user should have an ID")
		}
		if user.Email != req.Email {
			t.Fatalf("expected email %q, got %q", req.Email, user.Email)
		}
		if user.FirstName != req.FirstName {
			t.Fatalf("expected first name %q, got %q", req.FirstName, user.FirstName)
		}
		if user.LastName != req.LastName {
			t.Fatalf("expected last name %q, got %q", req.LastName, user.LastName)
		}
		if user.PhoneNumber != req.PhoneNumber {
			t.Fatalf("expected phone number %q, got %q", req.PhoneNumber, user.PhoneNumber)
		}

		if len(user.Roles) != 2 || user.Roles[0] != "MANAGER" || user.Roles[1] != "CLIENT" {
			t.Fatalf("expected roles [MANAGER, CLIENT], got %#v", user.Roles)
		}
	})

	t.Run("returns duplicate key error when email already exists", func(t *testing.T) {
		existingUser := seedTestUser(t, repo, []string{"ADMIN"})

		req := dto.UserRequest{
			Email:       existingUser.Email,
			FirstName:   "Different",
			LastName:    "Name",
			PhoneNumber: "+48111111111",
			Password:    "password",
			Roles:       []string{"CLIENT"},
		}

		_, err := repo.CreateUser(req)
		if err == nil {
			t.Fatal("expected error for duplicate email, got nil")
		}
		if err.Error() != "duplicate key" {
			t.Fatalf("expected duplicate key error, got %v", err)
		}
	})

	t.Run("creates user with single role", func(t *testing.T) {
		req := dto.UserRequest{
			Email:       "createuser.single@test.local",
			FirstName:   "Ewa",
			LastName:    "Sikora",
			PhoneNumber: "+48123123123",
			Password:    "hash",
			Roles:       []string{"CLIENT"},
		}

		user, err := repo.CreateUser(req)
		if err != nil {
			t.Fatalf("CreateUser returned error: %v", err)
		}

		t.Cleanup(func() {
			_, _ = repo.db.Exec(context.Background(), `DELETE FROM user_roles WHERE user_id = $1`, user.Id)
			_, _ = repo.db.Exec(context.Background(), `DELETE FROM users WHERE id = $1`, user.Id)
		})

		if len(user.Roles) != 1 || user.Roles[0] != "CLIENT" {
			t.Fatalf("expected roles [CLIENT], got %#v", user.Roles)
		}
	})
}
