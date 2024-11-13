package prettyprint

import (
	"bytes"
	"testing"

	"github.com/williamvannuffelen/tse/workitem"
)

func TestPrintTimeSpentPerDayTable(t *testing.T) {
	testCases := []struct {
		name            string
		days            []workitem.TimeSpentPerDay
		date            string
		expectedStrings []string
	}{
		{
			name: "Single day",
			days: []workitem.TimeSpentPerDay{
				{
					Day:       "Monday",
					Date:      "2023-10-02",
					TimeSpent: 4,
				},
			},
			date: "2023-10-02",
			expectedStrings: []string{
				"Day",
				"Date",
				"TimeSpent",
				"Monday",
				"2023-10-02",
				"4",
			},
		},
		{
			name: "Multiple days no filter",
			days: []workitem.TimeSpentPerDay{
				{
					Day:       "Monday",
					Date:      "2023-10-02",
					TimeSpent: 4,
				},
				{
					Day:       "Tuesday",
					Date:      "2023-10-03",
					TimeSpent: 2,
				},
			},
			expectedStrings: []string{
				"Day",
				"Date",
				"TimeSpent",
				"Monday",
				"2023-10-02",
				"4",
				"Tuesday",
				"2023-10-03",
				"2",
			},
		},
		{
			name: "Multiple days with filter",
			days: []workitem.TimeSpentPerDay{
				{
					Day:       "Monday",
					Date:      "2023-10-02",
					TimeSpent: 4,
				},
				{
					Day:       "Tuesday",
					Date:      "2023-10-03",
					TimeSpent: 2,
				},
			},
			date: "2023-10-03",
			expectedStrings: []string{
				"Day",
				"Date",
				"TimeSpent",
				"Tuesday",
				"2023-10-03",
				"2",
			},
		},
		{
			name: "No days",
			days: []workitem.TimeSpentPerDay{},
			expectedStrings: []string{
				"Day",
				"Date",
				"TimeSpent",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var buffer bytes.Buffer
			PrintTimeSpentPerDayTable(&buffer, tc.days, tc.date)

			got := buffer.String()
			for _, expected := range tc.expectedStrings {
				Contains(t, got, expected)
			}
		})
	}
}
