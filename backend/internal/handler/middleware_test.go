// backend/internal/handler/middleware_test.go
package handler_test

import (
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

// Mock AuthService
type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) GetSession(sessionID string) (*service.Session, error) {
	args := m.Called(sessionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*service.Session), args.Error(1)
}

func TestMiddleware_CORS(t *testing.T) {
	mockAuth := new(MockAuthService)
	mw := handler.NewMiddleware(mockAuth)

	handler := mw.CORS(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// OPTIONS リクエスト（プリフライト）
	req := httptest.NewRequest("OPTIONS", "/api/test", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
	assert.Contains(t, w.Header().Get("Access-Control-Allow-Methods"), "GET")
	assert.Contains(t, w.Header().Get("Access-Control-Allow-Methods"), "POST")
}

func TestMiddleware_Logging(t *testing.T) {
	mockAuth := new(MockAuthService)
	mw := handler.NewMiddleware(mockAuth)

	called := false
	handler := mw.Logging(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/api/test", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.True(t, called)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestMiddleware_Authentication_Success(t *testing.T) {
	mockAuth := new(MockAuthService)
	mw := handler.NewMiddleware(mockAuth)

	session := &service.Session{
		UserID:      uuid.New(),
		Email:       "test@example.com",
		Role:        domain.RoleGeneral,
		AccessToken: "valid-token",
		ExpiresAt:   time.Now().Add(1 * time.Hour),
	}

	mockAuth.On("GetSession", "session-id").Return(session, nil)

	called := false
	handler := mw.Authentication(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		// セッションがコンテキストに設定されているか確認
		sess, ok := r.Context().Value(handler.ContextKeySession).(*service.Session)
		assert.True(t, ok)
		assert.Equal(t, session.Email, sess.Email)
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/api/test", nil)
	req.Header.Set("Authorization", "Bearer session-id")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.True(t, called)
	assert.Equal(t, http.StatusOK, w.Code)
	mockAuth.AssertExpectations(t)
}

func TestMiddleware_Authentication_MissingHeader(t *testing.T) {
	mockAuth := new(MockAuthService)
	mw := handler.NewMiddleware(mockAuth)

	handler := mw.Authentication(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("Should not be called")
	}))

	req := httptest.NewRequest("GET", "/api/test", nil)
	// Authorization ヘッダーなし
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestMiddleware_Authentication_InvalidSession(t *testing.T) {
	mockAuth := new(MockAuthService)
	mw := handler.NewMiddleware(mockAuth)

	mockAuth.On("GetSession", "invalid-session").Return(nil, service.ErrSessionNotFound)

	handler := mw.Authentication(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("Should not be called")
	}))

	req := httptest.NewRequest("GET", "/api/test", nil)
	req.Header.Set("Authorization", "Bearer invalid-session")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	mockAuth.AssertExpectations(t)
}

func TestMiddleware_RateLimit(t *testing.T) {
	mockAuth := new(MockAuthService)
	mw := handler.NewMiddleware(mockAuth)

	handler := mw.RateLimit(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// 最初のリクエストは成功
	req := httptest.NewRequest("GET", "/api/test", nil)
	req.RemoteAddr = "192.168.1.1:12345"
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// レート制限は設定されているが、通常のリクエストでは制限されない
	// （バースト設定により）
}

func TestMiddleware_CSRF_GET(t *testing.T) {
	mockAuth := new(MockAuthService)
	mw := handler.NewMiddleware(mockAuth)

	called := false
	handler := mw.CSRF(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}))

	// GET リクエストは CSRF チェック不要
	req := httptest.NewRequest("GET", "/api/test", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.True(t, called)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestMiddleware_CSRF_POST_MissingToken(t *testing.T) {
	mockAuth := new(MockAuthService)
	mw := handler.NewMiddleware(mockAuth)

	handler := mw.CSRF(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("Should not be called")
	}))

	// POST リクエストで CSRF トークンなし
	req := httptest.NewRequest("POST", "/api/test", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestMiddleware_CSRF_POST_WithToken(t *testing.T) {
	mockAuth := new(MockAuthService)
	mw := handler.NewMiddleware(mockAuth)

	called := false
	handler := mw.CSRF(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}))

	// POST リクエストで CSRF トークンあり
	req := httptest.NewRequest("POST", "/api/test", nil)
	req.Header.Set("X-CSRF-Token", "valid-token")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.True(t, called)
	assert.Equal(t, http.StatusOK, w.Code)
}
