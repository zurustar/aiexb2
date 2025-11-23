// backend/internal/handler/middleware_test.go
package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/your-org/esms/internal/domain"
	"github.com/your-org/esms/internal/handler"
	"github.com/your-org/esms/internal/service"
)

func TestMiddleware_CORS(t *testing.T) {
	mockAuth := new(MockAuthService)
	mw := handler.NewMiddleware(mockAuth)

	tests := []struct {
		name           string
		origin         string
		method         string
		expectedCode   int
		expectedOrigin string
	}{
		{
			name:           "Allowed Origin",
			origin:         "http://localhost:3000",
			method:         "GET",
			expectedCode:   http.StatusOK,
			expectedOrigin: "http://localhost:3000",
		},
		{
			name:           "Disallowed Origin",
			origin:         "http://evil.com",
			method:         "GET",
			expectedCode:   http.StatusOK,
			expectedOrigin: "",
		},
		{
			name:           "Preflight Request",
			origin:         "http://localhost:3000",
			method:         "OPTIONS",
			expectedCode:   http.StatusOK,
			expectedOrigin: "http://localhost:3000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := mw.CORS(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}))

			req := httptest.NewRequest(tt.method, "/api/test", nil)
			req.Header.Set("Origin", tt.origin)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
			assert.Equal(t, tt.expectedOrigin, w.Header().Get("Access-Control-Allow-Origin"))

			if tt.expectedOrigin != "" {
				assert.Equal(t, "true", w.Header().Get("Access-Control-Allow-Credentials"))
			}
		})
	}
}

func TestMiddleware_Authentication(t *testing.T) {
	validSession := &service.Session{
		UserID:      uuid.New(),
		Email:       "test@example.com",
		Role:        domain.RoleGeneral,
		AccessToken: "valid-token",
		ExpiresAt:   time.Now().Add(1 * time.Hour),
	}

	tests := []struct {
		name         string
		authHeader   string
		setupMock    func(*MockAuthService)
		expectedCode int
		expectCall   bool
	}{
		{
			name:       "Success",
			authHeader: "Bearer valid-session-id",
			setupMock: func(m *MockAuthService) {
				m.On("GetSession", "valid-session-id").Return(validSession, nil)
			},
			expectedCode: http.StatusOK,
			expectCall:   true,
		},
		{
			name:         "Missing Header",
			authHeader:   "",
			setupMock:    func(m *MockAuthService) {},
			expectedCode: http.StatusUnauthorized,
			expectCall:   false,
		},
		{
			name:         "Invalid Header Format",
			authHeader:   "InvalidFormat",
			setupMock:    func(m *MockAuthService) {},
			expectedCode: http.StatusUnauthorized,
			expectCall:   false,
		},
		{
			name:       "Session Not Found",
			authHeader: "Bearer invalid-session-id",
			setupMock: func(m *MockAuthService) {
				m.On("GetSession", "invalid-session-id").Return(nil, service.ErrSessionNotFound)
			},
			expectedCode: http.StatusUnauthorized,
			expectCall:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAuth := new(MockAuthService)
			tt.setupMock(mockAuth)
			mw := handler.NewMiddleware(mockAuth)

			called := false
			h := mw.Authentication(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				called = true
				w.WriteHeader(http.StatusOK)
			}))

			req := httptest.NewRequest("GET", "/api/test", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			w := httptest.NewRecorder()

			h.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
			assert.Equal(t, tt.expectCall, called)
			mockAuth.AssertExpectations(t)
		})
	}
}

func TestMiddleware_RateLimit(t *testing.T) {
	mockAuth := new(MockAuthService)
	mw := handler.NewMiddleware(mockAuth)

	// レート制限のテストはタイミングに依存するため、
	// 簡易的な確認にとどめるか、RateLimiterのインターフェースをモック化する必要がある
	// ここでは、通常のリクエストが通ることを確認する

	handler := mw.RateLimit(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/api/test", nil)
	req.RemoteAddr = "192.168.1.1:12345"
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestMiddleware_CSRF(t *testing.T) {
	mockAuth := new(MockAuthService)
	mw := handler.NewMiddleware(mockAuth)

	tests := []struct {
		name         string
		method       string
		headerToken  string
		expectedCode int
	}{
		{
			name:         "GET Request (Skip Check)",
			method:       "GET",
			headerToken:  "",
			expectedCode: http.StatusOK,
		},
		{
			name:         "POST Request with Valid Token",
			method:       "POST",
			headerToken:  "valid-token",
			expectedCode: http.StatusOK,
		},
		{
			name:         "POST Request Missing Token",
			method:       "POST",
			headerToken:  "",
			expectedCode: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := mw.CSRF(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}))

			req := httptest.NewRequest(tt.method, "/api/test", nil)
			if tt.headerToken != "" {
				req.Header.Set("X-CSRF-Token", tt.headerToken)
			}
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
		})
	}
}
