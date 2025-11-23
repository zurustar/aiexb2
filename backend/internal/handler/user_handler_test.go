// backend/internal/handler/user_handler_test.go
package handler_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/your-org/esms/internal/domain"
	"github.com/your-org/esms/internal/handler"
	"github.com/your-org/esms/internal/service"
)

// MockUserRepository for testing
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestUserHandler_ListUsers(t *testing.T) {
	tests := []struct {
		name         string
		session      *service.Session
		expectedCode int
		expectedErr  string
	}{
		{
			name: "Admin access",
			session: &service.Session{
				UserID: uuid.New(),
				Role:   domain.RoleAdmin,
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "Non-admin forbidden",
			session: &service.Session{
				UserID: uuid.New(),
				Role:   domain.RoleGeneral,
			},
			expectedCode: http.StatusForbidden,
			expectedErr:  "FORBIDDEN",
		},
		{
			name:         "Unauthenticated",
			session:      nil,
			expectedCode: http.StatusUnauthorized,
			expectedErr:  "UNAUTHORIZED",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			h := handler.NewUserHandler(mockRepo)

			req := httptest.NewRequest("GET", "/api/v1/users", nil)
			if tt.session != nil {
				ctx := context.WithValue(req.Context(), handler.ContextKeySession, tt.session)
				req = req.WithContext(ctx)
			}
			w := httptest.NewRecorder()

			h.ListUsers(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
			if tt.expectedErr != "" {
				assert.Contains(t, w.Body.String(), tt.expectedErr)
			}
		})
	}
}

func TestUserHandler_GetUser(t *testing.T) {
	userID := uuid.New()

	tests := []struct {
		name         string
		session      *service.Session
		setupMock    func(*MockUserRepository)
		expectedCode int
		expectedErr  string
	}{
		{
			name: "Admin access success",
			session: &service.Session{
				UserID: uuid.New(),
				Role:   domain.RoleAdmin,
			},
			setupMock: func(m *MockUserRepository) {
				m.On("GetByID", mock.Anything, userID).Return(&domain.User{ID: userID}, nil)
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "Non-admin forbidden",
			session: &service.Session{
				UserID: uuid.New(),
				Role:   domain.RoleGeneral,
			},
			setupMock:    func(m *MockUserRepository) {},
			expectedCode: http.StatusForbidden,
			expectedErr:  "FORBIDDEN",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			h := handler.NewUserHandler(mockRepo)

			tt.setupMock(mockRepo)

			req := httptest.NewRequest("GET", "/api/v1/users/"+userID.String(), nil)
			req = mux.SetURLVars(req, map[string]string{"id": userID.String()})
			if tt.session != nil {
				ctx := context.WithValue(req.Context(), handler.ContextKeySession, tt.session)
				req = req.WithContext(ctx)
			}
			w := httptest.NewRecorder()

			h.GetUser(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
			if tt.expectedErr != "" {
				assert.Contains(t, w.Body.String(), tt.expectedErr)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserHandler_GetCurrentUser_Unauthorized(t *testing.T) {
	mockRepo := new(MockUserRepository)
	h := handler.NewUserHandler(mockRepo)

	req := httptest.NewRequest("GET", "/api/v1/users/me", nil)
	w := httptest.NewRecorder()

	h.GetCurrentUser(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "UNAUTHORIZED")
}

func TestUserHandler_GetCurrentUser_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	h := handler.NewUserHandler(mockRepo)

	userID := uuid.New()
	session := &service.Session{
		UserID: userID,
		Role:   domain.RoleGeneral,
	}

	mockRepo.On("GetByID", mock.Anything, userID).Return(&domain.User{ID: userID}, nil)

	req := httptest.NewRequest("GET", "/api/v1/users/me", nil)
	ctx := context.WithValue(req.Context(), handler.ContextKeySession, session)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	h.GetCurrentUser(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockRepo.AssertExpectations(t)
}
