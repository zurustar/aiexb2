package handler_test

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/your-org/esms/internal/service"
)

// MockAuthService for handler tests
type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) GetAuthURL(state string) string {
	args := m.Called(state)
	return args.String(0)
}

func (m *MockAuthService) HandleCallback(ctx context.Context, code, state string) (*service.Session, error) {
	args := m.Called(ctx, code, state)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*service.Session), args.Error(1)
}

func (m *MockAuthService) Logout(ctx context.Context, sessionID string, userID uuid.UUID) error {
	args := m.Called(ctx, sessionID, userID)
	return args.Error(0)
}

func (m *MockAuthService) RefreshSession(ctx context.Context, sessionID string) (*service.Session, error) {
	args := m.Called(ctx, sessionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*service.Session), args.Error(1)
}

func (m *MockAuthService) GetSession(sessionID string) (*service.Session, error) {
	args := m.Called(sessionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*service.Session), args.Error(1)
}
