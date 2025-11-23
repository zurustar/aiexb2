// backend/internal/service/approval_service.go
package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/your-org/esms/internal/domain"
	"github.com/your-org/esms/internal/repository"
)

var (
	ErrAlreadyApproved = errors.New("reservation is already approved")
	ErrAlreadyRejected = errors.New("reservation is already rejected")
	ErrNotApprover     = errors.New("user is not an approver for this reservation")
)

// ApprovalService は承認に関するビジネスロジックを提供します
type ApprovalService struct {
	reservationRepo repository.ReservationRepository
	userRepo        repository.UserRepository
	auditLogRepo    repository.AuditLogRepository
}

// NewApprovalService は新しいApprovalServiceを作成します
func NewApprovalService(
	reservationRepo repository.ReservationRepository,
	userRepo repository.UserRepository,
	auditLogRepo repository.AuditLogRepository,
) *ApprovalService {
	return &ApprovalService{
		reservationRepo: reservationRepo,
		userRepo:        userRepo,
		auditLogRepo:    auditLogRepo,
	}
}

// ApproveReservation は予約を承認します
func (s *ApprovalService) ApproveReservation(ctx context.Context, reservationID uuid.UUID, startAt time.Time, approverID uuid.UUID) error {
	// 予約取得
	reservation, err := s.reservationRepo.GetByID(ctx, reservationID, startAt)
	if err != nil {
		return fmt.Errorf("failed to get reservation: %w", err)
	}

	// 承認状態チェック
	if reservation.ApprovalStatus == domain.ApprovalStatusConfirmed {
		return ErrAlreadyApproved
	}
	if reservation.ApprovalStatus == domain.ApprovalStatusRejected {
		return ErrAlreadyRejected
	}

	// 承認者権限チェック
	approver, err := s.userRepo.GetByID(ctx, approverID)
	if err != nil {
		return fmt.Errorf("failed to get approver: %w", err)
	}

	if !s.canApprove(approver, reservation) {
		return ErrNotApprover
	}

	// 承認処理
	reservation.ApprovalStatus = domain.ApprovalStatusConfirmed
	reservation.UpdatedBy = &approverID
	reservation.UpdatedAt = time.Now()

	err = s.reservationRepo.Update(ctx, reservation)
	if err != nil {
		return fmt.Errorf("failed to update reservation: %w", err)
	}

	// 監査ログ記録
	auditLog := &domain.AuditLog{
		ID:         uuid.New(),
		UserID:     approverID,
		Action:     domain.AuditActionApprove,
		TargetType: "reservation",
		TargetID:   reservationID.String(),
		Details: map[string]interface{}{
			"title":           reservation.Title,
			"organizer_id":    reservation.OrganizerID.String(),
			"approval_status": string(domain.ApprovalStatusConfirmed),
		},
		CreatedAt: time.Now(),
	}
	_ = s.auditLogRepo.Create(ctx, auditLog)

	return nil
}

// RejectReservation は予約を却下します
func (s *ApprovalService) RejectReservation(ctx context.Context, reservationID uuid.UUID, startAt time.Time, approverID uuid.UUID, reason string) error {
	// 予約取得
	reservation, err := s.reservationRepo.GetByID(ctx, reservationID, startAt)
	if err != nil {
		return fmt.Errorf("failed to get reservation: %w", err)
	}

	// 承認状態チェック
	if reservation.ApprovalStatus == domain.ApprovalStatusConfirmed {
		return ErrAlreadyApproved
	}
	if reservation.ApprovalStatus == domain.ApprovalStatusRejected {
		return ErrAlreadyRejected
	}

	// 承認者権限チェック
	approver, err := s.userRepo.GetByID(ctx, approverID)
	if err != nil {
		return fmt.Errorf("failed to get approver: %w", err)
	}

	if !s.canApprove(approver, reservation) {
		return ErrNotApprover
	}

	// 却下処理
	reservation.ApprovalStatus = domain.ApprovalStatusRejected
	reservation.UpdatedBy = &approverID
	reservation.UpdatedAt = time.Now()

	err = s.reservationRepo.Update(ctx, reservation)
	if err != nil {
		return fmt.Errorf("failed to update reservation: %w", err)
	}

	// 監査ログ記録
	auditLog := &domain.AuditLog{
		ID:         uuid.New(),
		UserID:     approverID,
		Action:     domain.AuditActionReject,
		TargetType: "reservation",
		TargetID:   reservationID.String(),
		Details: map[string]interface{}{
			"title":           reservation.Title,
			"organizer_id":    reservation.OrganizerID.String(),
			"approval_status": string(domain.ApprovalStatusRejected),
			"reason":          reason,
		},
		CreatedAt: time.Now(),
	}
	_ = s.auditLogRepo.Create(ctx, auditLog)

	return nil
}

// canApprove は指定されたユーザーが予約を承認できるかチェックします
func (s *ApprovalService) canApprove(user *domain.User, reservation *domain.Reservation) bool {
	// 管理者は全ての予約を承認可能
	if user.Role == domain.RoleAdmin {
		return true
	}

	// マネージャーは自分以外の予約を承認可能
	if user.Role == domain.RoleManager && user.ID != reservation.OrganizerID {
		return true
	}

	return false
}

// GetPendingApprovals は承認待ちの予約一覧を取得します（簡易実装）
func (s *ApprovalService) GetPendingApprovals(ctx context.Context, approverID uuid.UUID) ([]*domain.Reservation, error) {
	// 承認者権限チェック
	approver, err := s.userRepo.GetByID(ctx, approverID)
	if err != nil {
		return nil, fmt.Errorf("failed to get approver: %w", err)
	}

	// 管理者またはマネージャーのみ承認待ち一覧を取得可能
	if approver.Role != domain.RoleAdmin && approver.Role != domain.RoleManager {
		return nil, ErrUnauthorized
	}

	// TODO: リポジトリに承認待ち一覧取得メソッドを追加する必要がある
	// 現時点では空のリストを返す
	return []*domain.Reservation{}, nil
}
