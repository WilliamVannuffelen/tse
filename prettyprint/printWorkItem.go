package prettyprint

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/williamvannuffelen/tse/workitem"
)

func PrintWorkItem(workItem *workitem.KiaraWorkItem) {
	keyColor := color.New(color.FgCyan).SprintFunc()
	valueColor := color.New(color.FgYellow).SprintFunc()
	fmt.Println("Added timesheet entry:")
	fmt.Println(keyColor("  \"Day\":"), valueColor(workItem.Day))
	fmt.Println(keyColor("  \"Date\":"), valueColor(workItem.Date))
	fmt.Println(keyColor("  \"Description\":"), valueColor(workItem.Description))
	fmt.Println(keyColor("  \"JiraRef\":"), valueColor(workItem.JiraRef))
	fmt.Println(keyColor("  \"TimeSpent\":"), valueColor(workItem.TimeSpent))
	fmt.Println(keyColor("  \"Project\":"), valueColor(workItem.Project))
	fmt.Println(keyColor("  \"AppRef\":"), valueColor(workItem.AppRef))
}
