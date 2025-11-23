// backend/tests/integration/repository_test.go
package integration_test

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/your-org/esms/internal/domain"
	"github.com/your-org/esms/internal/repository"
)

var testDB *sql.DB

// TestMain はテストスイートの初期化と終了処理を行います
func TestMain(m *testing.M) {
	// 統合テスト用のDB接続文字列を環境変数から取得
	// 例: TEST_DATABASE_URL=postgres://user:pass@localhost:5432/esms_test?sslmode=disable
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		// 環境変数が設定されていない場合はスキップ
		os.Exit(0)
	}

	var err error
	testDB, err = sql.Open("pgx", dsn)
	if err != nil {
		panic("failed to connect to test database: " + err.Error())
	}

	// DB接続確認
	if err := testDB.Ping(); err != nil {
		panic("failed to ping test database: " + err.Error())
	}

	// テスト実行
	code := m.Run()

	// クリーンアップ
	testDB.Close()

	os.Exit(code)
}

// cleanupDB はテストデータをクリーンアップします
func cleanupDB(t *testing.T) {
	_, err := testDB.Exec("TRUNCATE TABLE audit_logs, reservation_resources, reservation_participants, reservation_instances, reservations_2025, reservations_2026, reservations_2027, resources, users CASCADE")
	require.NoError(t, err)
}

func TestUserRepository_Integration(t *testing.T) {
	if testDB == nil {
		t.Skip("TEST_DATABASE_URL not set, skipping integration test")
	}

	cleanupDB(t)
	defer cleanupDB(t)

	repo := repository.NewUserRepository(testDB)
	ctx := context.Background()

	// Create
	user := &domain.User{
		ID:        uuid.New(),
		Email:     "integration@example.com",
		Name:      "Integration Test User",
		Role:      domain.RoleGeneral,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := repo.Create(ctx, user)
	require.NoError(t, err)

	// GetByID
	retrieved, err := repo.GetByID(ctx, user.ID)
	require.NoError(t, err)
	assert.Equal(t, user.Email, retrieved.Email)
	assert.Equal(t, user.Name, retrieved.Name)
	assert.Equal(t, user.Role, retrieved.Role)

	// GetByEmail
	byEmail, err := repo.GetByEmail(ctx, user.Email)
	require.NoError(t, err)
	assert.Equal(t, user.ID, byEmail.ID)

	// Update
	retrieved.Name = "Updated Name"
	err = repo.Update(ctx, retrieved)
	require.NoError(t, err)

	updated, err := repo.GetByID(ctx, user.ID)
	require.NoError(t, err)
	assert.Equal(t, "Updated Name", updated.Name)

	// Delete
	err = repo.Delete(ctx, user.ID)
	require.NoError(t, err)

	_, err = repo.GetByID(ctx, user.ID)
	assert.Equal(t, repository.ErrNotFound, err)
}

func TestResourceRepository_Integration(t *testing.T) {
	if testDB == nil {
		t.Skip("TEST_DATABASE_URL not set, skipping integration test")
	}

	cleanupDB(t)
	defer cleanupDB(t)

	repo := repository.NewResourceRepository(testDB)
	ctx := context.Background()

	// Create
	capacity := 10
	resource := &domain.Resource{
		ID:        uuid.New(),
		Name:      "Integration Test Room",
		Type:      domain.ResourceTypeMeetingRoom,
		Capacity:  &capacity,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := repo.Create(ctx, resource)
	require.NoError(t, err)

	// GetByID
	retrieved, err := repo.GetByID(ctx, resource.ID)
	require.NoError(t, err)
	assert.Equal(t, resource.Name, retrieved.Name)
	assert.Equal(t, resource.Type, retrieved.Type)

	// FindAvailable (空き時間検索)
	startAt := time.Now()
	endAt := startAt.Add(1 * time.Hour)
	available, err := repo.FindAvailable(ctx, startAt, endAt)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(available), 1)

	// Delete
	err = repo.Delete(ctx, resource.ID)
	require.NoError(t, err)
}

func TestAuditLogRepository_Integration(t *testing.T) {
	if testDB == nil {
		t.Skip("TEST_DATABASE_URL not set, skipping integration test")
	}

	cleanupDB(t)
	defer cleanupDB(t)

	repo := repository.NewAuditLogRepository(testDB)
	ctx := context.Background()

	// Create
	log := &domain.AuditLog{
		ID:         uuid.New(),
		UserID:     uuid.New(),
		Action:     domain.AuditActionCreate,
		TargetType: "reservation",
		TargetID:   uuid.New().String(),
		Details:    map[string]interface{}{"status": "confirmed", "count": 1},
		IPAddress:  "192.168.1.1",
		UserAgent:  "Integration-Test-Agent",
		CreatedAt:  time.Now(),
	}

	err := repo.Create(ctx, log)
	require.NoError(t, err)

	// GetByID
	retrieved, err := repo.GetByID(ctx, log.ID)
	require.NoError(t, err)
	assert.Equal(t, log.Action, retrieved.Action)
	assert.Equal(t, log.TargetType, retrieved.TargetType)
	assert.Equal(t, log.TargetID, retrieved.TargetID)
	assert.Equal(t, "confirmed", retrieved.Details["status"])

	// GetByEntityID
	entityID, _ := uuid.Parse(log.TargetID)
	logs, err := repo.GetByEntityID(ctx, entityID, 10)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(logs), 1)
}

func TestReservationRepository_Integration(t *testing.T) {
	if testDB == nil {
		t.Skip("TEST_DATABASE_URL not set, skipping integration test")
	}

	cleanupDB(t)
	defer cleanupDB(t)

	// まずユーザーとリソースを作成
	userRepo := repository.NewUserRepository(testDB)
	resourceRepo := repository.NewResourceRepository(testDB)
	reservationRepo := repository.NewReservationRepository(testDB)
	ctx := context.Background()

	user := &domain.User{
		ID:        uuid.New(),
		Email:     "organizer@example.com",
		Name:      "Organizer",
		Role:      domain.RoleGeneral,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	require.NoError(t, userRepo.Create(ctx, user))

	capacity := 5
	resource := &domain.Resource{
		ID:        uuid.New(),
		Name:      "Test Room",
		Type:      domain.ResourceTypeMeetingRoom,
		Capacity:  &capacity,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	require.NoError(t, resourceRepo.Create(ctx, resource))

	// 予約作成
	startAt := time.Date(2025, 6, 1, 10, 0, 0, 0, time.UTC)
	reservation := &domain.Reservation{
		ID:             uuid.New(),
		OrganizerID:    user.ID,
		Title:          "Integration Test Meeting",
		Description:    "Test",
		StartAt:        startAt,
		EndAt:          startAt.Add(1 * time.Hour),
		Timezone:       "Asia/Tokyo",
		ApprovalStatus: domain.ApprovalStatusConfirmed,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	instance := &domain.ReservationInstance{
		ID:                 uuid.New(),
		ReservationID:      reservation.ID,
		ReservationStartAt: reservation.StartAt,
		StartAt:            reservation.StartAt,
		EndAt:              reservation.EndAt,
		Status:             domain.ReservationStatusConfirmed,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	err := reservationRepo.CreateWithInstances(ctx, reservation, []*domain.ReservationInstance{instance}, []uuid.UUID{resource.ID})
	require.NoError(t, err)

	// GetByID
	retrieved, err := reservationRepo.GetByID(ctx, reservation.ID, reservation.StartAt)
	require.NoError(t, err)
	assert.Equal(t, reservation.Title, retrieved.Title)
	assert.Equal(t, reservation.OrganizerID, retrieved.OrganizerID)

	// GetInstancesByReservationID
	instances, err := reservationRepo.GetInstancesByReservationID(ctx, reservation.ID)
	require.NoError(t, err)
	assert.Len(t, instances, 1)
	assert.Equal(t, instance.Status, instances[0].Status)
}
