// backend/internal/service/reservation_service_test.go
package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/your-org/esms/internal/domain"
	"github.com/your-org/esms/internal/service"
)

func TestReservationService_CreateReservation_Success(t *testing.T) {
	mockReservationRepo := new(MockReservationRepository)
	mockResourceRepo := new(MockResourceRepository)
	mockUserRepo := new(MockUserRepository)
	mockAuditLogRepo := new(MockAuditLogRepository)

	svc := service.NewReservationService(mockReservationRepo, mockResourceRepo, mockUserRepo, mockAuditLogRepo)

	ctx := context.Background()
	userID := uuid.New()
	resourceID := uuid.New()
	startAt := time.Now().Add(24 * time.Hour)
	endAt := startAt.Add(1 * time.Hour)

	user := &domain.User{
		ID:       userID,
		Email:    "test@example.com",
		Name:     "Test User",
		Role:     domain.RoleGeneral,
		IsActive: true,
	}

	capacity := 10
	resource := &domain.Resource{
		ID:       resourceID,
		Name:     "Test Room",
		Type:     domain.ResourceTypeMeetingRoom,
		Capacity: &capacity,
		IsActive: true,
	}

	mockUserRepo.On("GetByID", ctx, userID).Return(user, nil)
	mockResourceRepo.On("GetByID", ctx, resourceID).Return(resource, nil)
	mockResourceRepo.On("FindAvailable", ctx, mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time")).Return([]*domain.Resource{resource}, nil)
	mockReservationRepo.On("CreateWithInstances", ctx, mock.AnythingOfType("*domain.Reservation"), mock.AnythingOfType("[]*domain.ReservationInstance"), mock.AnythingOfType("[]uuid.UUID")).Return(nil)
	mockAuditLogRepo.On("Create", ctx, mock.AnythingOfType("*domain.AuditLog")).Return(nil)

	req := &service.CreateReservationRequest{
		OrganizerID: userID,
		ResourceIDs: []uuid.UUID{resourceID},
		Title:       "Test Meeting",
		Description: "Test",
		StartAt:     startAt,
		EndAt:       endAt,
		Timezone:    "Asia/Tokyo",
	}

	reservation, err := svc.CreateReservation(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, reservation)
	assert.Equal(t, req.Title, reservation.Title)
	assert.Equal(t, req.OrganizerID, reservation.OrganizerID)

	mockUserRepo.AssertExpectations(t)
	mockResourceRepo.AssertExpectations(t)
	mockReservationRepo.AssertExpectations(t)
}

func TestReservationService_CreateReservation_ResourceNotAvailable(t *testing.T) {
	mockReservationRepo := new(MockReservationRepository)
	mockResourceRepo := new(MockResourceRepository)
	mockUserRepo := new(MockUserRepository)
	mockAuditLogRepo := new(MockAuditLogRepository)

	svc := service.NewReservationService(mockReservationRepo, mockResourceRepo, mockUserRepo, mockAuditLogRepo)

	ctx := context.Background()
	userID := uuid.New()
	resourceID := uuid.New()
	startAt := time.Now().Add(24 * time.Hour)
	endAt := startAt.Add(1 * time.Hour)

	user := &domain.User{
		ID:       userID,
		Role:     domain.RoleGeneral,
		IsActive: true,
	}

	capacity := 10
	resource := &domain.Resource{
		ID:       resourceID,
		Name:     "Test Room",
		Type:     domain.ResourceTypeMeetingRoom,
		Capacity: &capacity,
		IsActive: true,
	}

	mockUserRepo.On("GetByID", ctx, userID).Return(user, nil)
	mockResourceRepo.On("GetByID", ctx, resourceID).Return(resource, nil)
	mockResourceRepo.On("FindAvailable", ctx, mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time")).Return([]*domain.Resource{}, nil) // 空き無し

	req := &service.CreateReservationRequest{
		OrganizerID: userID,
		ResourceIDs: []uuid.UUID{resourceID},
		Title:       "Test Meeting",
		StartAt:     startAt,
		EndAt:       endAt,
	}

	reservation, err := svc.CreateReservation(ctx, req)

	assert.Error(t, err)
	assert.Equal(t, service.ErrResourceNotAvailable, err)
	assert.Nil(t, reservation)
}

func TestReservationService_CancelReservation_Success(t *testing.T) {
	mockReservationRepo := new(MockReservationRepository)
	mockResourceRepo := new(MockResourceRepository)
	mockUserRepo := new(MockUserRepository)
	mockAuditLogRepo := new(MockAuditLogRepository)

	svc := service.NewReservationService(mockReservationRepo, mockResourceRepo, mockUserRepo, mockAuditLogRepo)

	ctx := context.Background()
	reservationID := uuid.New()
	userID := uuid.New()
	startAt := time.Now()

	reservation := &domain.Reservation{
		ID:          reservationID,
		OrganizerID: userID,
		Title:       "Test Meeting",
		StartAt:     startAt,
	}

	mockReservationRepo.On("GetByID", ctx, reservationID, startAt).Return(reservation, nil)
	mockReservationRepo.On("Delete", ctx, reservationID, startAt).Return(nil)
	mockAuditLogRepo.On("Create", ctx, mock.AnythingOfType("*domain.AuditLog")).Return(nil)

	err := svc.CancelReservation(ctx, reservationID, startAt, userID)

	assert.NoError(t, err)
	mockReservationRepo.AssertExpectations(t)
}
