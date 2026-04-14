//go:build integration
// +build integration

package repository

import (
	"context"
	"testing"
)

func TestRepositoryDeleteUser(t *testing.T) {
	repo := setupTestRepository(t)

	t.Run("deletes user successfully", func(t *testing.T) {
		seed := seedTestUser(t, repo, []string{"ADMIN"})

		err := repo.DeleteUser(seed.ID)
		if err != nil {
			t.Fatalf("DeleteUser returned error: %v", err)
		}

		// Verify user is actually deleted
		deletedUser, err := repo.GetMe(seed.ID)
		if err == nil {
			t.Fatalf("expected user to be deleted, but GetMe returned: %#v", deletedUser)
		}
		if err.Error() != "user not found" {
			t.Fatalf("expected user not found error after delete, got %v", err)
		}
	})

	t.Run("deletes user roles when user is deleted", func(t *testing.T) {
		seed := seedTestUser(t, repo, []string{"ADMIN", "MANAGER"})

		err := repo.DeleteUser(seed.ID)
		if err != nil {
			t.Fatalf("DeleteUser returned error: %v", err)
		}

		// Verify roles are deleted too (cascade should handle this)
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
			t.Fatalf("expected user roles to be deleted, but found %d roles", roleCount)
		}
	})

	t.Run("returns user not found error when user does not exist", func(t *testing.T) {
		err := repo.DeleteUser("00000000-0000-0000-0000-000000000000")
		if err == nil {
			t.Fatal("expected error when deleting non-existent user, got nil")
		}
		if err.Error() != "user not found" {
			t.Fatalf("expected user not found error, got %v", err)
		}
	})
}
