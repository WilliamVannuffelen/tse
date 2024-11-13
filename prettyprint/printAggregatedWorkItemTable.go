package prettyprint

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/williamvannuffelen/tse/workitem"
	"io"
	"sort"
	"strings"
	"text/tabwriter"
)

func PrintAggregatedWorkItemTable(w io.Writer, workItems []workitem.AggregatedWorkItem, showProject bool, showAppRef bool, showJiraRef bool) {
	sort.Slice(workItems, func(i, j int) bool {
		return workItems[i].Description < workItems[j].Description
	})
	maxDescriptionLen := len("Description")
	maxJiraRefLen := len("JiraRef")
	maxProjectLen := len("Project")
	maxAppRefLen := len("AppRef")
	maxTotalTimeLen := len("TotalTimeSpent")

	hasDescription := false
	hasJiraRef := false
	hasProject := false
	hasAppRef := false
	hasTotalTime := false

	for _, entry := range workItems {
		if len(entry.Description) > 0 {
			hasDescription = true
			if len(entry.Description) > maxDescriptionLen {
				maxDescriptionLen = len(entry.Description)
			}
		}
		if len(entry.JiraRef) > 0 {
			hasJiraRef = true
			if len(entry.JiraRef) > maxJiraRefLen {
				maxJiraRefLen = len(entry.JiraRef)
			}
		}
		if len(entry.Project) > 0 {
			hasProject = true
			if len(entry.Project) > maxProjectLen {
				maxProjectLen = len(entry.Project)
			}
		}
		if len(entry.AppRef) > 0 {
			hasAppRef = true
			if len(entry.AppRef) > maxAppRefLen {
				maxAppRefLen = len(entry.AppRef)
			}
		}
		if entry.TotalTime > 0 {
			hasTotalTime = true
		}
	}

	tableColor := color.New(color.FgBlue, color.Bold).SprintFunc()
	descriptionHeaderColor := color.New(color.FgYellow, color.Bold).SprintFunc()
	jiraRefHeaderColor := color.New(color.FgMagenta, color.Bold).SprintFunc()
	projectHeaderColor := color.New(color.FgGreen, color.Bold).SprintFunc()
	appRefHeaderColor := color.New(color.FgCyan, color.Bold).SprintFunc()
	totalTimeHeaderColor := color.New(color.FgRed, color.Bold).SprintFunc()
	greenColor := color.New(color.FgGreen, color.Bold).SprintFunc()
	redColor := color.New(color.FgRed, color.Bold).SprintFunc()
	blueBoldColor := color.New(color.FgBlue, color.Bold).SprintFunc()

	writer := tabwriter.NewWriter(w, 0, 0, 1, ' ', tabwriter.Debug)

	header := tableColor("|")
	separator := tableColor("+")

	if hasDescription {
		header += fmt.Sprintf(" %s %s", descriptionHeaderColor(fmt.Sprintf("%-*s", maxDescriptionLen, "Description")), tableColor("|"))
		separator += tableColor(fmt.Sprintf("-%s-+", strings.Repeat("-", maxDescriptionLen)))
	}
	if hasJiraRef && showJiraRef {
		header += fmt.Sprintf(" %s %s", jiraRefHeaderColor(fmt.Sprintf("%-*s", maxJiraRefLen, "JiraRef")), tableColor("|"))
		separator += tableColor(fmt.Sprintf("-%s-+", strings.Repeat("-", maxJiraRefLen)))
	}
	if hasProject && showProject {
		header += fmt.Sprintf(" %s %s", projectHeaderColor(fmt.Sprintf("%-*s", maxProjectLen, "Project")), tableColor("|"))
		separator += tableColor(fmt.Sprintf("-%s-+", strings.Repeat("-", maxProjectLen)))
	}
	if hasAppRef && showAppRef {
		header += fmt.Sprintf(" %s %s", appRefHeaderColor(fmt.Sprintf("%-*s", maxAppRefLen, "AppRef")), tableColor("|"))
		separator += tableColor(fmt.Sprintf("-%s-+", strings.Repeat("-", maxAppRefLen)))
	}
	if hasTotalTime {
		header += fmt.Sprintf(" %s %s", totalTimeHeaderColor(fmt.Sprintf("%-*s", maxTotalTimeLen, "TotalTimeSpent")), tableColor("|"))
		separator += tableColor(fmt.Sprintf("-%s-+", strings.Repeat("-", maxTotalTimeLen)))
	}

	fmt.Fprintf(writer, "%s\n", separator)
	fmt.Fprintf(writer, "%s\n", header)
	fmt.Fprintf(writer, "%s\n", separator)

	totalTimeSpent := 0.0
	for _, entry := range workItems {
		timeSpentColor := redColor
		row := tableColor("|")
		if hasDescription {
			row += fmt.Sprintf(" %s %s", descriptionHeaderColor(fmt.Sprintf("%-*s", maxDescriptionLen, entry.Description)), tableColor("|"))
		}
		if hasJiraRef && showJiraRef {
			row += fmt.Sprintf(" %s %s", jiraRefHeaderColor(fmt.Sprintf("%-*s", maxJiraRefLen, entry.JiraRef)), tableColor("|"))
		}
		if hasProject && showProject {
			row += fmt.Sprintf(" %s %s", projectHeaderColor(fmt.Sprintf("%-*s", maxProjectLen, entry.Project)), tableColor("|"))
		}
		if hasAppRef && showAppRef {
			row += fmt.Sprintf(" %s %s", appRefHeaderColor(fmt.Sprintf("%-*s", maxAppRefLen, entry.AppRef)), tableColor("|"))
		}
		if hasTotalTime {
			row += fmt.Sprintf(" %s %s", timeSpentColor(fmt.Sprintf("%-*s", maxTotalTimeLen, fmt.Sprintf("%.2f", entry.TotalTime))), tableColor("|"))
		}
		fmt.Fprintf(writer, "%s\n", row)
		totalTimeSpent += entry.TotalTime
	}
	fmt.Fprintf(writer, "%s\n", separator)
	writer.Flush()

	timeSpentColor := redColor
	if totalTimeSpent >= 40 {
		timeSpentColor = greenColor
	}
	fmt.Println(blueBoldColor("Week Total Time Spent:"), timeSpentColor(fmt.Sprintf("%.2f hours", totalTimeSpent)))
}
