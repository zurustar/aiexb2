// backend/internal/domain/reservation_test.go
package domain_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/your-org/esms/internal/domain"
)

func TestReservation_Validate(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name        string
		reservation domain.Reservation
		wantErr     bool
	}{
		{
			name: "Valid single reservation",
			reservation: domain.Reservation{
				Title:   "Meeting",
				StartAt: now,
				EndAt:   now.Add(1 * time.Hour),
			},
			wantErr: false,
		},
		{
			name: "Valid recurring reservation",
			reservation: domain.Reservation{
				Title:   "Daily Standup",
				StartAt: now,
				EndAt:   now.Add(15 * time.Minute),
				RRule:   "FREQ=DAILY;COUNT=5",
			},
			wantErr: false,
		},
		{
			name: "Missing title",
			reservation: domain.Reservation{
				StartAt: now,
				EndAt:   now.Add(1 * time.Hour),
			},
			wantErr: true,
		},
		{
			name: "End time before start time",
			reservation: domain.Reservation{
				Title:   "Invalid Time",
				StartAt: now.Add(1 * time.Hour),
				EndAt:   now,
			},
			wantErr: true,
		},
		{
			name: "Invalid RRULE",
			reservation: domain.Reservation{
				Title:   "Invalid RRULE",
				StartAt: now,
				EndAt:   now.Add(1 * time.Hour),
				RRule:   "INVALID_RRULE",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.reservation.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestReservation_ExpandInstances(t *testing.T) {
	// 固定の日時を設定してテストの安定性を確保
	baseTime := time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)
	reservationID := uuid.New()

	tests := []struct {
		name          string
		reservation   domain.Reservation
		queryStart    time.Time
		queryEnd      time.Time
		wantCount     int
		checkInstance func(*testing.T, []domain.ReservationInstance)
	}{
		{
			name: "Single reservation within range",
			reservation: domain.Reservation{
				ID:      reservationID,
				Title:   "Single",
				StartAt: baseTime,
				EndAt:   baseTime.Add(1 * time.Hour),
			},
			queryStart: baseTime.Add(-1 * time.Hour),
			queryEnd:   baseTime.Add(2 * time.Hour),
			wantCount:  1,
			checkInstance: func(t *testing.T, instances []domain.ReservationInstance) {
				assert.Equal(t, reservationID, instances[0].ReservationID)
				assert.Equal(t, baseTime, instances[0].StartAt)
				assert.Equal(t, baseTime.Add(1*time.Hour), instances[0].EndAt)
				// パーティションキーの確認
				assert.Equal(t, baseTime, instances[0].ReservationStartAt)
			},
		},
		{
			name: "Single reservation outside range",
			reservation: domain.Reservation{
				ID:      reservationID,
				Title:   "Single",
				StartAt: baseTime,
				EndAt:   baseTime.Add(1 * time.Hour),
			},
			queryStart:    baseTime.Add(2 * time.Hour),
			queryEnd:      baseTime.Add(3 * time.Hour),
			wantCount:     0,
			checkInstance: nil,
		},
		{
			name: "Daily recurring reservation (5 days)",
			reservation: domain.Reservation{
				ID:      reservationID,
				Title:   "Daily",
				StartAt: baseTime,
				EndAt:   baseTime.Add(1 * time.Hour),
				RRule:   "FREQ=DAILY;COUNT=5",
			},
			queryStart: baseTime,
			queryEnd:   baseTime.Add(5 * 24 * time.Hour), // 5日間
			wantCount:  5,
			checkInstance: func(t *testing.T, instances []domain.ReservationInstance) {
				for i, instance := range instances {
					expectedStart := baseTime.Add(time.Duration(i) * 24 * time.Hour)
					assert.Equal(t, expectedStart, instance.StartAt)
					assert.Equal(t, expectedStart.Add(1*time.Hour), instance.EndAt)
					// 親の開始日時（パーティションキー）は常に一定
					assert.Equal(t, baseTime, instance.ReservationStartAt)
				}
			},
		},
		{
			name: "Weekly recurring reservation (subset)",
			reservation: domain.Reservation{
				ID:      reservationID,
				Title:   "Weekly",
				StartAt: baseTime,
				EndAt:   baseTime.Add(1 * time.Hour),
				RRule:   "FREQ=WEEKLY;COUNT=4", // 4週間
			},
			queryStart: baseTime.Add(1 * 24 * time.Hour),  // 2日目から
			queryEnd:   baseTime.Add(15 * 24 * time.Hour), // 15日目まで（2週分含まれるはず）
			wantCount:  2,                                 // 2週目と3週目
			checkInstance: func(t *testing.T, instances []domain.ReservationInstance) {
				// 1週後
				assert.Equal(t, baseTime.Add(7*24*time.Hour), instances[0].StartAt)
				// 2週後
				assert.Equal(t, baseTime.Add(14*24*time.Hour), instances[1].StartAt)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			instances, err := tt.reservation.ExpandInstances(tt.queryStart, tt.queryEnd)
			assert.NoError(t, err)
			assert.Len(t, instances, tt.wantCount)
			if tt.checkInstance != nil && len(instances) > 0 {
				tt.checkInstance(t, instances)
			}
		})
	}
}
