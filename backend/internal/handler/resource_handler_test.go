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

	req := httptest.NewRequest("GET", "/api/v1/resources", nil)
	w := httptest.NewRecorder()

	h.ListResources(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestResourceHandler_CreateResource(t *testing.T) {
	mockRepo := new(MockResourceRepository)
	h := handler.NewResourceHandler(mockRepo)

	mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil)

	body := strings.NewReader(`{
		"name": "Meeting Room A",
		"type": "MEETING_ROOM",
		"location": "Building 1, Floor 2",
		"capacity": 10
	}`)

	req := httptest.NewRequest("POST", "/api/v1/resources", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.CreateResource(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockRepo.AssertExpectations(t)
}
