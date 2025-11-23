// backend/internal/handler/resource_handler_test.go
package handler_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/your-org/esms/internal/domain"
	"github.com/your-org/esms/internal/handler"
)

// MockResourceRepository for testing
type MockResourceRepository struct {
	mock.Mock
}

func (m *MockResourceRepository) Create(ctx context.Context, resource *domain.Resource) error {
	args := m.Called(ctx, resource)
	return args.Error(0)
}

func (m *MockResourceRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Resource, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Resource), args.Error(1)
}

func (m *MockResourceRepository) Update(ctx context.Context, resource *domain.Resource) error {
	args := m.Called(ctx, resource)
	return args.Error(0)
}

func (m *MockResourceRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockResourceRepository) FindAvailable(ctx context.Context, startAt, endAt time.Time) ([]*domain.Resource, error) {
	args := m.Called(ctx, startAt, endAt)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Resource), args.Error(1)
}

func TestResourceHandler_ListResources(t *testing.T) {
	mockRepo := new(MockResourceRepository)
	h := handler.NewResourceHandler(mockRepo)

	tests := []struct {
		name         string
		queryParam   string
		expectedCode int
		expectedErr  string
	}{
		{
			name:         "No filter",
			queryParam:   "",
			expectedCode: http.StatusOK,
		},
		{
			name:         "Valid is_active=true",
			queryParam:   "?is_active=true",
			expectedCode: http.StatusOK,
		},
		{
			name:         "Valid is_active=false",
			queryParam:   "?is_active=false",
			expectedCode: http.StatusOK,
		},
		{
			name:         "Invalid is_active value",
			queryParam:   "?is_active=invalid",
			expectedCode: http.StatusBadRequest,
			expectedErr:  "INVALID_QUERY_PARAM",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/api/v1/resources"+tt.queryParam, nil)
			w := httptest.NewRecorder()

			h.ListResources(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
			if tt.expectedErr != "" {
				assert.Contains(t, w.Body.String(), tt.expectedErr)
			}
		})
	}
}

func TestResourceHandler_CreateResource(t *testing.T) {
	tests := []struct {
		name         string
		body         string
		setupMock    func(*MockResourceRepository)
		expectedCode int
		expectedErr  string
	}{
		{
			name: "Success",
			body: `{
				"name": "Meeting Room A",
				"type": "MEETING_ROOM",
				"location": "Building 1, Floor 2",
				"capacity": 10
			}`,
			setupMock: func(m *MockResourceRepository) {
				m.On("Create", mock.Anything, mock.Anything).Return(nil)
			},
			expectedCode: http.StatusCreated,
		},
		{
			name: "Missing name",
			body: `{
				"type": "MEETING_ROOM",
				"location": "Building 1, Floor 2"
			}`,
			setupMock:    func(m *MockResourceRepository) {},
			expectedCode: http.StatusBadRequest,
			expectedErr:  "INVALID_NAME",
		},
		{
			name: "Missing type",
			body: `{
				"name": "Meeting Room A",
				"location": "Building 1, Floor 2"
			}`,
			setupMock:    func(m *MockResourceRepository) {},
			expectedCode: http.StatusBadRequest,
			expectedErr:  "INVALID_TYPE",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockResourceRepository)
			h := handler.NewResourceHandler(mockRepo)

			tt.setupMock(mockRepo)

			req := httptest.NewRequest("POST", "/api/v1/resources", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			h.CreateResource(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
			if tt.expectedErr != "" {
				assert.Contains(t, w.Body.String(), tt.expectedErr)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestResourceHandler_GetResource_NotFound(t *testing.T) {
	mockRepo := new(MockResourceRepository)
	h := handler.NewResourceHandler(mockRepo)

	id := uuid.New()
	mockRepo.On("GetByID", mock.Anything, id).Return(nil, assert.AnError)

	req := httptest.NewRequest("GET", "/api/v1/resources/"+id.String(), nil)
	req = mux.SetURLVars(req, map[string]string{"id": id.String()})
	w := httptest.NewRecorder()

	h.GetResource(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "NOT_FOUND")
}
