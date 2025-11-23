// backend/internal/repository/resource_repository_test.go
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

func TestResourceRepository_FindAvailable(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewResourceRepository(db)
	ctx := context.Background()

	startAt := time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)
	endAt := startAt.Add(1 * time.Hour)

	capacity := 10
	expectedResource := &domain.Resource{
		ID:        uuid.New(),
		Name:      "Meeting Room A",
		Type:      domain.ResourceTypeMeetingRoom,
		Capacity:  &capacity,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	rows := sqlmock.NewRows([]string{"id", "name", "type", "capacity", "created_at", "updated_at"}).
		AddRow(expectedResource.ID, expectedResource.Name, expectedResource.Type, expectedResource.Capacity, expectedResource.CreatedAt, expectedResource.UpdatedAt)

	// クエリのマッチング
	// NOT EXISTS 句を含むクエリが正しく発行されるか確認
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT r.id, r.name, r.type, r.capacity, r.created_at, r.updated_at FROM resources r WHERE NOT EXISTS`)).
		WithArgs(startAt, endAt).
		WillReturnRows(rows)

	resources, err := repo.FindAvailable(ctx, startAt, endAt)
	assert.NoError(t, err)
	assert.Len(t, resources, 1)
	assert.Equal(t, expectedResource, resources[0])
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestResourceRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewResourceRepository(db)
	ctx := context.Background()

	capacity2 := 5
	resource := &domain.Resource{
		ID:        uuid.New(),
		Name:      "Room B",
		Type:      domain.ResourceTypeMeetingRoom,
		Capacity:  &capacity2,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO resources`)).
		WithArgs(resource.ID, resource.Name, resource.Type, resource.Capacity, resource.CreatedAt, resource.UpdatedAt).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Create(ctx, resource)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
