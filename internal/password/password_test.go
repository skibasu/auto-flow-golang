package password

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestHashPassword(t *testing.T) {
	tests := []struct {
		name            string
		input           string
		expectEmptyHash bool
		expectErr       bool
	}{
		{
			name:            "returns hash for valid password",
			input:           "Test1234%",
			expectEmptyHash: false,
			expectErr:       false,
		},
		{
			name:            "returns empty string for empty password",
			input:           "",
			expectEmptyHash: true,
			expectErr:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := HashPassword(tt.input)

			if (err != nil) != tt.expectErr {
				t.Fatalf("expected error: %v, got: %v", tt.expectErr, err)
			}

			if tt.expectEmptyHash {
				if result != "" {
					t.Fatalf("expected empty hash, got %q", result)
				}
				return
			}

			if result == "" {
				t.Fatal("expected non-empty hash, got empty string")
			}

			if compareErr := bcrypt.CompareHashAndPassword([]byte(result), []byte(tt.input)); compareErr != nil {
				t.Fatalf("hash does not match original password: %v", compareErr)
			}
		})
	}
}
