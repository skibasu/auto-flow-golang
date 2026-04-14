package appMiddleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHasRequiredRole(t *testing.T) {
	tests := []struct {
		name          string
		userRoles     []string
		requiredRoles []string
		expected      bool
	}{
		{
			name:          "user has one of required roles",
			userRoles:     []string{"USER"},
			requiredRoles: []string{"ADMIN", "USER"},
			expected:      true,
		},
		{
			name:          "user has exact required role",
			userRoles:     []string{"ADMIN"},
			requiredRoles: []string{"ADMIN"},
			expected:      true,
		},
		{
			name:          "user has no required role",
			userRoles:     []string{"USER"},
			requiredRoles: []string{"ADMIN"},
			expected:      false,
		},
		{
			name:          "user has multiple roles and one matches",
			userRoles:     []string{"USER", "MANAGER"},
			requiredRoles: []string{"ADMIN", "MANAGER"},
			expected:      true,
		},
		{
			name:          "empty user roles",
			userRoles:     []string{},
			requiredRoles: []string{"ADMIN"},
			expected:      false,
		},
		{
			name:          "empty required roles",
			userRoles:     []string{"ADMIN"},
			requiredRoles: []string{},
			expected:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := hasRequiredRole(tt.userRoles, tt.requiredRoles)
			if got != tt.expected {
				t.Errorf("hasRequiredRole(%v, %v) = %v, want %v",
					tt.userRoles, tt.requiredRoles, got, tt.expected)
			}
		})
	}
}

func TestRequireRole(t *testing.T) {
	mw := newTestMiddleware()

	tests := []struct {
		name           string
		contextUser    *UserContext // nil means no UserContext set
		requiredRoles  []string
		expectedStatus int
		expectNext     bool
	}{
		{
			name:           "no user context returns 401",
			contextUser:    nil,
			requiredRoles:  []string{"ADMIN"},
			expectedStatus: http.StatusUnauthorized,
			expectNext:     false,
		},
		{
			name:           "user has required role",
			contextUser:    &UserContext{Id: "u1", Roles: []string{"ADMIN"}},
			requiredRoles:  []string{"ADMIN"},
			expectedStatus: http.StatusOK,
			expectNext:     true,
		},
		{
			name:           "user missing required role returns 403",
			contextUser:    &UserContext{Id: "u2", Roles: []string{"USER"}},
			requiredRoles:  []string{"ADMIN"},
			expectedStatus: http.StatusForbidden,
			expectNext:     false,
		},
		{
			name:           "user has one of multiple required roles",
			contextUser:    &UserContext{Id: "u3", Roles: []string{"MANAGER"}},
			requiredRoles:  []string{"ADMIN", "MANAGER"},
			expectedStatus: http.StatusOK,
			expectNext:     true,
		},
		{
			name:           "user has extra roles but includes required",
			contextUser:    &UserContext{Id: "u4", Roles: []string{"USER", "ADMIN"}},
			requiredRoles:  []string{"ADMIN"},
			expectedStatus: http.StatusOK,
			expectNext:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			called := false
			req := httptest.NewRequest(http.MethodGet, "/", nil)

			if tt.contextUser != nil {
				ctx := context.WithValue(req.Context(), UserContextKey, *tt.contextUser)
				req = req.WithContext(ctx)
			}

			rr := httptest.NewRecorder()
			mw.RequireRole(tt.requiredRoles)(nextHandler(&called)).ServeHTTP(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rr.Code)
			}
			if called != tt.expectNext {
				t.Errorf("expected next called=%v, got %v", tt.expectNext, called)
			}
		})
	}
}
