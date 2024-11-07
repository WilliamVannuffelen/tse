package helpers

import (
	"strings"
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

func Contains(t *testing.T, got, expected string) {
	t.Helper()

	if !strings.Contains(got, expected) {
		t.Errorf(`assert.Contains(
		t,
		got:
		-%v-,
		should be present in:
		-%v-
		)`, expected, got)
	}
}

func createNowMock(date string) func() time.Time {
	return func() time.Time {
		t, _ := time.Parse("2006-01-02", date)
		return t
	}
}

func TestGetCurrentWeekDate(t *testing.T) {
	testCases := []struct {
		name     string
		date     string
		nowMock  func() time.Time
		expected string
	}{
		{
			name:     "start of the week",
			date:     "2024-09-23",
			nowMock:  createNowMock("2024-09-23"),
			expected: "2024-09-23",
		},
		{
			name:     "middle of the week",
			date:     "2024-09-25",
			nowMock:  createNowMock("2024-09-25"),
			expected: "2024-09-23",
		},
		{
			name:     "end of the week",
			date:     "2024-09-29",
			nowMock:  createNowMock("2024-09-29"),
			expected: "2024-09-23",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := GetCurrentWeekDate(tc.nowMock)
			Equal(t, got, tc.expected)
		})
	}
}

func TestNewErrorStackTraceString(t *testing.T) {
	got := NewErrorStackTraceString("test exception message")
	expected := "test exception message \n from"
	Contains(t, got, expected)
}

func TestGetStartOfWeek(t *testing.T) {
	testCases := []struct {
		name          string
		date          string
		startOfWeek   string
		expectedError string
		exact         bool
	}{
		{
			name:          "start of the week",
			date:          "2024-11-04",
			startOfWeek:   "2024-11-04",
			expectedError: "",
			exact:         true,
		},
		{
			name:          "middle of the week",
			date:          "2024-11-06",
			startOfWeek:   "2024-11-04",
			expectedError: "",
			exact:         true,
		},
		{
			name:          "end of the week",
			date:          "2024-11-10",
			startOfWeek:   "2024-11-04",
			expectedError: "",
			exact:         true,
		},
		{
			name:          "invalid date format",
			date:          "30-11-2024",
			startOfWeek:   "",
			expectedError: "failed to parse date \n from parsing time \"30-11-2024\" as \"2006-01-02\": cannot parse \"30-11-2024\" as \"2006\"",
			exact:         false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := GetStartOfWeek(tc.date)
			if tc.exact {
				if tc.expectedError != "" {
					Equal(t, err.Error(), tc.expectedError)
				} else {
					Equal(t, got, tc.startOfWeek)
				}
			} else {
				Contains(t, err.Error(), tc.expectedError)
			}
		})
	}
}

func TestGetDayOffset(t *testing.T) {
	testCases := []struct {
		name          string
		day           string
		offset        int
		expectedError string
		exact         bool
	}{
		{
			name:          "valid monday",
			day:           "mon",
			offset:        0,
			expectedError: "",
			exact:         true,
		},
		{
			name:          "valid tuesday",
			day:           "tue",
			offset:        1,
			expectedError: "",
			exact:         true,
		},
		{
			name:          "valid sunday",
			day:           "sun",
			offset:        6,
			expectedError: "",
			exact:         true,
		},
		{
			name:          "invalid day",
			day:           "invalid",
			offset:        0,
			expectedError: "invalid day provided. Valid values: mon, tue, wed, thu, fri, sat, sun \n from keyerror: 'invalid' not in days map",
			exact:         false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := getDayOffset(tc.day)
			if tc.exact {
				if tc.expectedError != "" {
					Equal(t, err.Error(), tc.expectedError)
				} else {
					Equal(t, got, tc.offset)
				}
			} else {
				Contains(t, err.Error(), tc.expectedError)
			}
		})
	}
}

func TestGetDateFromDay(t *testing.T) {
	testCases := []struct {
		name          string
		day           string
		date          string
		nowMock       func() time.Time
		expectedError string
		exact         bool
	}{
		{
			name:          "valid monday",
			day:           "mon",
			date:          "2024-09-23",
			nowMock:       createNowMock("2024-09-23"),
			expectedError: "",
			exact:         true,
		},
		{
			name:          "valid wednesday",
			day:           "wed",
			date:          "2024-09-25",
			nowMock:       createNowMock("2024-09-25"),
			expectedError: "",
			exact:         true,
		},
		{
			name:          "valid sunday",
			day:           "sun",
			date:          "2024-09-29",
			nowMock:       createNowMock("2024-09-29"),
			expectedError: "",
			exact:         true,
		},
		{
			name:          "invalid day",
			day:           "invalid",
			date:          "",
			nowMock:       createNowMock("2024-09-23"),
			expectedError: "invalid day provided. Valid values: mon, tue, wed, thu, fri, sat, sun \n from keyerror: 'invalid' not in days map",
			exact:         false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := GetDateFromDay(tc.day, tc.nowMock)
			if tc.exact {
				if tc.expectedError != "" {
					Equal(t, err.Error(), tc.expectedError)
				} else {
					Equal(t, got, tc.date)
				}
			} else {
				Contains(t, err.Error(), tc.expectedError)
			}
		})
	}
}
