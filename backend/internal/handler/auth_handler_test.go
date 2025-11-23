// backend/internal/handler/auth_handler_test.go
package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/your-org/esms/internal/handler"
)

func TestAuthHandler_Login(t *testing.T) {
	// Note: Requires proper AuthService interface implementation
	t.Skip("Requires interface-based AuthService")

	h := handler.NewAuthHandler(nil)

	req := httptest.NewRequest("GET", "/api/v1/auth/login?state=test-state", nil)
	w := httptest.NewRecorder()

	h.Login(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAuthHandler_Callback_MissingCode(t *testing.T) {
	h := handler.NewAuthHandler(nil)

	req := httptest.NewRequest("GET", "/api/v1/auth/callback?state=test", nil)
	w := httptest.NewRecorder()

	h.Callback(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Missing code parameter")
}
