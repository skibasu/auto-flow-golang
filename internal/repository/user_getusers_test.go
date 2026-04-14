//go:build integration
// +build integration

package repository

import (
	"testing"

	"github.com/skibasu/auto-flow-api/internal/dto"
)

func TestRepositoryGetUsers(t *testing.T) {
	repo := setupTestRepository(t)

	t.Run("returns users matching email filter", func(t *testing.T) {
		matched := seedTestUserWithData(t, repo, testUserSeed{
			Email:       "getusers.email-match@test.local",
			FirstName:   "Alicja",
			LastName:    "Nowak",
			PhoneNumber: "+48111111111",
			Roles:       []string{"ADMIN"},
		})
		_ = seedTestUserWithData(t, repo, testUserSeed{
			Email:       "other-user@test.local",
			FirstName:   "Piotr",
			LastName:    "Kowal",
			PhoneNumber: "+48222222222",
			Roles:       []string{"CLIENT"},
		})

		users, err := repo.GetUsers(dto.UsersFilterRequest{Email: "email-match"})
		if err != nil {
			t.Fatalf("GetUsers returned error: %v", err)
		}

		if len(*users) != 1 {
			t.Fatalf("expected 1 user, got %d", len(*users))
		}
		if (*users)[0].Email != matched.Email {
			t.Fatalf("expected email %q, got %q", matched.Email, (*users)[0].Email)
		}
	})

	t.Run("returns users matching role filter", func(t *testing.T) {
		adminUser := seedTestUserWithData(t, repo, testUserSeed{
			Email:       "getusers.admin@test.local",
			FirstName:   "Admin",
			LastName:    "User",
			PhoneNumber: "+48333333333",
			Roles:       []string{"ADMIN"},
		})
		_ = seedTestUserWithData(t, repo, testUserSeed{
			Email:       "getusers.manager@test.local",
			FirstName:   "Manager",
			LastName:    "User",
			PhoneNumber: "+48444444444",
			Roles:       []string{"MANAGER"},
		})

		users, err := repo.GetUsers(dto.UsersFilterRequest{Roles: []string{"ADMIN"}})
		if err != nil {
			t.Fatalf("GetUsers returned error: %v", err)
		}

		if len(*users) == 0 {
			t.Fatal("expected at least one admin user, got none")
		}

		found := false
		for _, user := range *users {
			if user.Email == adminUser.Email {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("expected to find admin user %q in results", adminUser.Email)
		}
	})

	t.Run("returns users matching combined filters", func(t *testing.T) {
		matched := seedTestUserWithData(t, repo, testUserSeed{
			Email:       "getusers.combined@test.local",
			FirstName:   "Marta",
			LastName:    "Zielinska",
			PhoneNumber: "+48555555555",
			Roles:       []string{"MANAGER"},
		})
		_ = seedTestUserWithData(t, repo, testUserSeed{
			Email:       "getusers.combined.other@test.local",
			FirstName:   "Marta",
			LastName:    "Inna",
			PhoneNumber: "+48666666666",
			Roles:       []string{"CLIENT"},
		})

		users, err := repo.GetUsers(dto.UsersFilterRequest{
			FirstName: "Marta",
			Roles:     []string{"MANAGER"},
		})
		if err != nil {
			t.Fatalf("GetUsers returned error: %v", err)
		}

		if len(*users) != 1 {
			t.Fatalf("expected 1 user for combined filter, got %d", len(*users))
		}
		if (*users)[0].Email != matched.Email {
			t.Fatalf("expected email %q, got %q", matched.Email, (*users)[0].Email)
		}
		if len((*users)[0].Roles) != 1 || (*users)[0].Roles[0] != "MANAGER" {
			t.Fatalf("expected roles [MANAGER], got %#v", (*users)[0].Roles)
		}
	})
}
