package jwt

import (
	"testing"
	"time"
)

func TestGenerateTokenAndParseToken_Success(t *testing.T) {
	secret := "test-secret"
	tokenType := "access"
	userID := "user-123"
	roles := []string{"ADMIN", "MANAGER"}

	token, err := GenerateToken(tokenType, userID, secret, roles, time.Hour)
	if err != nil {
		t.Fatalf("GenerateToken returned error: %v", err)
	}
	if token == "" {
		t.Fatal("GenerateToken returned empty token")
	}

	claims, err := ParseToken(token, secret)
	if err != nil {
		t.Fatalf("ParseToken returned error: %v", err)
	}

	if claims.Sub != userID {
		t.Errorf("Sub: expected %q, got %q", userID, claims.Sub)
	}
	if claims.Type != tokenType {
		t.Errorf("Type: expected %q, got %q", tokenType, claims.Type)
	}
	if len(claims.Roles) != len(roles) {
		t.Fatalf("Roles length: expected %d, got %d", len(roles), len(claims.Roles))
	}
	for i := range roles {
		if claims.Roles[i] != roles[i] {
			t.Errorf("Roles[%d]: expected %q, got %q", i, roles[i], claims.Roles[i])
		}
	}
	if claims.ExpiresAt == nil {
		t.Fatal("ExpiresAt should not be nil")
	}
}

func TestParseToken_WrongSecret(t *testing.T) {
	token, err := GenerateToken("access", "user-1", "correct-secret", []string{"USER"}, time.Hour)
	if err != nil {
		t.Fatalf("GenerateToken returned error: %v", err)
	}

	_, err = ParseToken(token, "wrong-secret")
	if err == nil {
		t.Fatal("expected error for wrong secret, got nil")
	}
}

func TestParseToken_ExpiredToken(t *testing.T) {
	token, err := GenerateToken("access", "user-1", "test-secret", []string{"USER"}, -time.Minute)
	if err != nil {
		t.Fatalf("GenerateToken returned error: %v", err)
	}

	_, err = ParseToken(token, "test-secret")
	if err == nil {
		t.Fatal("expected error for expired token, got nil")
	}
}

func TestParseToken_InvalidTokenFormat(t *testing.T) {
	_, err := ParseToken("not-a-jwt-token", "test-secret")
	if err == nil {
		t.Fatal("expected error for invalid token format, got nil")
	}
}
