package helpers

import (
	"testing"
	"time"
)

func Equal[v comparable](t *testing.T, got, expected v) {
	t.Helper()

	if expected != got {
		t.Errorf(`assert.Equal(
		t,
		got:
		%v,
		expected:
		%v
		)`, got, expected)
	}
}

func TestNewErrorStackTraceString(t *testing.T) {
	got := NewErrorStackTraceString("test message")
	expected := "helpers/helpers_test.go:23 - helpers.TestNewErrorStackTraceString - test message \n from"
	Equal(t, got, expected)
}

var nowFunc = time.Now
var mockNow = func() time.Time {
	return time.Date(2024, 11, 06, 0, 0, 0, 0, time.UTC)
}

func TestGetCurrentWeekDate(t *testing.T) {
	originalNowFunc := nowFunc
	nowFunc = mockNow
	defer func() { nowFunc = originalNowFunc }()

	expected := "2024-11-04"

	got := GetCurrentWeekDate()
	Equal(t, got, expected)
}

func TestGetStartOfWeek(t *testing.T) {
	testCases := []struct {
		name          string
		date          string
		startOfWeek   string
		expectedError string
	}{
		{
			name:          "middle of the week",
			date:          "2024-11-06",
			startOfWeek:   "2024-11-04",
			expectedError: "",
		},
		{
			name:          "start of the week",
			date:          "2024-11-04",
			startOfWeek:   "2024-11-04",
			expectedError: "",
		},
		{
			name:          "end of the week",
			date:          "2024-11-10",
			startOfWeek:   "2024-11-04",
			expectedError: "",
		},
		{
			name:          "invalid date format",
			date:          "30-11-2024",
			startOfWeek:   "",
			expectedError: "helpers/helpers.go:39 - helpers.GetStartOfWeek - failed to parse date \n from parsing time \"30-11-2024\" as \"2006-01-02\": cannot parse \"30-11-2024\" as \"2006\"",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := GetStartOfWeek(tc.date)
			if tc.expectedError != "" {
				Equal(t, err.Error(), tc.expectedError)
			} else {
				Equal(t, got, tc.startOfWeek)
			}
		})
	}
}
