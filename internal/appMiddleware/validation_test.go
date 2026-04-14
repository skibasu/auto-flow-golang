package appMiddleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/skibasu/auto-flow-api/internal/dto"
)

func TestValidateRequest(t *testing.T) {
	mw := newTestMiddleware()

	tests := []struct {
		name           string
		body           string
		protectData    bool
		expectedStatus int
		expectNext     bool
	}{
		{
			name:           "valid credentials body calls next",
			body:           `{"email":"test@test.com","password":"Password1!"}`,
			protectData:    false,
			expectedStatus: http.StatusOK,
			expectNext:     true,
		},
		{
			name:           "empty body returns 400",
			body:           "",
			protectData:    false,
			expectedStatus: http.StatusBadRequest,
			expectNext:     false,
		},
		{
			name:           "invalid JSON syntax returns 400",
			body:           `{"email": bad json`,
			protectData:    false,
			expectedStatus: http.StatusBadRequest,
			expectNext:     false,
		},
		{
			name:           "unknown field returns 400",
			body:           `{"email":"test@test.com","password":"Password1!","unknownField":"x"}`,
			protectData:    false,
			expectedStatus: http.StatusBadRequest,
			expectNext:     false,
		},
		{
			name:           "missing required field returns 400",
			body:           `{"email":"test@test.com"}`,
			protectData:    false,
			expectedStatus: http.StatusBadRequest,
			expectNext:     false,
		},
		{
			name:           "invalid email format returns 400",
			body:           `{"email":"not-an-email","password":"Password1!"}`,
			protectData:    false,
			expectedStatus: http.StatusBadRequest,
			expectNext:     false,
		},
		{
			name:           "weak password returns 400",
			body:           `{"email":"test@test.com","password":"weak"}`,
			protectData:    false,
			expectedStatus: http.StatusBadRequest,
			expectNext:     false,
		},
		{
			name:           "protectData=true with invalid body returns 401",
			body:           `{"email":"bad"}`,
			protectData:    true,
			expectedStatus: http.StatusUnauthorized,
			expectNext:     false,
		},
		{
			name:           "protectData=true with empty body returns 401",
			body:           "",
			protectData:    true,
			expectedStatus: http.StatusUnauthorized,
			expectNext:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			called := false
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			handler := ValidateRequest[dto.Credentials](mw, tt.protectData)(nextHandler(&called))
			handler.ServeHTTP(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d (body: %s)", tt.expectedStatus, rr.Code, rr.Body.String())
			}
			if called != tt.expectNext {
				t.Errorf("expected next called=%v, got %v", tt.expectNext, called)
			}
		})
	}
}

func TestValidateRequest_StoresBodyInContext(t *testing.T) {
	mw := newTestMiddleware()

	body := `{"email":"context@test.com","password":"Password1!"}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	var got dto.Credentials
	captureHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got = GetValidatedBody[dto.Credentials](r)
		w.WriteHeader(http.StatusOK)
	})

	ValidateRequest[dto.Credentials](mw, false)(captureHandler).ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	if got.Email != "context@test.com" {
		t.Errorf("expected email %q, got %q", "context@test.com", got.Email)
	}
}

func TestGetValidatedBody_PanicsWithoutMiddleware(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic when GetValidatedBody called without middleware")
		}
	}()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	GetValidatedBody[dto.Credentials](req)
}

func TestValidateRequest_UserRequest(t *testing.T) {
	mw := newTestMiddleware()

	tests := []struct {
		name           string
		body           string
		expectedStatus int
		expectNext     bool
	}{
		{
			name:           "valid user request",
			body:           `{"email":"user@test.com","firstName":"Jan","lastName":"Kowalski","phoneNumber":"+48123456789","roles":["CLIENT"]}`,
			expectedStatus: http.StatusOK,
			expectNext:     true,
		},
		{
			name:           "invalid role value",
			body:           `{"email":"user@test.com","firstName":"Jan","lastName":"Kowalski","phoneNumber":"+48123456789","roles":["SUPERUSER"]}`,
			expectedStatus: http.StatusBadRequest,
			expectNext:     false,
		},
		{
			name: "invalid phone number",
			body: `{"email":"user@test.com","firstName":"Jan","lastName":"Kowalski","phoneNumber":"not-a-phone","roles":["CLIENT"]}`,

			expectedStatus: http.StatusBadRequest,
			expectNext:     false,
		},
		{
			name:           "missing required roles field",
			body:           `{"email":"user@test.com","firstName":"Jan","lastName":"Kowalski","phoneNumber":"+48123456789"}`,
			expectedStatus: http.StatusBadRequest,
			expectNext:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			called := false
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			ValidateRequest[dto.UserRequest](mw, false)(nextHandler(&called)).ServeHTTP(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d (body: %s)", tt.expectedStatus, rr.Code, rr.Body.String())
			}
			if called != tt.expectNext {
				t.Errorf("expected next called=%v, got %v", tt.expectNext, called)
			}
		})
	}
}
