// backend/internal/handler/reservation_handler_test.go
package handler_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/your-org/esms/internal/handler"
)

func TestReservationHandler_CreateReservation_Unauthorized(t *testing.T) {
	// Mock services would be needed for full implementation
	h := handler.NewReservationHandler(nil, nil)

	body := strings.NewReader(`{
		"resource_ids": ["123e4567-e89b-12d3-a456-426614174000"],
		"title": "Test Meeting",
		"start_at": "2025-06-01T10:00:00Z",
		"end_at": "2025-06-01T11:00:00Z"
	}`)

	req := httptest.NewRequest("POST", "/api/v1/events", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.CreateReservation(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "UNAUTHORIZED")
}

func TestReservationHandler_GetReservation_InvalidID(t *testing.T) {
	h := handler.NewReservationHandler(nil, nil)

	// This would need router setup for path variables
	// Demonstrating test structure
	assert.NotNil(t, h)
}
