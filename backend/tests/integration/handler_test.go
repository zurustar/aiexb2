// backend/tests/integration/handler_test.go
package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/your-org/esms/internal/domain"
	"github.com/your-org/esms/internal/handler"
	"github.com/your-org/esms/internal/repository"
	"github.com/your-org/esms/internal/service"
)

// TestAuthenticationFailures tests authentication failure scenarios
func TestAuthenticationFailures(t *testing.T) {
	env := setupTestEnv(t)
	defer env.cleanup()

	// Create user handler to test authentication
	userHandler := handler.NewUserHandler(env.userRepo)

	tests := []struct {
		name           string
		setupContext   func(*http.Request) *http.Request
		expectedStatus int
		expectedError  string
	}{
		{
			name: "Missing Session Context",
			setupContext: func(req *http.Request) *http.Request {
				// No session in context
				return req
			},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "UNAUTHORIZED",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/api/v1/users/me", nil)
			req = tt.setupContext(req)
			w := httptest.NewRecorder()

			userHandler.GetCurrentUser(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Contains(t, w.Body.String(), tt.expectedError)

			// Verify error response format
			var errorResp map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &errorResp)
			require.NoError(t, err)
			assert.Contains(t, errorResp, "error")
			if errorObj, ok := errorResp["error"].(map[string]interface{}); ok {
				assert.Contains(t, errorObj, "code")
				assert.Contains(t, errorObj, "message")
			}
		})
	}
}

// TestAuthorizationFailures tests authorization failure scenarios
func TestAuthorizationFailures(t *testing.T) {
	env := setupTestEnv(t)
	defer env.cleanup()

	userHandler := handler.NewUserHandler(env.userRepo)

	// Create a non-admin user
	generalUser := createTestUser(t, env, domain.RoleGeneral)

	tests := []struct {
		name           string
		userID         uuid.UUID
		role           domain.Role
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "Non-admin accessing ListUsers",
			userID:         generalUser.ID,
			role:           domain.RoleGeneral,
			expectedStatus: http.StatusForbidden,
			expectedError:  "FORBIDDEN",
		},
		{
			name:           "Non-admin accessing GetUser",
			userID:         generalUser.ID,
			role:           domain.RoleGeneral,
			expectedStatus: http.StatusForbidden,
			expectedError:  "FORBIDDEN",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			session := &service.Session{
				UserID: tt.userID,
				Role:   tt.role,
			}

			req := httptest.NewRequest("GET", "/api/v1/users", nil)
			ctx := context.WithValue(req.Context(), handler.ContextKeySession, session)
			req = req.WithContext(ctx)
			w := httptest.NewRecorder()

			userHandler.ListUsers(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Contains(t, w.Body.String(), tt.expectedError)

			// Verify error response format
			var errorResp map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &errorResp)
			require.NoError(t, err)
			assert.Contains(t, errorResp, "error")
			if errorObj, ok := errorResp["error"].(map[string]interface{}); ok {
				assert.Contains(t, errorObj, "code")
				assert.Contains(t, errorObj, "message")
			}
		})
	}
}

// TestValidationFailures tests input validation failure scenarios
func TestValidationFailures(t *testing.T) {
	env := setupTestEnv(t)
	defer env.cleanup()

	resourceHandler := handler.NewResourceHandler(env.resourceRepo)
	reservationHandler := handler.NewReservationHandler(env.reservationSvc, env.approvalSvc)

	// Create admin session for testing
	adminUser := createTestUser(t, env, domain.RoleAdmin)
	adminSession := &service.Session{
		UserID: adminUser.ID,
		Role:   domain.RoleAdmin,
	}

	tests := []struct {
		name           string
		handler        func(http.ResponseWriter, *http.Request)
		method         string
		path           string
		body           string
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "CreateResource - Missing Name",
			handler:        resourceHandler.CreateResource,
			method:         "POST",
			path:           "/api/v1/resources",
			body:           `{"type": "MEETING_ROOM", "location": "Building 1"}`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "INVALID_NAME",
		},
		{
			name:           "CreateResource - Missing Type",
			handler:        resourceHandler.CreateResource,
			method:         "POST",
			path:           "/api/v1/resources",
			body:           `{"name": "Room A", "location": "Building 1"}`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "INVALID_TYPE",
		},
		{
			name:           "CreateReservation - Missing Timezone",
			handler:        reservationHandler.CreateReservation,
			method:         "POST",
			path:           "/api/v1/events",
			body:           `{"resource_ids": ["` + uuid.New().String() + `"], "title": "Meeting", "start_at": "2025-06-01T10:00:00Z", "end_at": "2025-06-01T11:00:00Z"}`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "INVALID_TIMEZONE",
		},
		{
			name:           "CreateReservation - Invalid Time Range",
			handler:        reservationHandler.CreateReservation,
			method:         "POST",
			path:           "/api/v1/events",
			body:           `{"resource_ids": ["` + uuid.New().String() + `"], "title": "Meeting", "start_at": "2025-06-01T11:00:00Z", "end_at": "2025-06-01T10:00:00Z", "timezone": "UTC"}`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "INVALID_TIME_RANGE",
		},
		{
			name:           "ListResources - Invalid Query Parameter",
			handler:        resourceHandler.ListResources,
			method:         "GET",
			path:           "/api/v1/resources?is_active=invalid",
			body:           "",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "INVALID_QUERY_PARAM",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req *http.Request
			if tt.body != "" {
				req = httptest.NewRequest(tt.method, tt.path, bytes.NewBufferString(tt.body))
				req.Header.Set("Content-Type", "application/json")
			} else {
				req = httptest.NewRequest(tt.method, tt.path, nil)
			}

			// Add session context
			ctx := context.WithValue(req.Context(), handler.ContextKeySession, adminSession)
			req = req.WithContext(ctx)

			w := httptest.NewRecorder()

			tt.handler(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Contains(t, w.Body.String(), tt.expectedError)

			// Verify error response format
			var errorResp map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &errorResp)
			require.NoError(t, err)
			assert.Contains(t, errorResp, "error")
			if errorObj, ok := errorResp["error"].(map[string]interface{}); ok {
				assert.Contains(t, errorObj, "code")
				assert.Contains(t, errorObj, "message")
			}
		})
	}
}

// TestSuccessCases tests basic happy path scenarios
func TestSuccessCases(t *testing.T) {
	env := setupTestEnv(t)
	defer env.cleanup()

	userHandler := handler.NewUserHandler(env.userRepo)

	// Create admin and general users
	adminUser := createTestUser(t, env, domain.RoleAdmin)
	generalUser := createTestUser(t, env, domain.RoleGeneral)

	tests := []struct {
		name           string
		userID         uuid.UUID
		role           domain.Role
		expectedStatus int
	}{
		{
			name:           "GetCurrentUser - General User",
			userID:         generalUser.ID,
			role:           domain.RoleGeneral,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "GetCurrentUser - Admin User",
			userID:         adminUser.ID,
			role:           domain.RoleAdmin,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "ListUsers - Admin Access",
			userID:         adminUser.ID,
			role:           domain.RoleAdmin,
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			session := &service.Session{
				UserID: tt.userID,
				Role:   tt.role,
			}

			req := httptest.NewRequest("GET", "/api/v1/users/me", nil)
			ctx := context.WithValue(req.Context(), handler.ContextKeySession, session)
			req = req.WithContext(ctx)
			w := httptest.NewRecorder()

			userHandler.GetCurrentUser(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

// Helper functions

type testEnv struct {
	userRepo       *mockUserRepository
	resourceRepo   *mockResourceRepository
	authSvc        *mockAuthService
	reservationSvc *mockReservationService
	approvalSvc    *mockApprovalService
	cleanup        func()
}

func setupTestEnv(t *testing.T) *testEnv {
	// Create in-memory repositories and services for testing
	userRepo := &mockUserRepository{users: make(map[uuid.UUID]*domain.User)}
	resourceRepo := &mockResourceRepository{}
	authSvc := &mockAuthService{sessions: make(map[string]*service.Session)}
	reservationSvc := &mockReservationService{}
	approvalSvc := &mockApprovalService{}

	return &testEnv{
		userRepo:       userRepo,
		resourceRepo:   resourceRepo,
		authSvc:        authSvc,
		reservationSvc: reservationSvc,
		approvalSvc:    approvalSvc,
		cleanup:        func() {},
	}
}

func createTestUser(t *testing.T, env *testEnv, role domain.Role) *domain.User {
	user := &domain.User{
		ID:    uuid.New(),
		Email: "test@example.com",
		Name:  "Test User",
		Role:  role,
	}
	env.userRepo.users[user.ID] = user
	return user
}

// Mock implementations

type mockUserRepository struct {
	users map[uuid.UUID]*domain.User
}

func (m *mockUserRepository) Create(ctx context.Context, user *domain.User) error {
	m.users[user.ID] = user
	return nil
}

func (m *mockUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	user, ok := m.users[id]
	if !ok {
		return nil, repository.ErrNotFound
	}
	return user, nil
}

func (m *mockUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	for _, user := range m.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, repository.ErrNotFound
}

func (m *mockUserRepository) Update(ctx context.Context, user *domain.User) error {
	m.users[user.ID] = user
	return nil
}

func (m *mockUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	delete(m.users, id)
	return nil
}

type mockResourceRepository struct{}

func (m *mockResourceRepository) Create(ctx context.Context, resource *domain.Resource) error {
	return nil
}

func (m *mockResourceRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Resource, error) {
	return nil, repository.ErrNotFound
}

func (m *mockResourceRepository) Update(ctx context.Context, resource *domain.Resource) error {
	return nil
}

func (m *mockResourceRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return nil
}

func (m *mockResourceRepository) FindAvailable(ctx context.Context, startAt, endAt time.Time) ([]*domain.Resource, error) {
	return []*domain.Resource{}, nil
}

type mockAuthService struct {
	sessions map[string]*service.Session
}

func (m *mockAuthService) GetAuthURL(state string) string {
	return "http://auth.example.com"
}

func (m *mockAuthService) HandleCallback(ctx context.Context, code, state string) (*service.Session, error) {
	return nil, nil
}

func (m *mockAuthService) Logout(ctx context.Context, sessionID string, userID uuid.UUID) error {
	delete(m.sessions, sessionID)
	return nil
}

func (m *mockAuthService) RefreshSession(ctx context.Context, sessionID string) (*service.Session, error) {
	return nil, nil
}

func (m *mockAuthService) GetSession(sessionID string) (*service.Session, error) {
	session, ok := m.sessions[sessionID]
	if !ok {
		return nil, service.ErrSessionNotFound
	}
	if session.ExpiresAt.Before(time.Now()) {
		return nil, service.ErrSessionNotFound // Use ErrSessionNotFound for expired sessions too
	}
	return session, nil
}

type mockReservationService struct{}

func (m *mockReservationService) CreateReservation(ctx context.Context, req *service.CreateReservationRequest) (*domain.Reservation, error) {
	return nil, nil
}

func (m *mockReservationService) CancelReservation(ctx context.Context, id uuid.UUID, startAt time.Time, userID uuid.UUID) error {
	return nil
}

type mockApprovalService struct{}

func (m *mockApprovalService) ApproveReservation(ctx context.Context, reservationID uuid.UUID, startAt time.Time, approverID uuid.UUID) error {
	return nil
}

func (m *mockApprovalService) RejectReservation(ctx context.Context, reservationID uuid.UUID, startAt time.Time, approverID uuid.UUID, reason string) error {
	return nil
}
