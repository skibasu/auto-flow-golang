//go:build integration
// +build integration

package repository

import (
	"context"
	"testing"

	"github.com/skibasu/auto-flow-api/internal/dto"
)

func strPtr(s string) *string {
	return &s
}

func TestRepositoryUpdateUser(t *testing.T) {
	repo := setupTestRepository(t)

	t.Run("returns error when no fields and no roles to update", func(t *testing.T) {
		seed := seedTestUser(t, repo, []string{"ADMIN"})

		req := dto.UpdateUserRequest{} // All nil/empty

		_, err := repo.UpdateUser(seed.ID, req)
		if err == nil {
			t.Fatal("expected error for no fields to update, got nil")
		}
		if err.Error() != "no fields to update" {
			t.Fatalf("expected no fields to update error, got %v", err)
		}
	})

	t.Run("updates user field successfully", func(t *testing.T) {
		seed := seedTestUser(t, repo, []string{"USER"})

		req := dto.UpdateUserRequest{
			FirstName: strPtr("Updated"),
		}

		user, err := repo.UpdateUser(seed.ID, req)
		if err != nil {
			t.Fatalf("UpdateUser returned error: %v", err)
		}

		if user.FirstName != "Updated" {
			t.Fatalf("expected first name Updated, got %q", user.FirstName)
		}
		// Should keep original last name, email, etc
		if user.Email != seed.Email {
			t.Fatalf("expected email unchanged, got %q", user.Email)
		}
	})

	t.Run("updates multiple user fields", func(t *testing.T) {
		seed := seedTestUser(t, repo, []string{"ADMIN"})

		newEmail := "updated.user@test.local"
		req := dto.UpdateUserRequest{
			Email:     &newEmail,
			FirstName: strPtr("Katarzyna"),
			LastName:  strPtr("Zielinska"),
		}

		user, err := repo.UpdateUser(seed.ID, req)
		if err != nil {
			t.Fatalf("UpdateUser returned error: %v", err)
		}

		if user.Email != newEmail {
			t.Fatalf("expected email %q, got %q", newEmail, user.Email)
		}
		if user.FirstName != "Katarzyna" {
			t.Fatalf("expected first name Katarzyna, got %q", user.FirstName)
		}
		if user.LastName != "Zielinska" {
			t.Fatalf("expected last name Zielinska, got %q", user.LastName)
		}
	})

	t.Run("updates user roles successfully", func(t *testing.T) {
		seed := seedTestUser(t, repo, []string{"USER"})

		req := dto.UpdateUserRequest{
			Roles: []string{"ADMIN", "MANAGER"},
		}

		user, err := repo.UpdateUser(seed.ID, req)
		if err != nil {
			t.Fatalf("UpdateUser returned error: %v", err)
		}

		if len(user.Roles) != 2 || user.Roles[0] != "ADMIN" || user.Roles[1] != "MANAGER" {
			t.Fatalf("expected roles [ADMIN, MANAGER], got %#v", user.Roles)
		}

		// Verify old role is gone
		var oldRoleCount int
		err = repo.db.QueryRow(
			context.Background(),
			`SELECT COUNT(*) FROM user_roles WHERE user_id = $1 AND role = 'USER'`,
			seed.ID,
		).Scan(&oldRoleCount)
		if err != nil {
			t.Fatalf("failed to check old roles: %v", err)
		}
		if oldRoleCount != 0 {
			t.Fatal("expected old role USER to be deleted")
		}
	})

	t.Run("clears roles when empty roles list provided", func(t *testing.T) {
		seed := seedTestUser(t, repo, []string{"ADMIN", "MANAGER"})

		req := dto.UpdateUserRequest{
			Roles: []string{}, // Empty but not nil
		}

		user, err := repo.UpdateUser(seed.ID, req)
		if err != nil {
			t.Fatalf("UpdateUser returned error: %v", err)
		}

		if len(user.Roles) != 0 {
			t.Fatalf("expected empty roles, got %#v", user.Roles)
		}

		// Verify all roles deleted from DB
		var roleCount int
		err = repo.db.QueryRow(
			context.Background(),
			`SELECT COUNT(*) FROM user_roles WHERE user_id = $1`,
			seed.ID,
		).Scan(&roleCount)
		if err != nil {
			t.Fatalf("failed to check roles: %v", err)
		}
		if roleCount != 0 {
			t.Fatalf("expected all roles to be deleted, but found %d", roleCount)
		}
	})

	t.Run("updates fields and roles together", func(t *testing.T) {
		seed := seedTestUser(t, repo, []string{"USER"})

		newPhone := "+99988776655"
		req := dto.UpdateUserRequest{
			PhoneNumber: &newPhone,
			Roles:       []string{"ADMIN"},
		}

		user, err := repo.UpdateUser(seed.ID, req)
		if err != nil {
			t.Fatalf("UpdateUser returned error: %v", err)
		}

		if user.PhoneNumber != newPhone {
			t.Fatalf("expected phone %q, got %q", newPhone, user.PhoneNumber)
		}
		if len(user.Roles) != 1 || user.Roles[0] != "ADMIN" {
			t.Fatalf("expected roles [ADMIN], got %#v", user.Roles)
		}
	})

	t.Run("returns user not found error for non-existent ID", func(t *testing.T) {
		req := dto.UpdateUserRequest{
			FirstName: strPtr("Any"),
		}

		_, err := repo.UpdateUser("00000000-0000-0000-0000-000000000000", req)
		if err == nil {
			t.Fatal("expected error for non-existent user, got nil")
		}
		if err.Error() != "user not found" {
			t.Fatalf("expected user not found error, got %v", err)
		}
	})
}
