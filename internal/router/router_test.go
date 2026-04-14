package router

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/skibasu/auto-flow-api/internal/appMiddleware"
	"github.com/skibasu/auto-flow-api/internal/config"
	"github.com/skibasu/auto-flow-api/internal/jwt"
)

// --- MOCK HANDLER ---
type mockPublicHandler struct{}

func (m *mockPublicHandler) Auth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"ok": true}`))
	}
}

func (m *mockPublicHandler) RefreshToken() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
}

type mockPrivateHandler struct{}

func (m *mockPrivateHandler) GetMe() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"ok": true}`))
	}
}
func (m *mockPrivateHandler) GetRepair() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"ok": true}`))
	}
}

type mockAdminHandler struct{}

func (m *mockAdminHandler) GetUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"ok": true}`))
	}
}
func (m *mockAdminHandler) CreateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"ok": true}`))
	}
}

func (m *mockAdminHandler) UpdateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"ok": true}`))
	}
}
func (m *mockAdminHandler) DeleteUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"ok": true}`))
	}
}

func Test_InitializePublicRoutes(t *testing.T) {
	r := NewRouter()

	cfg := config.Config{
		Secret: "test-secret",
	}
	mw := appMiddleware.NewAppMiddleware(cfg)

	handlers := &mockPublicHandler{}

	r.InitializePublicRoutes(handlers, mw)

	tests := []struct {
		name           string
		method         string
		path           string
		body           string
		expectedStatus int
	}{
		{
			name:           "POST /auth valid",
			method:         http.MethodPost,
			path:           "/auth",
			body:           `{"email":"test@test.com","password":"Test123!"}`,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "POST /auth invalid body",
			method:         http.MethodPost,
			path:           "/auth",
			body:           `{}`,
			expectedStatus: http.StatusUnauthorized, //protectData=true
		},
		{
			name:           "POST /refresh",
			method:         http.MethodPost,
			path:           "/refresh",
			body:           ``,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "GET /auth not allowed",
			method:         http.MethodGet,
			path:           "/auth",
			body:           ``,
			expectedStatus: http.StatusMethodNotAllowed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")

			rec := httptest.NewRecorder()

			r.ServeHTTP(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("expected %d, got %d", tt.expectedStatus, rec.Code)
			}
		})
	}

}

func Test_InitializePrivateRoutes(t *testing.T) {
	r := NewRouter()

	cfg := config.Config{
		Secret: "test-secret",
	}
	mw := appMiddleware.NewAppMiddleware(cfg)

	handlers := &mockPrivateHandler{}

	r.InitializePrivateRoutes(handlers, mw)

	managerToken, err := jwt.GenerateToken("access", "user-1", cfg.Secret, []string{"MANAGER"}, time.Hour)
	if err != nil {
		t.Fatalf("failed to generate manager token: %v", err)
	}

	userToken, err := jwt.GenerateToken("access", "user-2", cfg.Secret, []string{"CLIENT"}, time.Hour)
	if err != nil {
		t.Fatalf("failed to generate user token: %v", err)
	}

	tests := []struct {
		name           string
		method         string
		path           string
		authHeader     string
		expectedStatus int
	}{
		{
			name:           "GET /me unauthorized without token",
			method:         http.MethodGet,
			path:           "/me/",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "GET /me unauthorized invalid auth format",
			method:         http.MethodGet,
			path:           "/me/",
			authHeader:     "Token invalid",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "GET /me authorized",
			method:         http.MethodGet,
			path:           "/me/",
			authHeader:     "Bearer " + managerToken,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "GET /repairs forbidden for user role",
			method:         http.MethodGet,
			path:           "/repairs/",
			authHeader:     "Bearer " + userToken,
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "GET /repairs authorized for manager role",
			method:         http.MethodGet,
			path:           "/repairs/",
			authHeader:     "Bearer " + managerToken,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "POST /me method not allowed",
			method:         http.MethodPost,
			path:           "/me/",
			authHeader:     "Bearer " + managerToken,
			expectedStatus: http.StatusMethodNotAllowed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			rec := httptest.NewRecorder()

			r.ServeHTTP(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("expected %d, got %d", tt.expectedStatus, rec.Code)
			}
		})
	}
}

func Test_InitializeAdminRoutes(t *testing.T) {
	r := NewRouter()

	cfg := config.Config{
		Secret: "test-secret",
	}
	mw := appMiddleware.NewAppMiddleware(cfg)

	handlers := &mockAdminHandler{}

	r.InitializeAdminRoutes(handlers, mw)

	adminToken, err := jwt.GenerateToken("access", "admin-1", cfg.Secret, []string{"ADMIN"}, time.Hour)
	if err != nil {
		t.Fatalf("failed to generate admin token: %v", err)
	}

	managerToken, err := jwt.GenerateToken("access", "manager-1", cfg.Secret, []string{"MANAGER"}, time.Hour)
	if err != nil {
		t.Fatalf("failed to generate manager token: %v", err)
	}

	expiredToken, err := jwt.GenerateToken("access", "admin-2", cfg.Secret, []string{"ADMIN"}, -time.Hour)
	if err != nil {
		t.Fatalf("failed to generate expired token: %v", err)
	}

	validCreateBody := `{"email":"admin.user@test.com","password":"Test123!","firstName":"John","lastName":"Doe","phoneNumber":"+48123456789","roles":["ADMIN"]}`
	validPatchBody := `{"firstName":"Johnny"}`

	tests := []struct {
		name           string
		method         string
		path           string
		body           string
		authHeader     string
		expectedStatus int
	}{
		{
			name:           "GET /users unauthorized without token",
			method:         http.MethodGet,
			path:           "/users/",
			body:           "",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "GET /users forbidden for non admin",
			method:         http.MethodGet,
			path:           "/users/",
			body:           "",
			authHeader:     "Bearer " + managerToken,
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "GET /users authorized for admin",
			method:         http.MethodGet,
			path:           "/users/",
			body:           "",
			authHeader:     "Bearer " + adminToken,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "POST /users valid body",
			method:         http.MethodPost,
			path:           "/users/",
			body:           validCreateBody,
			authHeader:     "Bearer " + adminToken,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "POST /users invalid body",
			method:         http.MethodPost,
			path:           "/users/",
			body:           `{}`,
			authHeader:     "Bearer " + adminToken,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "PATCH /users/{id} valid body",
			method:         http.MethodPatch,
			path:           "/users/123",
			body:           validPatchBody,
			authHeader:     "Bearer " + adminToken,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "DELETE /users/{id} authorized",
			method:         http.MethodDelete,
			path:           "/users/123",
			body:           "",
			authHeader:     "Bearer " + adminToken,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "DELETE /users/{id} unauthorized without token",
			method:         http.MethodDelete,
			path:           "/users/123",
			body:           "",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "DELETE /users/{id} forbidden for non admin",
			method:         http.MethodDelete,
			path:           "/users/123",
			body:           "",
			authHeader:     "Bearer " + managerToken,
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "GET /users expired token",
			method:         http.MethodGet,
			path:           "/users/",
			body:           "",
			authHeader:     "Bearer " + expiredToken,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "PUT /users method not allowed",
			method:         http.MethodPut,
			path:           "/users/",
			body:           "",
			authHeader:     "Bearer " + adminToken,
			expectedStatus: http.StatusMethodNotAllowed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, strings.NewReader(tt.body))
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			if tt.body != "" {
				req.Header.Set("Content-Type", "application/json")
			}

			rec := httptest.NewRecorder()

			r.ServeHTTP(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("expected %d, got %d", tt.expectedStatus, rec.Code)
			}
		})
	}
}
