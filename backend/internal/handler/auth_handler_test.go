// backend/internal/handler/auth_handler_test.go
package handler_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/your-org/esms/internal/domain"
	"github.com/your-org/esms/internal/handler"
	"github.com/your-org/esms/internal/service"
)

func TestAuthHandler_Login(t *testing.T) {
	mockAuth := new(MockAuthService)
	h := handler.NewAuthHandler(mockAuth)

	mockAuth.On("GetAuthURL", "test-state").Return("http://auth.example.com?state=test-state")

	req := httptest.NewRequest("GET", "/api/v1/auth/login?state=test-state", nil)
	w := httptest.NewRecorder()

	h.Login(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "http://auth.example.com?state=test-state")
	mockAuth.AssertExpectations(t)
}

func TestAuthHandler_Callback_Success(t *testing.T) {
	mockAuth := new(MockAuthService)
	h := handler.NewAuthHandler(mockAuth)

	session := &service.Session{
		UserID:    uuid.New(),
		Email:     "test@example.com",
		Name:      "Test User",
		Role:      domain.RoleGeneral,
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}

	mockAuth.On("HandleCallback", mock.Anything, "valid-code", "valid-state").Return(session, nil)

	req := httptest.NewRequest("GET", "/api/v1/auth/callback?code=valid-code&state=valid-state", nil)
	w := httptest.NewRecorder()

	h.Callback(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Cookie check
	cookies := w.Result().Cookies()
	assert.NotEmpty(t, cookies)
	var sessionCookie *http.Cookie
	for _, c := range cookies {
		if c.Name == "session_id" {
			sessionCookie = c
			break
		}
	}
	assert.NotNil(t, sessionCookie)
	assert.Equal(t, "session-token", sessionCookie.Value) // Mock value from handler
	assert.True(t, sessionCookie.HttpOnly)
	assert.True(t, sessionCookie.Secure)
	assert.Equal(t, http.SameSiteLaxMode, sessionCookie.SameSite)
}

func TestAuthHandler_Callback_MissingCode(t *testing.T) {
	mockAuth := new(MockAuthService)
	h := handler.NewAuthHandler(mockAuth)

	req := httptest.NewRequest("GET", "/api/v1/auth/callback?state=test", nil)
	w := httptest.NewRecorder()

	h.Callback(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Missing code parameter")
}

func TestAuthHandler_Callback_AuthFailed(t *testing.T) {
	mockAuth := new(MockAuthService)
	h := handler.NewAuthHandler(mockAuth)

	mockAuth.On("HandleCallback", mock.Anything, "invalid-code", "valid-state").Return(nil, errors.New("auth failed"))

	req := httptest.NewRequest("GET", "/api/v1/auth/callback?code=invalid-code&state=valid-state", nil)
	w := httptest.NewRecorder()

	h.Callback(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthHandler_Logout(t *testing.T) {
	mockAuth := new(MockAuthService)
	h := handler.NewAuthHandler(mockAuth)

	userID := uuid.New()
	session := &service.Session{UserID: userID}

	mockAuth.On("Logout", mock.Anything, "session-id", userID).Return(nil)

	req := httptest.NewRequest("POST", "/api/v1/auth/logout", nil)
	// Set context with session
	ctx := context.WithValue(req.Context(), handler.ContextKeySession, session)
	req = req.WithContext(ctx)
	// Set cookie
	req.AddCookie(&http.Cookie{Name: "session_id", Value: "session-id"})

	w := httptest.NewRecorder()

	h.Logout(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Verify cookie cleared
	cookies := w.Result().Cookies()
	var sessionCookie *http.Cookie
	for _, c := range cookies {
		if c.Name == "session_id" {
			sessionCookie = c
			break
		}
	}
	assert.NotNil(t, sessionCookie)
	assert.Equal(t, "", sessionCookie.Value)
	assert.True(t, sessionCookie.MaxAge < 0)
}
