//go:build integration
// +build integration

package repository

import "testing"

func TestRepositoryGetMe(t *testing.T) {
	repo := setupTestRepository(t)

	t.Run("returns user with roles when user exists", func(t *testing.T) {
		seed := seedTestUser(t, repo, []string{"ADMIN"})

		user, err := repo.GetMe(seed.ID)
		if err != nil {
			t.Fatalf("GetMe returned error: %v", err)
		}

		if user.Id != seed.ID {
			t.Fatalf("expected id %q, got %q", seed.ID, user.Id)
		}
		if user.Email != seed.Email {
			t.Fatalf("expected email %q, got %q", seed.Email, user.Email)
		}
		if user.FirstName != seed.FirstName {
			t.Fatalf("expected first name %q, got %q", seed.FirstName, user.FirstName)
		}
		if user.LastName != seed.LastName {
			t.Fatalf("expected last name %q, got %q", seed.LastName, user.LastName)
		}
		if user.PhoneNumber != seed.PhoneNumber {
			t.Fatalf("expected phone number %q, got %q", seed.PhoneNumber, user.PhoneNumber)
		}
		if len(user.Roles) != 1 || user.Roles[0] != "ADMIN" {
			t.Fatalf("expected roles [ADMIN], got %#v", user.Roles)
		}
	})

	t.Run("returns user not found when record does not exist", func(t *testing.T) {
		_, err := repo.GetMe("00000000-0000-0000-0000-000000000000")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "user not found" {
			t.Fatalf("expected user not found error, got %v", err)
		}
	})
}
