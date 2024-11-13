package prettyprint

import (
	"bytes"
	"testing"

	"github.com/williamvannuffelen/tse/workitem"
)

// minimal test, just validating all object properties are present
// not worth testing formatting of pretty-printed table
func TestPrintAggregatedWorkItemTable(t *testing.T) {
	testCases := []struct {
		name            string
		workItems       []workitem.AggregatedWorkItem
		expectedStrings []string
	}{
		{
			name: "Single work item",
			workItems: []workitem.AggregatedWorkItem{
				{
					Description: "Implement feature X",
					JiraRef:     "JIRA-123",
					TotalTime:   4,
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
			workItems: []workitem.AggregatedWorkItem{
				{
					Description: "Implement feature X",
					JiraRef:     "JIRA-123",
					TotalTime:   4,
					Project:     "Project A",
					AppRef:      "App-456",
				},
				{
					Description: "Fix bug Y",
					JiraRef:     "JIRA-124",
					TotalTime:   2,
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
				"JIRA-124",
				"2",
				"Project B",
				"App-457",
			},
		},
		{
			name:            "No work items",
			workItems:       []workitem.AggregatedWorkItem{},
			expectedStrings: []string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var buffer bytes.Buffer
			PrintAggregatedWorkItemTable(&buffer, tc.workItems, true, true, true)

			got := buffer.String()
			for _, expected := range tc.expectedStrings {
				Contains(t, got, expected)
			}
		})
	}
}
