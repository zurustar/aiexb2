// backend/tests/integration/service_test.go
package integration_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/your-org/esms/internal/domain"
	"github.com/your-org/esms/internal/repository"
	"github.com/your-org/esms/internal/service"
)

// TestServiceIntegration_ReservationFlow は予約フロー全体の統合テスト
func TestServiceIntegration_ReservationFlow(t *testing.T) {
	if testDB == nil {
		t.Skip("TEST_DATABASE_URL not set, skipping integration test")
	}

	cleanupDB(t)
	defer cleanupDB(t)

	ctx := context.Background()

	// リポジトリ初期化
	userRepo := repository.NewUserRepository(testDB)
	resourceRepo := repository.NewResourceRepository(testDB)
	reservationRepo := repository.NewReservationRepository(testDB)
	auditLogRepo := repository.NewAuditLogRepository(testDB)

	// サービス初期化
	reservationService := service.NewReservationService(
		reservationRepo,
		resourceRepo,
		userRepo,
		auditLogRepo,
	)
	approvalService := service.NewApprovalService(
		reservationRepo,
		userRepo,
		auditLogRepo,
	)

	// テストデータ準備: ユーザー
	organizer := &domain.User{
		ID:        uuid.New(),
		Email:     "organizer@example.com",
		Name:      "Organizer",
		Role:      domain.RoleGeneral,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	require.NoError(t, userRepo.Create(ctx, organizer))

	approver := &domain.User{
		ID:        uuid.New(),
		Email:     "approver@example.com",
		Name:      "Approver",
		Role:      domain.RoleManager,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	require.NoError(t, userRepo.Create(ctx, approver))

	// テストデータ準備: リソース
	capacity := 10
	resource := &domain.Resource{
		ID:        uuid.New(),
		Name:      "Meeting Room A",
		Type:      domain.ResourceTypeMeetingRoom,
		Capacity:  &capacity,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	require.NoError(t, resourceRepo.Create(ctx, resource))

	// 1. 予約作成
	startAt := time.Date(2025, 6, 1, 10, 0, 0, 0, time.UTC)
	req := &service.CreateReservationRequest{
		OrganizerID: organizer.ID,
		ResourceIDs: []uuid.UUID{resource.ID},
		Title:       "Integration Test Meeting",
		Description: "Testing reservation flow",
		StartAt:     startAt,
		EndAt:       startAt.Add(1 * time.Hour),
		Timezone:    "Asia/Tokyo",
	}

	reservation, err := reservationService.CreateReservation(ctx, req)
	require.NoError(t, err)
	assert.NotNil(t, reservation)
	assert.Equal(t, "Integration Test Meeting", reservation.Title)
	assert.Equal(t, domain.ApprovalStatusConfirmed, reservation.ApprovalStatus)

	// 2. 予約を取得して確認
	retrieved, err := reservationRepo.GetByID(ctx, reservation.ID, reservation.StartAt)
	require.NoError(t, err)
	assert.Equal(t, reservation.Title, retrieved.Title)

	// 3. 承認フロー（既に確認済みの場合はエラー）
	err = approvalService.ApproveReservation(ctx, reservation.ID, reservation.StartAt, approver.ID)
	assert.Error(t, err) // 既に承認済みなのでエラー
	assert.Equal(t, service.ErrAlreadyApproved, err)

	// 4. 予約キャンセル
	err = reservationService.CancelReservation(ctx, reservation.ID, reservation.StartAt, organizer.ID)
	require.NoError(t, err)

	// 5. キャンセル後は取得できない
	_, err = reservationRepo.GetByID(ctx, reservation.ID, reservation.StartAt)
	assert.Equal(t, repository.ErrNotFound, err)
}

// TestServiceIntegration_ApprovalFlow は承認フロー全体の統合テスト
func TestServiceIntegration_ApprovalFlow(t *testing.T) {
	if testDB == nil {
		t.Skip("TEST_DATABASE_URL not set, skipping integration test")
	}

	cleanupDB(t)
	defer cleanupDB(t)

	ctx := context.Background()

	// リポジトリ初期化
	userRepo := repository.NewUserRepository(testDB)
	resourceRepo := repository.NewResourceRepository(testDB)
	reservationRepo := repository.NewReservationRepository(testDB)
	auditLogRepo := repository.NewAuditLogRepository(testDB)

	// サービス初期化
	approvalService := service.NewApprovalService(
		reservationRepo,
		userRepo,
		auditLogRepo,
	)

	// テストデータ準備
	organizer := &domain.User{
		ID:        uuid.New(),
		Email:     "user@example.com",
		Name:      "User",
		Role:      domain.RoleGeneral,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	require.NoError(t, userRepo.Create(ctx, organizer))

	manager := &domain.User{
		ID:        uuid.New(),
		Email:     "manager@example.com",
		Name:      "Manager",
		Role:      domain.RoleManager,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	require.NoError(t, userRepo.Create(ctx, manager))

	// 承認待ち予約を作成
	startAt := time.Date(2025, 6, 1, 14, 0, 0, 0, time.UTC)
	reservation := &domain.Reservation{
		ID:             uuid.New(),
		OrganizerID:    organizer.ID,
		Title:          "Pending Reservation",
		StartAt:        startAt,
		EndAt:          startAt.Add(1 * time.Hour),
		Timezone:       "Asia/Tokyo",
		ApprovalStatus: domain.ApprovalStatusPending,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	require.NoError(t, reservationRepo.Create(ctx, reservation))

	// 1. マネージャーが承認
	err := approvalService.ApproveReservation(ctx, reservation.ID, reservation.StartAt, manager.ID)
	require.NoError(t, err)

	// 2. 承認後のステータス確認
	approved, err := reservationRepo.GetByID(ctx, reservation.ID, reservation.StartAt)
	require.NoError(t, err)
	assert.Equal(t, domain.ApprovalStatusConfirmed, approved.ApprovalStatus)

	// 3. 監査ログ確認
	logs, err := auditLogRepo.GetByEntityID(ctx, reservation.ID, 10)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(logs), 1)

	// 承認アクションのログを確認
	var foundApproval bool
	for _, log := range logs {
		if log.Action == domain.AuditActionApprove {
			foundApproval = true
			assert.Equal(t, manager.ID, log.UserID)
			break
		}
	}
	assert.True(t, foundApproval, "Approval audit log should exist")
}

// TestServiceIntegration_ResourceConflict はリソース競合の統合テスト
func TestServiceIntegration_ResourceConflict(t *testing.T) {
	if testDB == nil {
		t.Skip("TEST_DATABASE_URL not set, skipping integration test")
	}

	cleanupDB(t)
	defer cleanupDB(t)

	ctx := context.Background()

	// リポジトリ初期化
	userRepo := repository.NewUserRepository(testDB)
	resourceRepo := repository.NewResourceRepository(testDB)
	reservationRepo := repository.NewReservationRepository(testDB)
	auditLogRepo := repository.NewAuditLogRepository(testDB)

	// サービス初期化
	reservationService := service.NewReservationService(
		reservationRepo,
		resourceRepo,
		userRepo,
		auditLogRepo,
	)

	// テストデータ準備
	user := &domain.User{
		ID:        uuid.New(),
		Email:     "user@example.com",
		Name:      "User",
		Role:      domain.RoleGeneral,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	require.NoError(t, userRepo.Create(ctx, user))

	capacity := 5
	resource := &domain.Resource{
		ID:        uuid.New(),
		Name:      "Shared Room",
		Type:      domain.ResourceTypeMeetingRoom,
		Capacity:  &capacity,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	require.NoError(t, resourceRepo.Create(ctx, resource))

	// 1. 最初の予約を作成
	startAt := time.Date(2025, 6, 1, 15, 0, 0, 0, time.UTC)
	req1 := &service.CreateReservationRequest{
		OrganizerID: user.ID,
		ResourceIDs: []uuid.UUID{resource.ID},
		Title:       "First Reservation",
		StartAt:     startAt,
		EndAt:       startAt.Add(1 * time.Hour),
		Timezone:    "Asia/Tokyo",
	}

	reservation1, err := reservationService.CreateReservation(ctx, req1)
	require.NoError(t, err)
	assert.NotNil(t, reservation1)

	// 2. 同じ時間帯に別の予約を試みる（競合）
	req2 := &service.CreateReservationRequest{
		OrganizerID: user.ID,
		ResourceIDs: []uuid.UUID{resource.ID},
		Title:       "Conflicting Reservation",
		StartAt:     startAt.Add(30 * time.Minute), // 重複する時間帯
		EndAt:       startAt.Add(90 * time.Minute),
		Timezone:    "Asia/Tokyo",
	}

	reservation2, err := reservationService.CreateReservation(ctx, req2)
	assert.Error(t, err)
	assert.Nil(t, reservation2)
	assert.Equal(t, service.ErrResourceNotAvailable, err)

	// 3. 代替リソースを検索
	alternatives, err := reservationService.FindAlternativeResources(
		ctx,
		req2.StartAt,
		req2.EndAt,
		domain.ResourceTypeMeetingRoom,
	)
	require.NoError(t, err)
	// 現在利用可能なリソースはない（resource1つのみで使用中）
	assert.Len(t, alternatives, 0)
}
