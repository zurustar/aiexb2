// backend/internal/service/mock_test.go
package service_test

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/your-org/esms/internal/domain"
)

// Mock Repositories

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
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

func (m *MockUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) List(ctx context.Context, offset, limit int) ([]*domain.User, int64, error) {
	args := m.Called(ctx, offset, limit)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*domain.User), args.Get(1).(int64), args.Error(2)
}

type MockAuditLogRepository struct {
	mock.Mock
}

func (m *MockAuditLogRepository) Create(ctx context.Context, log *domain.AuditLog) error {
	args := m.Called(ctx, log)
	return args.Error(0)
}

func (m *MockAuditLogRepository) List(ctx context.Context, filter domain.AuditLogFilter, offset, limit int) ([]*domain.AuditLog, int64, error) {
	args := m.Called(ctx, filter, offset, limit)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*domain.AuditLog), args.Get(1).(int64), args.Error(2)
}

func (m *MockAuditLogRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.AuditLog, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.AuditLog), args.Error(1)
}

func (m *MockAuditLogRepository) GetByEntityID(ctx context.Context, entityID uuid.UUID, limit int) ([]*domain.AuditLog, error) {
	args := m.Called(ctx, entityID, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.AuditLog), args.Error(1)
}

type MockReservationRepository struct {
	mock.Mock
}

func (m *MockReservationRepository) Create(ctx context.Context, reservation *domain.Reservation) error {
	args := m.Called(ctx, reservation)
	return args.Error(0)
}

func (m *MockReservationRepository) CreateWithInstances(ctx context.Context, reservation *domain.Reservation, instances []*domain.ReservationInstance, resourceIDs []uuid.UUID) error {
	args := m.Called(ctx, reservation, instances, resourceIDs)
	return args.Error(0)
}

func (m *MockReservationRepository) GetByID(ctx context.Context, id uuid.UUID, startAt time.Time) (*domain.Reservation, error) {
	args := m.Called(ctx, id, startAt)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Reservation), args.Error(1)
}

func (m *MockReservationRepository) Update(ctx context.Context, reservation *domain.Reservation) error {
	args := m.Called(ctx, reservation)
	return args.Error(0)
}

func (m *MockReservationRepository) Delete(ctx context.Context, id uuid.UUID, startAt time.Time) error {
	args := m.Called(ctx, id, startAt)
	return args.Error(0)
}

func (m *MockReservationRepository) GetInstancesByReservationID(ctx context.Context, reservationID uuid.UUID) ([]*domain.ReservationInstance, error) {
	args := m.Called(ctx, reservationID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.ReservationInstance), args.Error(1)
}

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
