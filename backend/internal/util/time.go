// backend/internal/util/time.go
package util

import (
	"time"
)

// TimeRange は時間の範囲を表す構造体
type TimeRange struct {
	Start time.Time
	End   time.Time
}

// JST は日本標準時のロケーション
var JST *time.Location

func init() {
	var err error
	JST, err = time.LoadLocation("Asia/Tokyo")
	if err != nil {
		// フォールバック: 固定オフセット
		JST = time.FixedZone("Asia/Tokyo", 9*60*60)
	}
}

// NowJST は現在の日本時間を返します
func NowJST() time.Time {
	return time.Now().In(JST)
}

// ToJST は指定された時間を日本時間に変換します
func ToJST(t time.Time) time.Time {
	return t.In(JST)
}

// IsOverlapping は2つの時間範囲が重複しているかを判定します
// start1 <= end1, start2 <= end2 であることを前提とします
// 接しているだけ（end1 == start2）の場合は重複とみなしません
func IsOverlapping(start1, end1, start2, end2 time.Time) bool {
	return start1.Before(end2) && end1.After(start2)
}

// IsBusinessHour は指定された時間が営業時間内（9:00 - 18:00, 平日）かを判定します
func IsBusinessHour(t time.Time) bool {
	local := ToJST(t)

	// 土日は営業時間外
	weekday := local.Weekday()
	if weekday == time.Saturday || weekday == time.Sunday {
		return false
	}

	hour := local.Hour()
	// 9:00 <= t < 18:00
	return hour >= 9 && hour < 18
}

// TruncateToMinute は秒以下を切り捨てます
func TruncateToMinute(t time.Time) time.Time {
	return t.Truncate(time.Minute)
}
