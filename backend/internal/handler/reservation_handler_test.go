// backend/internal/handler/reservation_handler_test.go
package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/your-org/esms/internal/domain"
	"github.com/your-org/esms/internal/handler"
	"github.com/your-org/esms/internal/service"
)

func TestReservationHandler_CreateReservation(t *testing.T) {
	mockRes := new(MockReservationService)
	mockApp := new(MockApprovalService)
	h := handler.NewReservationHandler(mockRes, mockApp)

	userID := uuid.New()
	session := &service.Session{UserID: userID}
	resourceID := uuid.New()

	tests := []struct {
		name          string
		body          map[string]interface{}
		setupMock     func()
		expectedCode  int
		expectedError string
	}{
		{
			name: "Success",
			body: map[string]interface{}{
				"resource_ids": []string{resourceID.String()},
				"title":        "Test Meeting",
				"start_at":     "2025-06-01T10:00:00Z",
				"end_at":       "2025-06-01T11:00:00Z",
				"timezone":     "UTC",
			},
			setupMock: func() {
				mockRes.On("CreateReservation", mock.Anything, mock.Anything).Return(&domain.Reservation{ID: uuid.New()}, nil)
			},
			expectedCode: http.StatusCreated,
		},
		{
			name: "Validation Error - Missing Timezone",
			body: map[string]interface{}{
				"resource_ids": []string{resourceID.String()},
				"title":        "Test Meeting",
				"start_at":     "2025-06-01T10:00:00Z",
				"end_at":       "2025-06-01T11:00:00Z",
			},
			setupMock:     func() {},
			expectedCode:  http.StatusBadRequest,
			expectedError: "INVALID_TIMEZONE",
		},
		{
			name: "Validation Error - Invalid Time Range",
			body: map[string]interface{}{
				"resource_ids": []string{resourceID.String()},
				"title":        "Test Meeting",
				"start_at":     "2025-06-01T11:00:00Z",
				"end_at":       "2025-06-01T10:00:00Z",
				"timezone":     "UTC",
			},
			setupMock:     func() {},
			expectedCode:  http.StatusBadRequest,
			expectedError: "INVALID_TIME_RANGE",
		},
		{
			name: "Conflict Error - Resource Not Available",
			body: map[string]interface{}{
				"resource_ids": []string{resourceID.String()},
				"title":        "Test Meeting",
				"start_at":     "2025-06-01T10:00:00Z",
				"end_at":       "2025-06-01T11:00:00Z",
				"timezone":     "UTC",
			},
			setupMock: func() {
				mockRes.On("CreateReservation", mock.Anything, mock.Anything).Return(nil, service.ErrResourceNotAvailable)
			},
			expectedCode:  http.StatusConflict,
			expectedError: "RESOURCE_CONFLICT",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mocks
			mockRes.ExpectedCalls = nil
			mockRes.Calls = nil

			tt.setupMock()

			bodyBytes, _ := json.Marshal(tt.body)
			req := httptest.NewRequest("POST", "/api/v1/events", bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")

			// Set session in context
			ctx := context.WithValue(req.Context(), handler.ContextKeySession, session)
			req = req.WithContext(ctx)

			w := httptest.NewRecorder()
			h.CreateReservation(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
			if tt.expectedError != "" {
				assert.Contains(t, w.Body.String(), tt.expectedError)
			}
		})
	}
}

func TestReservationHandler_CreateReservation_Unauthorized(t *testing.T) {
	mockRes := new(MockReservationService)
	mockApp := new(MockApprovalService)
	h := handler.NewReservationHandler(mockRes, mockApp)

	req := httptest.NewRequest("POST", "/api/v1/events", nil)
	w := httptest.NewRecorder()

	h.CreateReservation(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
