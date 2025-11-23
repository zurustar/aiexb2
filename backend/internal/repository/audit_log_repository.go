// backend/internal/repository/audit_log_repository.go
package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/your-org/esms/internal/domain"
)

// AuditLogRepository は監査ログデータへのアクセスを提供するインターフェース
type AuditLogRepository interface {
	Create(ctx context.Context, log *domain.AuditLog) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.AuditLog, error)
	GetByEntityID(ctx context.Context, entityID uuid.UUID, limit int) ([]*domain.AuditLog, error)
}

// postgresAuditLogRepository はPostgreSQLを使用したAuditLogRepositoryの実装
type postgresAuditLogRepository struct {
	db *sql.DB
}

// NewAuditLogRepository は新しいAuditLogRepositoryを作成します
func NewAuditLogRepository(db *sql.DB) AuditLogRepository {
	return &postgresAuditLogRepository{db: db}
}

func (r *postgresAuditLogRepository) Create(ctx context.Context, log *domain.AuditLog) error {
	// DetailsをJSONに変換
	detailsJSON, err := json.Marshal(log.Details)
	if err != nil {
		return fmt.Errorf("failed to marshal details: %w", err)
	}

	query := `
		INSERT INTO audit_logs (id, user_id, action, target_type, target_id, details, ip_address, user_agent, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err = r.db.ExecContext(ctx, query,
		log.ID,
		log.UserID,
		log.Action,
		log.TargetType,
		log.TargetID,
		detailsJSON,
		log.IPAddress,
		log.UserAgent,
		log.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create audit log: %w", err)
	}
	return nil
}

func (r *postgresAuditLogRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.AuditLog, error) {
	query := `
		SELECT id, user_id, action, target_type, target_id, details, ip_address, user_agent, created_at
		FROM audit_logs
		WHERE id = $1
	`
	row := r.db.QueryRowContext(ctx, query, id)

	var log domain.AuditLog
	var detailsJSON []byte
	err := row.Scan(
		&log.ID,
		&log.UserID,
		&log.Action,
		&log.TargetType,
		&log.TargetID,
		&detailsJSON,
		&log.IPAddress,
		&log.UserAgent,
		&log.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get audit log by id: %w", err)
	}

	// JSONをDetailsにデシリアライズ
	if err := json.Unmarshal(detailsJSON, &log.Details); err != nil {
		return nil, fmt.Errorf("failed to unmarshal details: %w", err)
	}

	return &log, nil
}

func (r *postgresAuditLogRepository) GetByEntityID(ctx context.Context, entityID uuid.UUID, limit int) ([]*domain.AuditLog, error) {
	query := `
		SELECT id, user_id, action, target_type, target_id, details, ip_address, user_agent, created_at
		FROM audit_logs
		WHERE target_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	`
	// entityIDはUUIDだが、target_idはVARCHARなので文字列に変換
	rows, err := r.db.QueryContext(ctx, query, entityID.String(), limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get audit logs by entity id: %w", err)
	}
	defer rows.Close()

	var logs []*domain.AuditLog
	for rows.Next() {
		var log domain.AuditLog
		var detailsJSON []byte
		err := rows.Scan(
			&log.ID,
			&log.UserID,
			&log.Action,
			&log.TargetType,
			&log.TargetID,
			&detailsJSON,
			&log.IPAddress,
			&log.UserAgent,
			&log.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan audit log: %w", err)
		}

		// JSONをDetailsにデシリアライズ
		if err := json.Unmarshal(detailsJSON, &log.Details); err != nil {
			return nil, fmt.Errorf("failed to unmarshal details: %w", err)
		}

		logs = append(logs, &log)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return logs, nil
}
