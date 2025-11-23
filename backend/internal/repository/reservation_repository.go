// backend/internal/repository/reservation_repository.go
package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/your-org/esms/internal/domain"
)

// ReservationRepository は予約データへのアクセスを提供するインターフェース
type ReservationRepository interface {
	Create(ctx context.Context, reservation *domain.Reservation) error
	CreateWithInstances(ctx context.Context, reservation *domain.Reservation, instances []*domain.ReservationInstance, resourceIDs []uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID, startAt time.Time) (*domain.Reservation, error)
	Update(ctx context.Context, reservation *domain.Reservation) error
	Delete(ctx context.Context, id uuid.UUID, startAt time.Time) error
	GetInstancesByReservationID(ctx context.Context, reservationID uuid.UUID) ([]*domain.ReservationInstance, error)
}

// postgresReservationRepository はPostgreSQLを使用したReservationRepositoryの実装
type postgresReservationRepository struct {
	db *sql.DB
}

// NewReservationRepository は新しいReservationRepositoryを作成します
func NewReservationRepository(db *sql.DB) ReservationRepository {
	return &postgresReservationRepository{db: db}
}

func (r *postgresReservationRepository) Create(ctx context.Context, reservation *domain.Reservation) error {
	query := `
		INSERT INTO reservations (id, organizer_id, title, description, start_at, end_at, rrule, is_private, timezone, approval_status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`
	_, err := r.db.ExecContext(ctx, query,
		reservation.ID,
		reservation.OrganizerID,
		reservation.Title,
		reservation.Description,
		reservation.StartAt,
		reservation.EndAt,
		reservation.RRule,
		reservation.IsPrivate,
		reservation.Timezone,
		reservation.ApprovalStatus,
		reservation.CreatedAt,
		reservation.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create reservation: %w", err)
	}
	return nil
}

// CreateWithInstances はトランザクション内で予約、予約インスタンス、リソース割り当てを作成します
func (r *postgresReservationRepository) CreateWithInstances(ctx context.Context, reservation *domain.Reservation, instances []*domain.ReservationInstance, resourceIDs []uuid.UUID) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// 予約を作成
	query := `
		INSERT INTO reservations (id, organizer_id, title, description, start_at, end_at, rrule, is_private, timezone, approval_status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`
	_, err = tx.ExecContext(ctx, query,
		reservation.ID,
		reservation.OrganizerID,
		reservation.Title,
		reservation.Description,
		reservation.StartAt,
		reservation.EndAt,
		reservation.RRule,
		reservation.IsPrivate,
		reservation.Timezone,
		reservation.ApprovalStatus,
		reservation.CreatedAt,
		reservation.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create reservation: %w", err)
	}

	// 予約インスタンスを作成
	instanceQuery := `
		INSERT INTO reservation_instances (id, reservation_id, reservation_start_at, start_at, end_at, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	for _, instance := range instances {
		_, err = tx.ExecContext(ctx, instanceQuery,
			instance.ID,
			instance.ReservationID,
			instance.ReservationStartAt,
			instance.StartAt,
			instance.EndAt,
			instance.Status,
			instance.CreatedAt,
			instance.UpdatedAt,
		)
		if err != nil {
			return fmt.Errorf("failed to create reservation instance: %w", err)
		}

		// リソース割り当てを作成
		resourceQuery := `
			INSERT INTO reservation_resources (reservation_instance_id, resource_id, created_at)
			VALUES ($1, $2, $3)
		`
		for _, resourceID := range resourceIDs {
			_, err = tx.ExecContext(ctx, resourceQuery,
				instance.ID,
				resourceID,
				time.Now(),
			)
			if err != nil {
				return fmt.Errorf("failed to create reservation resource: %w", err)
			}
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *postgresReservationRepository) GetByID(ctx context.Context, id uuid.UUID, startAt time.Time) (*domain.Reservation, error) {
	query := `
		SELECT id, organizer_id, title, description, start_at, end_at, rrule, is_private, timezone, approval_status, version, created_at, updated_at
		FROM reservations
		WHERE id = $1 AND start_at = $2
	`
	row := r.db.QueryRowContext(ctx, query, id, startAt)

	var reservation domain.Reservation
	err := row.Scan(
		&reservation.ID,
		&reservation.OrganizerID,
		&reservation.Title,
		&reservation.Description,
		&reservation.StartAt,
		&reservation.EndAt,
		&reservation.RRule,
		&reservation.IsPrivate,
		&reservation.Timezone,
		&reservation.ApprovalStatus,
		&reservation.Version,
		&reservation.CreatedAt,
		&reservation.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get reservation by id: %w", err)
	}
	return &reservation, nil
}

func (r *postgresReservationRepository) Update(ctx context.Context, reservation *domain.Reservation) error {
	reservation.UpdatedAt = time.Now()
	query := `
		UPDATE reservations
		SET title = $1, description = $2, approval_status = $3, updated_at = $4, version = version + 1
		WHERE id = $5 AND start_at = $6
	`
	result, err := r.db.ExecContext(ctx, query,
		reservation.Title,
		reservation.Description,
		reservation.ApprovalStatus,
		reservation.UpdatedAt,
		reservation.ID,
		reservation.StartAt,
	)
	if err != nil {
		return fmt.Errorf("failed to update reservation: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *postgresReservationRepository) Delete(ctx context.Context, id uuid.UUID, startAt time.Time) error {
	query := `DELETE FROM reservations WHERE id = $1 AND start_at = $2`
	result, err := r.db.ExecContext(ctx, query, id, startAt)
	if err != nil {
		return fmt.Errorf("failed to delete reservation: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *postgresReservationRepository) GetInstancesByReservationID(ctx context.Context, reservationID uuid.UUID) ([]*domain.ReservationInstance, error) {
	query := `
		SELECT id, reservation_id, reservation_start_at, start_at, end_at, status, created_at, updated_at
		FROM reservation_instances
		WHERE reservation_id = $1
		ORDER BY start_at
	`
	rows, err := r.db.QueryContext(ctx, query, reservationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get reservation instances: %w", err)
	}
	defer rows.Close()

	var instances []*domain.ReservationInstance
	for rows.Next() {
		var instance domain.ReservationInstance
		err := rows.Scan(
			&instance.ID,
			&instance.ReservationID,
			&instance.ReservationStartAt,
			&instance.StartAt,
			&instance.EndAt,
			&instance.Status,
			&instance.CreatedAt,
			&instance.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan reservation instance: %w", err)
		}
		instances = append(instances, &instance)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return instances, nil
}
