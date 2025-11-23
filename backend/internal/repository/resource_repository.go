// backend/internal/repository/resource_repository.go
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

// ResourceRepository はリソースデータへのアクセスを提供するインターフェース
type ResourceRepository interface {
	Create(ctx context.Context, resource *domain.Resource) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Resource, error)
	Update(ctx context.Context, resource *domain.Resource) error
	Delete(ctx context.Context, id uuid.UUID) error
	FindAvailable(ctx context.Context, startAt, endAt time.Time) ([]*domain.Resource, error)
}

// postgresResourceRepository はPostgreSQLを使用したResourceRepositoryの実装
type postgresResourceRepository struct {
	db *sql.DB
}

// NewResourceRepository は新しいResourceRepositoryを作成します
func NewResourceRepository(db *sql.DB) ResourceRepository {
	return &postgresResourceRepository{db: db}
}

func (r *postgresResourceRepository) Create(ctx context.Context, resource *domain.Resource) error {
	query := `
		INSERT INTO resources (id, name, type, capacity, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.ExecContext(ctx, query,
		resource.ID,
		resource.Name,
		resource.Type,
		resource.Capacity,
		resource.CreatedAt,
		resource.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create resource: %w", err)
	}
	return nil
}

func (r *postgresResourceRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Resource, error) {
	query := `
		SELECT id, name, type, capacity, created_at, updated_at
		FROM resources
		WHERE id = $1
	`
	row := r.db.QueryRowContext(ctx, query, id)

	var resource domain.Resource
	err := row.Scan(
		&resource.ID,
		&resource.Name,
		&resource.Type,
		&resource.Capacity,
		&resource.CreatedAt,
		&resource.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get resource by id: %w", err)
	}
	return &resource, nil
}

func (r *postgresResourceRepository) Update(ctx context.Context, resource *domain.Resource) error {
	resource.UpdatedAt = time.Now()
	query := `
		UPDATE resources
		SET name = $1, type = $2, capacity = $3, updated_at = $4
		WHERE id = $5
	`
	result, err := r.db.ExecContext(ctx, query,
		resource.Name,
		resource.Type,
		resource.Capacity,
		resource.UpdatedAt,
		resource.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update resource: %w", err)
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

func (r *postgresResourceRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM resources WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete resource: %w", err)
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

// FindAvailable は指定された期間に空いているリソースを取得します
// 重複する予約が存在しないリソースを返します
func (r *postgresResourceRepository) FindAvailable(ctx context.Context, startAt, endAt time.Time) ([]*domain.Resource, error) {
	// 指定期間に重複する予約があるリソースを除外するクエリ
	// reservation_instances テーブルを使用
	// 重複条件: (start < endAt AND end > startAt)
	query := `
		SELECT r.id, r.name, r.type, r.capacity, r.created_at, r.updated_at
		FROM resources r
		WHERE NOT EXISTS (
			SELECT 1
			FROM reservation_instances ri
			WHERE ri.resource_id = r.id
			  AND ri.start_at < $2
			  AND ri.end_at > $1
		)
		ORDER BY r.name
	`

	rows, err := r.db.QueryContext(ctx, query, startAt, endAt)
	if err != nil {
		return nil, fmt.Errorf("failed to find available resources: %w", err)
	}
	defer rows.Close()

	var resources []*domain.Resource
	for rows.Next() {
		var r domain.Resource
		err := rows.Scan(
			&r.ID,
			&r.Name,
			&r.Type,
			&r.Capacity,
			&r.CreatedAt,
			&r.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan resource: %w", err)
		}
		resources = append(resources, &r)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return resources, nil
}
