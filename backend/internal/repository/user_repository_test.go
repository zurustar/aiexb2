// backend/internal/repository/user_repository_test.go
package repository_test

import (
	"context"
	"database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/your-org/esms/internal/domain"
	"github.com/your-org/esms/internal/repository"
)

func TestUserRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewUserRepository(db)
	ctx := context.Background()

	user := &domain.User{
		ID:        uuid.New(),
		Email:     "test@example.com",
		Name:      "Test User",
		Role:      domain.RoleGeneral,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO users`)).
		WithArgs(user.ID, user.Email, user.Name, user.Role, user.CreatedAt, user.UpdatedAt).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Create(ctx, user)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_GetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewUserRepository(db)
	ctx := context.Background()

	expectedUser := &domain.User{
		ID:        uuid.New(),
		Email:     "test@example.com",
		Name:      "Test User",
		Role:      domain.RoleGeneral,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	rows := sqlmock.NewRows([]string{"id", "email", "name", "role", "created_at", "updated_at"}).
		AddRow(expectedUser.ID, expectedUser.Email, expectedUser.Name, expectedUser.Role, expectedUser.CreatedAt, expectedUser.UpdatedAt)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, email, name, role, created_at, updated_at FROM users WHERE id = $1`)).
		WithArgs(expectedUser.ID).
		WillReturnRows(rows)

	user, err := repo.GetByID(ctx, expectedUser.ID)
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_GetByID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewUserRepository(db)
	ctx := context.Background()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, email, name, role, created_at, updated_at FROM users WHERE id = $1`)).
		WithArgs(sqlmock.AnyArg()).
		WillReturnError(sql.ErrNoRows)

	user, err := repo.GetByID(ctx, uuid.New())
	assert.Error(t, err)
	assert.Equal(t, repository.ErrNotFound, err)
	assert.Nil(t, user)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewUserRepository(db)
	ctx := context.Background()

	user := &domain.User{
		ID:    uuid.New(),
		Email: "updated@example.com",
		Name:  "Updated User",
		Role:  domain.RoleManager,
	}

	// Update内部でtime.Now()が呼ばれるため、AnyArgを使用するか、実装側で時刻を受け取るようにするか。
	// ここではAnyArgを使用する。
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE users SET email = $1, name = $2, role = $3, updated_at = $4 WHERE id = $5`)).
		WithArgs(user.Email, user.Name, user.Role, sqlmock.AnyArg(), user.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Update(ctx, user)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewUserRepository(db)
	ctx := context.Background()

	id := uuid.New()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM users WHERE id = $1`)).
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Delete(ctx, id)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
