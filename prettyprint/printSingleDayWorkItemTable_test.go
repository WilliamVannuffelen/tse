package prettyprint

import (
	"bytes"
	"testing"

	"github.com/williamvannuffelen/tse/workitem"
)

// minimal test, just validating all object properties are present
// not worth testing formatting of pretty-printed table
func TestPrintSingleDayWorkItemTable(t *testing.T) {
	testCases := []struct {
		name            string
		date            string
		workItems       []workitem.KiaraWorkItem
		expectedStrings []string
	}{
		{
			name: "Single work item",
			date: "2021-01-01",
			workItems: []workitem.KiaraWorkItem{
				{
					Day:         "Monday",
					Date:        "2021-01-01",
					Description: "Implement feature X",
					JiraRef:     "JIRA-123",
					TimeSpent:   "4",
					Project:     "Project A",
					AppRef:      "App-456",
				},
			},
			expectedStrings: []string{
				"Description",
				"JiraRef",
				"Project",
				"AppRef",
				"Implement feature X",
				"JIRA-123",
				"4",
				"Project A",
				"App-456",
			},
		},
		{
			name: "Multiple work items",
			date: "2021-01-01",
			workItems: []workitem.KiaraWorkItem{
				{
					Day:         "Monday",
					Date:        "2021-01-01",
					Description: "Implement feature X",
					JiraRef:     "JIRA-123",
					TimeSpent:   "4",
					Project:     "Project A",
					AppRef:      "App-456",
				},
				{
					Day:         "Tuesday",
					Date:        "2021-01-02",
					Description: "Fix bug Y",
					JiraRef:     "JIRA-124",
					TimeSpent:   "2",
					Project:     "Project B",
					AppRef:      "App-457",
				},
			},
			expectedStrings: []string{
				"Description",
				"JiraRef",
				"Project",
				"AppRef",
				"Implement feature X",
				"JIRA-123",
				"4",
				"Project A",
				"App-456",
			},
		},
		{
			name:            "No work items",
			date:            "2021-01-01",
			workItems:       []workitem.KiaraWorkItem{},
			expectedStrings: []string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var buffer bytes.Buffer
			PrintSingleDayWorkItemTable(&buffer, tc.workItems, tc.date, true, true, true)

			got := buffer.String()
			for _, expected := range tc.expectedStrings {
				Contains(t, got, expected)
			}
		})
	}
}
