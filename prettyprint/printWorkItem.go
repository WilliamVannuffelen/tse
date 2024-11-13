package prettyprint

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/williamvannuffelen/tse/workitem"
	"io"
)

func PrintWorkItem(w io.Writer, workItem *workitem.KiaraWorkItem) {
	keyColor := color.New(color.FgCyan).SprintFunc()
	valueColor := color.New(color.FgYellow).SprintFunc()
	fmt.Fprintln(w, "Added timesheet entry:")
	fmt.Fprintln(w, keyColor("  \"Day\":"), valueColor(workItem.Day))
	fmt.Fprintln(w, keyColor("  \"Date\":"), valueColor(workItem.Date))
	fmt.Fprintln(w, keyColor("  \"Description\":"), valueColor(workItem.Description))
	fmt.Fprintln(w, keyColor("  \"JiraRef\":"), valueColor(workItem.JiraRef))
	fmt.Fprintln(w, keyColor("  \"TimeSpent\":"), valueColor(workItem.TimeSpent))
	fmt.Fprintln(w, keyColor("  \"Project\":"), valueColor(workItem.Project))
	fmt.Fprintln(w, keyColor("  \"AppRef\":"), valueColor(workItem.AppRef))
}
