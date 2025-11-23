// backend/internal/repository/audit_log_repository_test.go
package repository_test

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/your-org/esms/internal/domain"
	"github.com/your-org/esms/internal/repository"
)

func TestAuditLogRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewAuditLogRepository(db)
	ctx := context.Background()

	log := &domain.AuditLog{
		ID:         uuid.New(),
		UserID:     uuid.New(),
		Action:     "CREATE",
		TargetType: "reservation",
		TargetID:   uuid.New().String(),
		Details:    map[string]interface{}{"status": "confirmed"},
		IPAddress:  "127.0.0.1",
		UserAgent:  "test-agent",
		CreatedAt:  time.Now(),
	}

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO audit_logs`)).
		WithArgs(
			log.ID,
			log.UserID,
			log.Action,
			log.TargetType,
			log.TargetID,
			sqlmock.AnyArg(), // Details (JSONB)
			log.IPAddress,
			log.UserAgent,
			log.CreatedAt,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Create(ctx, log)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAuditLogRepository_GetByEntityID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewAuditLogRepository(db)
	ctx := context.Background()

	entityID := uuid.New()
	expectedLog := &domain.AuditLog{
		ID:         uuid.New(),
		UserID:     uuid.New(),
		Action:     "UPDATE",
		TargetType: "reservation",
		TargetID:   entityID.String(),
		Details:    map[string]interface{}{"status": "cancelled"},
		IPAddress:  "127.0.0.1",
		UserAgent:  "test-agent",
		CreatedAt:  time.Now(),
	}

	rows := sqlmock.NewRows([]string{"id", "user_id", "action", "target_type", "target_id", "details", "ip_address", "user_agent", "created_at"}).
		AddRow(expectedLog.ID, expectedLog.UserID, expectedLog.Action, expectedLog.TargetType, expectedLog.TargetID, `{"status":"cancelled"}`, expectedLog.IPAddress, expectedLog.UserAgent, expectedLog.CreatedAt)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, user_id, action, target_type, target_id, details, ip_address, user_agent, created_at FROM audit_logs WHERE target_id = $1`)).
		WithArgs(entityID.String(), 10).
		WillReturnRows(rows)

	logs, err := repo.GetByEntityID(ctx, entityID, 10)
	assert.NoError(t, err)
	assert.Len(t, logs, 1)
	assert.Equal(t, expectedLog.ID, logs[0].ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}
