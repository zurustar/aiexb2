// backend/internal/domain/reservation.go
package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/teambition/rrule-go"
)

// ApprovalStatus は予約の承認状態を表す型
type ApprovalStatus string

const (
	ApprovalStatusPending   ApprovalStatus = "PENDING"   // 承認待ち
	ApprovalStatusConfirmed ApprovalStatus = "CONFIRMED" // 確定済み
	ApprovalStatusRejected  ApprovalStatus = "REJECTED"  // 却下
)

// ReservationStatus は予約インスタンスのステータスを表す型
type ReservationStatus string

const (
	ReservationStatusConfirmed ReservationStatus = "CONFIRMED"  // 確定
	ReservationStatusCancelled ReservationStatus = "CANCELLED"  // キャンセル
	ReservationStatusCheckedIn ReservationStatus = "CHECKED_IN" // チェックイン済み
	ReservationStatusCompleted ReservationStatus = "COMPLETED"  // 完了
	ReservationStatusNoShow    ReservationStatus = "NO_SHOW"    // 無断キャンセル
)

// Reservation は予約（親）エンティティを表す構造体
type Reservation struct {
	ID             uuid.UUID
	OrganizerID    uuid.UUID
	Title          string
	Description    string
	StartAt        time.Time
	EndAt          time.Time
	RRule          string // iCalendar RFC 5545
	IsPrivate      bool
	Timezone       string
	ApprovalStatus ApprovalStatus
	UpdatedBy      *uuid.UUID
	Version        int
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      *time.Time

	// Relations
	Organizer *User
}

// ReservationInstance は予約インスタンス（展開後）を表す構造体
type ReservationInstance struct {
	ID                 uuid.UUID
	ReservationID      uuid.UUID
	ReservationStartAt time.Time // 親予約の開始日時（パーティションキー）
	StartAt            time.Time
	EndAt              time.Time
	OriginalStartAt    *time.Time // 繰り返し例外時の元の開始日時
	Status             ReservationStatus
	CheckedInAt        *time.Time
	CreatedAt          time.Time
	UpdatedAt          time.Time

	// Relations
	Reservation  *Reservation
	Resources    []*Resource
	Participants []*User
}

// IsRecurring は繰り返し予約かどうかを判定します
func (r *Reservation) IsRecurring() bool {
	return r.RRule != ""
}

// Validate は予約の整合性を検証します
func (r *Reservation) Validate() error {
	if r.Title == "" {
		return errors.New("title is required")
	}
	if r.StartAt.After(r.EndAt) {
		return errors.New("start time must be before end time")
	}
	if r.IsRecurring() {
		// RRULEの簡易検証
		if _, err := rrule.StrToRRule(r.RRule); err != nil {
			return errors.New("invalid rrule format")
		}
	}
	return nil
}

// ExpandInstances は指定された期間内のインスタンスを展開します
func (r *Reservation) ExpandInstances(start, end time.Time) ([]ReservationInstance, error) {
	if !r.IsRecurring() {
		// 単発予約の場合
		if r.StartAt.Before(end) && r.EndAt.After(start) {
			return []ReservationInstance{{
				ID:                 uuid.New(),
				ReservationID:      r.ID,
				ReservationStartAt: r.StartAt,
				StartAt:            r.StartAt,
				EndAt:              r.EndAt,
				Status:             ReservationStatusConfirmed,
			}}, nil
		}
		return []ReservationInstance{}, nil
	}

	// 繰り返し予約の展開
	rule, err := rrule.StrToRRule(r.RRule)
	if err != nil {
		return nil, err
	}

	// RRULEの開始時間を設定
	rule.DTStart(r.StartAt)

	// 指定期間内の日時を取得
	// rrule-go の Between は start <= time < end
	dates := rule.Between(start, end, true)

	duration := r.EndAt.Sub(r.StartAt)
	instances := make([]ReservationInstance, 0, len(dates))

	for _, date := range dates {
		instances = append(instances, ReservationInstance{
			ID:                 uuid.New(),
			ReservationID:      r.ID,
			ReservationStartAt: r.StartAt,
			StartAt:            date,
			EndAt:              date.Add(duration),
			Status:             ReservationStatusConfirmed,
		})
	}

	return instances, nil
}
