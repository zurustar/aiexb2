// backend/internal/service/reservation_service.go
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
	ErrResourceNotAvailable = errors.New("resource is not available for the requested time")
	ErrInvalidTimeRange     = errors.New("invalid time range")
	ErrUnauthorized         = errors.New("unauthorized")
)

// ReservationService は予約に関するビジネスロジックを提供します
type ReservationService struct {
	reservationRepo repository.ReservationRepository
	resourceRepo    repository.ResourceRepository
	userRepo        repository.UserRepository
	auditLogRepo    repository.AuditLogRepository
}

// NewReservationService は新しいReservationServiceを作成します
func NewReservationService(
	reservationRepo repository.ReservationRepository,
	resourceRepo repository.ResourceRepository,
	userRepo repository.UserRepository,
	auditLogRepo repository.AuditLogRepository,
) *ReservationService {
	return &ReservationService{
		reservationRepo: reservationRepo,
		resourceRepo:    resourceRepo,
		userRepo:        userRepo,
		auditLogRepo:    auditLogRepo,
	}
}

// CreateReservationRequest は予約作成リクエスト
type CreateReservationRequest struct {
	OrganizerID uuid.UUID
	ResourceIDs []uuid.UUID
	Title       string
	Description string
	StartAt     time.Time
	EndAt       time.Time
	RRule       string
	IsPrivate   bool
	Timezone    string
}

// CreateReservation は新しい予約を作成します
func (s *ReservationService) CreateReservation(ctx context.Context, req *CreateReservationRequest) (*domain.Reservation, error) {
	// 入力検証
	if req.StartAt.After(req.EndAt) || req.StartAt.Equal(req.EndAt) {
		return nil, ErrInvalidTimeRange
	}

	// ユーザー存在確認
	user, err := s.userRepo.GetByID(ctx, req.OrganizerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// リソース存在確認と権限チェック
	for _, resourceID := range req.ResourceIDs {
		resource, err := s.resourceRepo.GetByID(ctx, resourceID)
		if err != nil {
			return nil, fmt.Errorf("failed to get resource: %w", err)
		}

		if !resource.CanBeReservedBy(user) {
			return nil, ErrUnauthorized
		}
	}

	// リソースの空き状況確認
	availableResources, err := s.resourceRepo.FindAvailable(ctx, req.StartAt, req.EndAt)
	if err != nil {
		return nil, fmt.Errorf("failed to find available resources: %w", err)
	}

	// リクエストされたリソースが全て利用可能か確認
	availableMap := make(map[uuid.UUID]bool)
	for _, r := range availableResources {
		availableMap[r.ID] = true
	}

	for _, resourceID := range req.ResourceIDs {
		if !availableMap[resourceID] {
			return nil, ErrResourceNotAvailable
		}
	}

	// 予約作成
	reservation := &domain.Reservation{
		ID:             uuid.New(),
		OrganizerID:    req.OrganizerID,
		Title:          req.Title,
		Description:    req.Description,
		StartAt:        req.StartAt,
		EndAt:          req.EndAt,
		RRule:          req.RRule,
		IsPrivate:      req.IsPrivate,
		Timezone:       req.Timezone,
		ApprovalStatus: domain.ApprovalStatusConfirmed,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// 予約インスタンス生成
	instances, err := reservation.ExpandInstances(req.StartAt, req.EndAt)
	if err != nil {
		return nil, fmt.Errorf("failed to expand instances: %w", err)
	}

	// ドメインインスタンスをリポジトリ用に変換
	repoInstances := make([]*domain.ReservationInstance, len(instances))
	for i, inst := range instances {
		repoInstances[i] = &inst
	}

	// トランザクション内で予約とインスタンスを作成
	err = s.reservationRepo.CreateWithInstances(ctx, reservation, repoInstances, req.ResourceIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to create reservation: %w", err)
	}

	// 監査ログ記録
	auditLog := &domain.AuditLog{
		ID:         uuid.New(),
		UserID:     req.OrganizerID,
		Action:     domain.AuditActionCreate,
		TargetType: "reservation",
		TargetID:   reservation.ID.String(),
		Details: map[string]interface{}{
			"title":     req.Title,
			"start_at":  req.StartAt,
			"end_at":    req.EndAt,
			"resources": len(req.ResourceIDs),
		},
		CreatedAt: time.Now(),
	}
	_ = s.auditLogRepo.Create(ctx, auditLog) // エラーは無視（監査ログ失敗で予約失敗にしない）

	return reservation, nil
}

// CancelReservation は予約をキャンセルします
func (s *ReservationService) CancelReservation(ctx context.Context, reservationID uuid.UUID, startAt time.Time, userID uuid.UUID) error {
	// 予約取得
	reservation, err := s.reservationRepo.GetByID(ctx, reservationID, startAt)
	if err != nil {
		return fmt.Errorf("failed to get reservation: %w", err)
	}

	// 権限チェック（主催者のみキャンセル可能）
	if reservation.OrganizerID != userID {
		return ErrUnauthorized
	}

	// 予約削除
	err = s.reservationRepo.Delete(ctx, reservationID, startAt)
	if err != nil {
		return fmt.Errorf("failed to delete reservation: %w", err)
	}

	// 監査ログ記録
	auditLog := &domain.AuditLog{
		ID:         uuid.New(),
		UserID:     userID,
		Action:     domain.AuditActionCancel,
		TargetType: "reservation",
		TargetID:   reservationID.String(),
		Details: map[string]interface{}{
			"title": reservation.Title,
		},
		CreatedAt: time.Now(),
	}
	_ = s.auditLogRepo.Create(ctx, auditLog)

	return nil
}

// FindAlternativeResources は代替リソースを提案します
func (s *ReservationService) FindAlternativeResources(ctx context.Context, startAt, endAt time.Time, resourceType domain.ResourceType) ([]*domain.Resource, error) {
	// 指定時間帯に利用可能なリソースを検索
	availableResources, err := s.resourceRepo.FindAvailable(ctx, startAt, endAt)
	if err != nil {
		return nil, fmt.Errorf("failed to find available resources: %w", err)
	}

	// リソースタイプでフィルタリング
	var filtered []*domain.Resource
	for _, r := range availableResources {
		if r.Type == resourceType {
			filtered = append(filtered, r)
		}
	}

	return filtered, nil
}
