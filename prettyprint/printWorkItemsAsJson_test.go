package prettyprint

import (
	"bytes"
	"encoding/json"
	"sort"
	"strings"
	"testing"

	"github.com/williamvannuffelen/tse/workitem"
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

func sortItems(items []workitem.KiaraWorkItem) {
	sort.Slice(items, func(i, j int) bool {
		if items[i].Description == items[j].Description {
			return items[i].Day < items[j].Day
		}
		return items[i].Description < items[j].Description
	})
}

func validateWorkItemOutput(t *testing.T, actual, expected WorkItemOutput) {
	Equal(t, actual.Week, expected.Week)
	Equal(t, actual.TimeSpent, expected.TimeSpent)
	Equal(t, actual.TimeSpentMonday, expected.TimeSpentMonday)
	Equal(t, actual.TimeSpentTuesday, expected.TimeSpentTuesday)
	Equal(t, actual.TimeSpentWednesday, expected.TimeSpentWednesday)
	Equal(t, actual.TimeSpentThursday, expected.TimeSpentThursday)
	Equal(t, actual.TimeSpentFriday, expected.TimeSpentFriday)
	Equal(t, actual.TimeSpentSaturday, expected.TimeSpentSaturday)
	Equal(t, actual.TimeSpentSunday, expected.TimeSpentSunday)
	Equal(t, len(actual.Items), len(expected.Items))

	sortItems(actual.Items)
	sortItems(expected.Items)

	for i := range actual.Items {
		Equal(t, actual.Items[i], expected.Items[i])
	}
}

func TestPrintWorkItemsAsJson(t *testing.T) {
	var buffer bytes.Buffer
	testCases := []struct {
		name            string
		workItems       []workitem.KiaraWorkItem
		timeSpentPerDay []workitem.TimeSpentPerDay
		firstDateOfWeek string
		expectedOutput  WorkItemOutput
	}{
		{
			name: "Basic test case",
			workItems: []workitem.KiaraWorkItem{
				{Description: "Task 1", TimeSpent: "2.5", Day: "Mon"},
				{Description: "Task 2", TimeSpent: "3.0", Day: "Tue"},
				{Description: "Task 3", TimeSpent: "2.0", Day: "Mon"},
				{Description: "Task 3", TimeSpent: "2.0", Day: "Mon"},
				{Description: "Task 2", TimeSpent: "1.0", Day: "Tue"},
				{Description: "Task 3", TimeSpent: "1.0", Day: "Wed"},
				{Description: "Task 4", TimeSpent: "3.0", Day: "Wed"},
				{Description: "Task 3", TimeSpent: "2.0", Day: "Thu"},
				{Description: "Task 3", TimeSpent: "2.0", Day: "Fri"},
				{Description: "Task 3", TimeSpent: "2.0", Day: "Sat"},
				{Description: "Task 3", TimeSpent: "2.0", Day: "Sun"},
			},

			firstDateOfWeek: "2023-10-01",
			expectedOutput: WorkItemOutput{
				Week:               "2023-10-01",
				TimeSpent:          22.5,
				TimeSpentMonday:    6.5,
				TimeSpentTuesday:   4.0,
				TimeSpentWednesday: 4.0,
				TimeSpentThursday:  2.0,
				TimeSpentFriday:    2.0,
				TimeSpentSaturday:  2.0,
				TimeSpentSunday:    2.0,
				Items: []workitem.KiaraWorkItem{
					{Description: "Task 1", TimeSpent: "2.5", Day: "Mon"},
					{Description: "Task 2", TimeSpent: "3.0", Day: "Tue"},
					{Description: "Task 3", TimeSpent: "2.0", Day: "Mon"},
					{Description: "Task 3", TimeSpent: "2.0", Day: "Mon"},
					{Description: "Task 2", TimeSpent: "1.0", Day: "Tue"},
					{Description: "Task 3", TimeSpent: "1.0", Day: "Wed"},
					{Description: "Task 4", TimeSpent: "3.0", Day: "Wed"},
					{Description: "Task 3", TimeSpent: "2.0", Day: "Thu"},
					{Description: "Task 3", TimeSpent: "2.0", Day: "Fri"},
					{Description: "Task 3", TimeSpent: "2.0", Day: "Sat"},
					{Description: "Task 3", TimeSpent: "2.0", Day: "Sun"},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			buffer.Reset()
			PrintWorkItemsAsJson(&buffer, tc.workItems, tc.timeSpentPerDay, tc.firstDateOfWeek)
			var output WorkItemOutput
			err := json.Unmarshal(buffer.Bytes(), &output)
			if err != nil {
				t.Fatalf("failed to unmarshal json: %v", err)
			}
			validateWorkItemOutput(t, output, tc.expectedOutput)
		})
	}
}
