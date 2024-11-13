package prettyprint

import (
	"bytes"
	"testing"

	"github.com/fatih/color"
	"github.com/williamvannuffelen/tse/workitem"
)

func TestPrintWorkItem(t *testing.T) {
	var buffer bytes.Buffer
	workItem := &workitem.KiaraWorkItem{
		Day:         "Monday",
		Date:        "2023-10-02",
		Description: "Implement feature X",
		JiraRef:     "JIRA-123",
		TimeSpent:   "4h",
		Project:     "Project A",
		AppRef:      "App-456",
	}

	PrintWorkItem(&buffer, workItem)

	keyColor := color.New(color.FgCyan).SprintFunc()
	valueColor := color.New(color.FgYellow).SprintFunc()
	expectedOutput := "Added timesheet entry:\n" +
		keyColor("  \"Day\":") + " " + valueColor("Monday") + "\n" +
		keyColor("  \"Date\":") + " " + valueColor("2023-10-02") + "\n" +
		keyColor("  \"Description\":") + " " + valueColor("Implement feature X") + "\n" +
		keyColor("  \"JiraRef\":") + " " + valueColor("JIRA-123") + "\n" +
		keyColor("  \"TimeSpent\":") + " " + valueColor("4h") + "\n" +
		keyColor("  \"Project\":") + " " + valueColor("Project A") + "\n" +
		keyColor("  \"AppRef\":") + " " + valueColor("App-456") + "\n"

	if buffer.String() != expectedOutput {
		t.Errorf("expected %q but got %q", expectedOutput, buffer.String())
	}
}
