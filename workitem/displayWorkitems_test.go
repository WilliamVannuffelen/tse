package workitem

import (
	"sort"
	"strings"
	"testing"
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
		%v,
		should be present in:
		%v
		)`, expected, got)
	}
}

func TestTimeSpentToFLoat(t *testing.T) {
	testCases := []struct {
		name          string
		timeSpent     string
		expected      float64
		expectedError string
		exact         bool
	}{
		{
			name:          "valid time spent",
			timeSpent:     "2.5",
			expected:      2.5,
			expectedError: "",
			exact:         true,
		},
		{
			name:          "invalid time spent",
			timeSpent:     "invalid",
			expected:      0,
			expectedError: "failed to parse time spent",
			exact:         false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := timeSpentToFloat(tc.timeSpent)
			if tc.expectedError != "" {
				if tc.exact {
					Equal(t, err.Error(), tc.expectedError)
				} else {
					Contains(t, err.Error(), tc.expectedError)
				}
			} else {
				Equal(t, got, tc.expected)
			}
		})
	}
}

func equalAggregatedWorkItems(got, expected []AggregatedWorkItem) bool {
	if len(got) != len(expected) {
		return false
	}
	for i := range got {
		if got[i].Description != expected[i].Description ||
			got[i].AppRef != expected[i].AppRef ||
			got[i].JiraRef != expected[i].JiraRef ||
			got[i].Project != expected[i].Project ||
			got[i].TotalTime != expected[i].TotalTime {
			return false
		}
	}
	return true
}

func TestAggregateWorkItems(t *testing.T) {
	testCases := []struct {
		name            string
		items           []KiaraWorkItem
		aggregatedItems []AggregatedWorkItem
		expectedError   string
		exact           bool
	}{
		{
			name: "Valid work items",
			items: []KiaraWorkItem{
				{Description: "Task 1", AppRef: "App 1", JiraRef: "OPS-305", Project: "Project 1", TimeSpent: "2"},
				{Description: "Task 2", AppRef: "App 2", JiraRef: "OPS-306", Project: "Project 2", TimeSpent: "1"},
				{Description: "Task 1", AppRef: "App 1", JiraRef: "OPS-305", Project: "Project 1", TimeSpent: "5"},
			},
			aggregatedItems: []AggregatedWorkItem{
				{Description: "Task 1", AppRef: "App 1", JiraRef: "OPS-305", Project: "Project 1", TotalTime: 7},
				{Description: "Task 2", AppRef: "App 2", JiraRef: "OPS-306", Project: "Project 2", TotalTime: 1},
			},
			expectedError: "",
			exact:         true,
		},
		{
			name: "Invalid time format",
			items: []KiaraWorkItem{
				{Description: "Task 1", AppRef: "App1", JiraRef: "Jira1", Project: "Project1", TimeSpent: "invalid"},
			},
			aggregatedItems: nil,
			expectedError:   "failed to aggregate work items",
			exact:           false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := AggregateWorkItems(tc.items)
			if tc.expectedError != "" {
				if tc.exact {
					Equal(t, err.Error(), tc.expectedError)
				} else {
					Contains(t, err.Error(), tc.expectedError)
				}
			} else {
				sort.Slice(got, func(i, j int) bool {
					return got[i].Description < got[j].Description
				})
				sort.Slice(tc.aggregatedItems, func(i, j int) bool {
					return tc.aggregatedItems[i].Description < tc.aggregatedItems[j].Description
				})
				if !equalAggregatedWorkItems(got, tc.aggregatedItems) {
					t.Errorf("AggregateWorkItems() = %v, expected %v", got, tc.aggregatedItems)
				}
			}
		})
	}
}

func TestCalculateTotalTimeSpent(t *testing.T) {
	testCases := []struct {
		name          string
		entries       []KiaraWorkItem
		expected      float64
		expectedError string
		exact         bool
	}{
		{
			name: "Valid work items",
			entries: []KiaraWorkItem{
				{Description: "Task 1", AppRef: "App 1", JiraRef: "OPS-305", Project: "Project 1", TimeSpent: "2"},
				{Description: "Task 2", AppRef: "App 2", JiraRef: "OPS-306", Project: "Project 2", TimeSpent: "1"},
				{Description: "Task 1", AppRef: "App 1", JiraRef: "OPS-305", Project: "Project 1", TimeSpent: "5"},
			},
			expected:      8,
			expectedError: "",
			exact:         true,
		},
		{
			name: "Invalid time format",
			entries: []KiaraWorkItem{
				{Description: "Task 1", AppRef: "App1", JiraRef: "Jira1", Project: "Project1", TimeSpent: "invalid"},
			},
			expected:      0,
			expectedError: "failed to calculate total time spent",
			exact:         false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := CalculateTotalTimeSpent(tc.entries)
			if tc.expectedError != "" {
				if tc.exact {
					Equal(t, err.Error(), tc.expectedError)
				} else {
					Contains(t, err.Error(), tc.expectedError)
				}
			} else {
				Equal(t, got, tc.expected)
			}
		})
	}
}

func equalTimeSpentPerDay(got, expected []TimeSpentPerDay) bool {
	if len(got) != len(expected) {
		return false
	}
	for i := range got {
		if got[i].Day != expected[i].Day ||
			got[i].Date != expected[i].Date ||
			got[i].TimeSpent != expected[i].TimeSpent {
			return false
		}
	}
	return true
}

func TestCalculateTimeSpentPerDay(t *testing.T) {
	testCases := []struct {
		name            string
		entries         []KiaraWorkItem
		timeSpentPerDay []TimeSpentPerDay
		expectedError   string
		exact           bool
	}{
		{
			name: "Valid work items",
			entries: []KiaraWorkItem{
				{Description: "Task 1", AppRef: "App 1", JiraRef: "OPS-305", Project: "Project 1", TimeSpent: "2", Date: "2024-09-30", Day: "Monday"},
				{Description: "Task 2", AppRef: "App 2", JiraRef: "OPS-306", Project: "Project 2", TimeSpent: "1", Date: "2024-10-01", Day: "Tuesday"},
				{Description: "Task 1", AppRef: "App 1", JiraRef: "OPS-305", Project: "Project 1", TimeSpent: "5", Date: "2024-09-30", Day: "Monday"},
			},
			timeSpentPerDay: []TimeSpentPerDay{
				{Day: "Monday", Date: "2024-09-30", TimeSpent: 7},
				{Day: "Tuesday", Date: "2024-10-01", TimeSpent: 1},
			},
			expectedError: "",
			exact:         true,
		},
		{
			name: "Invalid time format",
			entries: []KiaraWorkItem{
				{Description: "Task 1", AppRef: "App1", JiraRef: "Jira1", Project: "Project1", TimeSpent: "invalid", Date: "2021-01-01"},
			},
			timeSpentPerDay: nil,
			expectedError:   "failed to calculate time spent per day",
			exact:           false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := CalculateTimeSpentPerDay(tc.entries)
			if tc.expectedError != "" {
				if tc.exact {
					Equal(t, err.Error(), tc.expectedError)
				} else {
					Contains(t, err.Error(), tc.expectedError)
				}
			} else {
				sort.Slice(got, func(i, j int) bool {
					return got[i].Date < got[j].Date
				})
				sort.Slice(tc.timeSpentPerDay, func(i, j int) bool {
					return tc.timeSpentPerDay[i].Date < tc.timeSpentPerDay[j].Date
				})
				if !equalTimeSpentPerDay(got, tc.timeSpentPerDay) {
					t.Errorf("CalculateTimeSpentPerDay() = %v, expected %v", got, tc.timeSpentPerDay)
				}
			}
		})
	}
}
