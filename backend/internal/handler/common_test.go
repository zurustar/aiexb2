package handler_test

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/your-org/esms/internal/domain"
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

// MockReservationService for handler tests
type MockReservationService struct {
	mock.Mock
}

func (m *MockReservationService) CreateReservation(ctx context.Context, req *service.CreateReservationRequest) (*domain.Reservation, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Reservation), args.Error(1)
}

func (m *MockReservationService) CancelReservation(ctx context.Context, id uuid.UUID, startAt time.Time, userID uuid.UUID) error {
	args := m.Called(ctx, id, startAt, userID)
	return args.Error(0)
}

// MockApprovalService for handler tests
type MockApprovalService struct {
	mock.Mock
}

func (m *MockApprovalService) ApproveReservation(ctx context.Context, reservationID uuid.UUID, startAt time.Time, approverID uuid.UUID) error {
	args := m.Called(ctx, reservationID, startAt, approverID)
	return args.Error(0)
}

func (m *MockApprovalService) RejectReservation(ctx context.Context, reservationID uuid.UUID, startAt time.Time, approverID uuid.UUID, reason string) error {
	args := m.Called(ctx, reservationID, startAt, approverID, reason)
	return args.Error(0)
}
