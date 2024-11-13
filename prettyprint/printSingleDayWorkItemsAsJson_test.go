package prettyprint

import (
	"bytes"
	"encoding/json"
	"sort"
	"testing"

	"github.com/williamvannuffelen/tse/workitem"
)

func TestPrintSingleDayWorkItemsAsJson(t *testing.T) {
	testCases := []struct {
		name      string
		workItems []workitem.KiaraWorkItem
		date      string
		expected  Output
	}{
		{
			name: "Single work item",
			workItems: []workitem.KiaraWorkItem{
				{
					Day:         "Monday",
					Date:        "2023-10-02",
					Description: "Implement feature X",
					JiraRef:     "JIRA-123",
					TimeSpent:   "4",
					Project:     "Project A",
					AppRef:      "App-456",
				},
			},
			date: "2023-10-02",
			expected: Output{
				Day:       "Monday",
				Date:      "2023-10-02",
				TimeSpent: 4.0,
				Items: []workitem.KiaraWorkItem{
					{
						Day:         "Monday",
						Date:        "2023-10-02",
						Description: "Implement feature X",
						JiraRef:     "JIRA-123",
						TimeSpent:   "4",
						Project:     "Project A",
						AppRef:      "App-456",
					},
				},
			},
		},
		{
			name: "Multiple work items same day",
			workItems: []workitem.KiaraWorkItem{
				{
					Day:         "Monday",
					Date:        "2023-10-02",
					Description: "Implement feature X",
					JiraRef:     "JIRA-123",
					TimeSpent:   "4",
					Project:     "Project A",
					AppRef:      "App-456",
				},
				{
					Day:         "Monday",
					Date:        "2023-10-02",
					Description: "Fix bug Y",
					JiraRef:     "JIRA-124",
					TimeSpent:   "2",
					Project:     "Project B",
					AppRef:      "App-457",
				},
			},
			date: "2023-10-02",
			expected: Output{
				Day:       "Monday",
				Date:      "2023-10-02",
				TimeSpent: 6.0,
				Items: []workitem.KiaraWorkItem{
					{
						Day:         "Monday",
						Date:        "2023-10-02",
						Description: "Implement feature X",
						JiraRef:     "JIRA-123",
						TimeSpent:   "4",
						Project:     "Project A",
						AppRef:      "App-456",
					},
					{
						Day:         "Monday",
						Date:        "2023-10-02",
						Description: "Fix bug Y",
						JiraRef:     "JIRA-124",
						TimeSpent:   "2",
						Project:     "Project B",
						AppRef:      "App-457",
					},
				},
			},
		},
		{
			name:      "No work items",
			workItems: []workitem.KiaraWorkItem{},
			date:      "2023-10-02",
			expected: Output{
				Day:       "",
				Date:      "2023-10-02",
				TimeSpent: 0.0,
				Items:     []workitem.KiaraWorkItem{},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var buffer bytes.Buffer
			PrintSingleDayWorkItemsAsJson(&buffer, tc.workItems, tc.date)

			var got = Output{}
			err := json.Unmarshal(buffer.Bytes(), &got)
			if err != nil {
				t.Fatalf("failed to unmarshal JSON: %v", err)
			}
			Equal(t, got.Date, tc.expected.Date)
			Equal(t, got.Day, tc.expected.Day)
			Equal(t, got.TimeSpent, tc.expected.TimeSpent)

			if len(got.Items) != len(tc.expected.Items) {
				t.Fatalf("expected %d items, got %d", len(tc.expected.Items), len(got.Items))
			}
			if len(got.Items) > 1 {
				sort.Slice(got.Items, func(i, j int) bool {
					return got.Items[i].Description < got.Items[j].Description
				})
				sort.Slice(tc.expected.Items, func(i, j int) bool {
					return tc.expected.Items[i].Description < tc.expected.Items[j].Description
				})
			}
			for i, item := range got.Items {
				Equal(t, item.Day, tc.expected.Items[i].Day)
				Equal(t, item.Date, tc.expected.Items[i].Date)
				Equal(t, item.Description, tc.expected.Items[i].Description)
				Equal(t, item.JiraRef, tc.expected.Items[i].JiraRef)
				Equal(t, item.TimeSpent, tc.expected.Items[i].TimeSpent)
				Equal(t, item.Project, tc.expected.Items[i].Project)
				Equal(t, item.AppRef, tc.expected.Items[i].AppRef)
			}

		})
	}
}
