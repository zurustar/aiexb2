// backend/internal/util/time_test.go
package util_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/your-org/esms/internal/util"
)

func TestIsOverlapping(t *testing.T) {
	base := time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)

	tests := []struct {
		name   string
		start1 time.Time
		end1   time.Time
		start2 time.Time
		end2   time.Time
		want   bool
	}{
		{
			name:   "Completely overlapping",
			start1: base,
			end1:   base.Add(2 * time.Hour),
			start2: base,
			end2:   base.Add(2 * time.Hour),
			want:   true,
		},
		{
			name:   "Partial overlap (start)",
			start1: base,
			end1:   base.Add(2 * time.Hour),
			start2: base.Add(1 * time.Hour),
			end2:   base.Add(3 * time.Hour),
			want:   true,
		},
		{
			name:   "Partial overlap (end)",
			start1: base.Add(1 * time.Hour),
			end1:   base.Add(3 * time.Hour),
			start2: base,
			end2:   base.Add(2 * time.Hour),
			want:   true,
		},
		{
			name:   "Included",
			start1: base,
			end1:   base.Add(3 * time.Hour),
			start2: base.Add(1 * time.Hour),
			end2:   base.Add(2 * time.Hour),
			want:   true,
		},
		{
			name:   "Touching (not overlapping)",
			start1: base,
			end1:   base.Add(1 * time.Hour),
			start2: base.Add(1 * time.Hour),
			end2:   base.Add(2 * time.Hour),
			want:   false,
		},
		{
			name:   "Separate",
			start1: base,
			end1:   base.Add(1 * time.Hour),
			start2: base.Add(2 * time.Hour),
			end2:   base.Add(3 * time.Hour),
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, util.IsOverlapping(tt.start1, tt.end1, tt.start2, tt.end2))
		})
	}
}

func TestIsBusinessHour(t *testing.T) {
	// 2025-11-24 (Mon) - 平日
	monday := time.Date(2025, 11, 24, 0, 0, 0, 0, util.JST)
	// 2025-11-23 (Sun) - 休日
	sunday := time.Date(2025, 11, 23, 0, 0, 0, 0, util.JST)

	tests := []struct {
		name string
		time time.Time
		want bool
	}{
		{"Monday 8:59", monday.Add(8*time.Hour + 59*time.Minute), false},
		{"Monday 9:00", monday.Add(9 * time.Hour), true},
		{"Monday 12:00", monday.Add(12 * time.Hour), true},
		{"Monday 17:59", monday.Add(17*time.Hour + 59*time.Minute), true},
		{"Monday 18:00", monday.Add(18 * time.Hour), false},
		{"Sunday 12:00", sunday.Add(12 * time.Hour), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, util.IsBusinessHour(tt.time))
		})
	}
}
