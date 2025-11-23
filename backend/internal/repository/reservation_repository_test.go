// backend/internal/repository/reservation_repository_test.go
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

func TestReservationRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewReservationRepository(db)
	ctx := context.Background()

	reservation := &domain.Reservation{
		ID:             uuid.New(),
		OrganizerID:    uuid.New(),
		Title:          "Team Meeting",
		Description:    "Weekly sync",
		StartAt:        time.Now(),
		EndAt:          time.Now().Add(1 * time.Hour),
		ApprovalStatus: domain.ApprovalStatusConfirmed,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO reservations`)).
		WithArgs(
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
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Create(ctx, reservation)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestReservationRepository_CreateWithInstances(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewReservationRepository(db)
	ctx := context.Background()

	reservation := &domain.Reservation{
		ID:             uuid.New(),
		OrganizerID:    uuid.New(),
		Title:          "Team Meeting",
		StartAt:        time.Now(),
		EndAt:          time.Now().Add(1 * time.Hour),
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

	resourceID := uuid.New()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO reservations`)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO reservation_instances`)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO reservation_resources`)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = repo.CreateWithInstances(ctx, reservation, []*domain.ReservationInstance{instance}, []uuid.UUID{resourceID})
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
