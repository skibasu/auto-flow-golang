package appMiddleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/skibasu/auto-flow-api/internal/appErrors"
	"github.com/skibasu/auto-flow-api/internal/config"
	jwtpkg "github.com/skibasu/auto-flow-api/internal/jwt"
)

const testSecret = "test-secret-key"

func newTestMiddleware() *AppMiddleware {
	return NewAppMiddleware(config.Config{Secret: testSecret})
}

func newValidToken(t *testing.T, userID string, roles []string) string {
	t.Helper()
	token, err := jwtpkg.GenerateToken("access", userID, testSecret, roles, time.Minute)
	if err != nil {
		t.Fatalf("failed to generate test token: %v", err)
	}
	return token
}

func decodeErrorResponse(t *testing.T, body []byte) appErrors.ErrorResponse {
	t.Helper()
	var resp appErrors.ErrorResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		t.Fatalf("failed to decode error response: %v", err)
	}
	return resp
}

// nextHandler is a simple handler that records it was called
func nextHandler(called *bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		*called = true
		w.WriteHeader(http.StatusOK)
	})
}

func TestAuthMiddleware(t *testing.T) {
	mw := newTestMiddleware()

	tests := []struct {
		name           string
		authHeader     string
		expectedStatus int
		expectNext     bool
	}{
		{
			name:           "missing Authorization header",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
			expectNext:     false,
		},
		{
			name:           "invalid format - no Bearer prefix",
			authHeader:     "Token some-token-here",
			expectedStatus: http.StatusUnauthorized,
			expectNext:     false,
		},
		{
			name:           "invalid format - only one part",
			authHeader:     "Bearer",
			expectedStatus: http.StatusUnauthorized,
			expectNext:     false,
		},
		{
			name:           "invalid token value",
			authHeader:     "Bearer invalid.token.value",
			expectedStatus: http.StatusUnauthorized,
			expectNext:     false,
		},
		{
			name: "expired token",
			authHeader: "Bearer " + func() string {
				tok, _ := jwtpkg.GenerateToken("access", "user-1", testSecret, []string{"USER"}, -time.Minute)
				return tok
			}(),
			expectedStatus: http.StatusUnauthorized,
			expectNext:     false,
		},
		{
			name: "valid token calls next handler",
			authHeader: "Bearer " + func() string {
				tok, _ := jwtpkg.GenerateToken("access", "user-abc", testSecret, []string{"ADMIN"}, time.Minute)
				return tok
			}(),
			expectedStatus: http.StatusOK,
			expectNext:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			called := false
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			rr := httptest.NewRecorder()

			mw.AuthMiddleware(nextHandler(&called)).ServeHTTP(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rr.Code)
			}
			if called != tt.expectNext {
				t.Errorf("expected next called=%v, got %v", tt.expectNext, called)
			}
		})
	}
}

func TestAuthMiddleware_SetsUserContext(t *testing.T) {
	mw := newTestMiddleware()
	userID := "user-xyz"
	roles := []string{"ADMIN", "USER"}
	token := newValidToken(t, userID, roles)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	var gotUser UserContext
	captureHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotUser = r.Context().Value(UserContextKey).(UserContext)
		w.WriteHeader(http.StatusOK)
	})

	mw.AuthMiddleware(captureHandler).ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	if gotUser.Id != userID {
		t.Errorf("expected user ID %q, got %q", userID, gotUser.Id)
	}
	if len(gotUser.Roles) != len(roles) {
		t.Errorf("expected roles %v, got %v", roles, gotUser.Roles)
	}
}
