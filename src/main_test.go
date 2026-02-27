package main

import (
	"strings"
	"testing"
	"time"
)

func TestGetCurrentDayWithTime(t *testing.T) {
	tests := []struct {
		name        string
		startDate   string
		inputTime   time.Time
		expectedDay int
	}{
		{
			name:        "Day 1 - Aug 13 2025",
			startDate:   "2025-08-13",
			inputTime:   time.Date(2025, 8, 13, 21, 0, 0, 0, time.FixedZone("Bangkok", 7*60*60)),
			expectedDay: 1,
		},
		{
			name:        "Day 2 - Aug 14 2025",
			startDate:   "2025-08-13",
			inputTime:   time.Date(2025, 8, 14, 21, 0, 0, 0, time.FixedZone("Bangkok", 7*60*60)),
			expectedDay: 2,
		},
		{
			name:        "Day 199 - Feb 27 2026",
			startDate:   "2025-08-13",
			inputTime:   time.Date(2026, 2, 27, 21, 0, 0, 0, time.FixedZone("Bangkok", 7*60*60)),
			expectedDay: 199,
		},
		{
			name:        "Day 200 - Feb 28 2026",
			startDate:   "2025-08-13",
			inputTime:   time.Date(2026, 2, 28, 21, 0, 0, 0, time.FixedZone("Bangkok", 7*60*60)),
			expectedDay: 200,
		},
		{
			name:        "Day 100 - Nov 20 2025",
			startDate:   "2025-08-13",
			inputTime:   time.Date(2025, 11, 20, 21, 0, 0, 0, time.FixedZone("Bangkok", 7*60*60)),
			expectedDay: 100,
		},
		{
			name:        "Exactly at midnight Feb 28 2026",
			startDate:   "2025-08-13",
			inputTime:   time.Date(2026, 2, 28, 0, 0, 0, 0, time.FixedZone("Bangkok", 7*60*60)),
			expectedDay: 200, // Midnight Feb 28 is already day 200
		},
		{
			name:        "Day 200 at midnight Feb 28 2026",
			startDate:   "2025-08-13",
			inputTime:   time.Date(2026, 2, 28, 0, 0, 1, 0, time.FixedZone("Bangkok", 7*60*60)),
			expectedDay: 200,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			day, err := GetCurrentDayWithTime(tt.startDate, tt.inputTime)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if day != tt.expectedDay {
				t.Errorf("expected day %d, got %d", tt.expectedDay, day)
			}
		})
	}
}

func TestSplitMessage(t *testing.T) {
	tests := []struct {
		name     string
		msg      string
		maxLen   int
		expected int
	}{
		{"short message", "hello", 10, 1},
		{"exact length", "hello", 5, 1},
		{"long message", strings.Repeat("a", 50), 20, 3},
		{"empty message", "", 10, 1},
		{"unicode Thai", "สวัสดี", 3, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := splitMessage(tt.msg, tt.maxLen)
			if len(result) != tt.expected {
				t.Errorf("expected %d chunks, got %d", tt.expected, len(result))
			}
		})
	}
}
