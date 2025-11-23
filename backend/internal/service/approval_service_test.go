// backend/internal/service/approval_service_test.go
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

func TestApprovalService_ApproveReservation_Success(t *testing.T) {
	mockReservationRepo := new(MockReservationRepository)
	mockUserRepo := new(MockUserRepository)
	mockAuditLogRepo := new(MockAuditLogRepository)

	svc := service.NewApprovalService(mockReservationRepo, mockUserRepo, mockAuditLogRepo)

	ctx := context.Background()
	reservationID := uuid.New()
	approverID := uuid.New()
	startAt := time.Now()

	reservation := &domain.Reservation{
		ID:             reservationID,
		OrganizerID:    uuid.New(),
		Title:          "Test Meeting",
		StartAt:        startAt,
		ApprovalStatus: domain.ApprovalStatusPending,
	}

	approver := &domain.User{
		ID:       approverID,
		Role:     domain.RoleAdmin,
		IsActive: true,
	}

	mockReservationRepo.On("GetByID", ctx, reservationID, startAt).Return(reservation, nil)
	mockUserRepo.On("GetByID", ctx, approverID).Return(approver, nil)
	mockReservationRepo.On("Update", ctx, mock.MatchedBy(func(r *domain.Reservation) bool {
		return r.ApprovalStatus == domain.ApprovalStatusConfirmed
	})).Return(nil)
	mockAuditLogRepo.On("Create", ctx, mock.AnythingOfType("*domain.AuditLog")).Return(nil)

	err := svc.ApproveReservation(ctx, reservationID, startAt, approverID)

	assert.NoError(t, err)
	mockReservationRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

func TestApprovalService_ApproveReservation_AlreadyApproved(t *testing.T) {
	mockReservationRepo := new(MockReservationRepository)
	mockUserRepo := new(MockUserRepository)
	mockAuditLogRepo := new(MockAuditLogRepository)

	svc := service.NewApprovalService(mockReservationRepo, mockUserRepo, mockAuditLogRepo)

	ctx := context.Background()
	reservationID := uuid.New()
	approverID := uuid.New()
	startAt := time.Now()

	reservation := &domain.Reservation{
		ID:             reservationID,
		ApprovalStatus: domain.ApprovalStatusConfirmed, // 既に承認済み
	}

	mockReservationRepo.On("GetByID", ctx, reservationID, startAt).Return(reservation, nil)

	err := svc.ApproveReservation(ctx, reservationID, startAt, approverID)

	assert.Error(t, err)
	assert.Equal(t, service.ErrAlreadyApproved, err)
}

func TestApprovalService_ApproveReservation_NotApprover(t *testing.T) {
	mockReservationRepo := new(MockReservationRepository)
	mockUserRepo := new(MockUserRepository)
	mockAuditLogRepo := new(MockAuditLogRepository)

	svc := service.NewApprovalService(mockReservationRepo, mockUserRepo, mockAuditLogRepo)

	ctx := context.Background()
	reservationID := uuid.New()
	approverID := uuid.New()
	startAt := time.Now()

	reservation := &domain.Reservation{
		ID:             reservationID,
		OrganizerID:    uuid.New(),
		ApprovalStatus: domain.ApprovalStatusPending,
	}

	approver := &domain.User{
		ID:       approverID,
		Role:     domain.RoleGeneral, // 一般ユーザーは承認不可
		IsActive: true,
	}

	mockReservationRepo.On("GetByID", ctx, reservationID, startAt).Return(reservation, nil)
	mockUserRepo.On("GetByID", ctx, approverID).Return(approver, nil)

	err := svc.ApproveReservation(ctx, reservationID, startAt, approverID)

	assert.Error(t, err)
	assert.Equal(t, service.ErrNotApprover, err)
}

func TestApprovalService_RejectReservation_Success(t *testing.T) {
	mockReservationRepo := new(MockReservationRepository)
	mockUserRepo := new(MockUserRepository)
	mockAuditLogRepo := new(MockAuditLogRepository)

	svc := service.NewApprovalService(mockReservationRepo, mockUserRepo, mockAuditLogRepo)

	ctx := context.Background()
	reservationID := uuid.New()
	approverID := uuid.New()
	startAt := time.Now()
	reason := "リソースが不足しています"

	reservation := &domain.Reservation{
		ID:             reservationID,
		OrganizerID:    uuid.New(),
		Title:          "Test Meeting",
		StartAt:        startAt,
		ApprovalStatus: domain.ApprovalStatusPending,
	}

	approver := &domain.User{
		ID:       approverID,
		Role:     domain.RoleManager,
		IsActive: true,
	}

	mockReservationRepo.On("GetByID", ctx, reservationID, startAt).Return(reservation, nil)
	mockUserRepo.On("GetByID", ctx, approverID).Return(approver, nil)
	mockReservationRepo.On("Update", ctx, mock.MatchedBy(func(r *domain.Reservation) bool {
		return r.ApprovalStatus == domain.ApprovalStatusRejected
	})).Return(nil)
	mockAuditLogRepo.On("Create", ctx, mock.AnythingOfType("*domain.AuditLog")).Return(nil)

	err := svc.RejectReservation(ctx, reservationID, startAt, approverID, reason)

	assert.NoError(t, err)
	mockReservationRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

func TestApprovalService_RejectReservation_ManagerCannotRejectOwnReservation(t *testing.T) {
	mockReservationRepo := new(MockReservationRepository)
	mockUserRepo := new(MockUserRepository)
	mockAuditLogRepo := new(MockAuditLogRepository)

	svc := service.NewApprovalService(mockReservationRepo, mockUserRepo, mockAuditLogRepo)

	ctx := context.Background()
	managerID := uuid.New()
	reservationID := uuid.New()
	startAt := time.Now()

	reservation := &domain.Reservation{
		ID:             reservationID,
		OrganizerID:    managerID, // マネージャー自身が主催者
		ApprovalStatus: domain.ApprovalStatusPending,
	}

	manager := &domain.User{
		ID:       managerID,
		Role:     domain.RoleManager,
		IsActive: true,
	}

	mockReservationRepo.On("GetByID", ctx, reservationID, startAt).Return(reservation, nil)
	mockUserRepo.On("GetByID", ctx, managerID).Return(manager, nil)

	err := svc.RejectReservation(ctx, reservationID, startAt, managerID, "test")

	assert.Error(t, err)
	assert.Equal(t, service.ErrNotApprover, err)
}
