package prettyprint

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/williamvannuffelen/tse/workitem"
	"os"
	"sort"
	"strings"
	"text/tabwriter"
)

func PrintSingleDayWorkItemTable(workItems []workitem.KiaraWorkItem, date string, showProject bool, showAppRef bool, showJiraRef bool) {
	filteredWorkItems := []workitem.KiaraWorkItem{}
	for _, item := range workItems {
		if item.Date == date {
			filteredWorkItems = append(filteredWorkItems, item)
		}
	}

	sort.Slice(filteredWorkItems, func(i, j int) bool {
		return filteredWorkItems[i].Description < filteredWorkItems[j].Description
	})

	maxDescriptionLen := len("Description")
	maxJiraRefLen := len("JiraRef")
	maxProjectLen := len("Project")
	maxAppRefLen := len("AppRef")
	maxTimeSpentLen := len("TimeSpent")

	hasDescription := false
	hasJiraRef := false
	hasProject := false
	hasAppRef := false
	hasTimeSpent := false

	for _, entry := range filteredWorkItems {
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
		if len(entry.TimeSpent) > 0 {
			hasTimeSpent = true
		}
	}

	tableColor := color.New(color.FgBlue, color.Bold).SprintFunc()
	descriptionHeaderColor := color.New(color.FgYellow, color.Bold).SprintFunc()
	jiraRefHeaderColor := color.New(color.FgMagenta, color.Bold).SprintFunc()
	projectHeaderColor := color.New(color.FgGreen, color.Bold).SprintFunc()
	appRefHeaderColor := color.New(color.FgCyan, color.Bold).SprintFunc()
	timeSpentHeaderColor := color.New(color.FgRed, color.Bold).SprintFunc()

	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)
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
	if hasTimeSpent {
		header += fmt.Sprintf(" %s %s", timeSpentHeaderColor(fmt.Sprintf("%-*s", maxTimeSpentLen, "TimeSpent")), tableColor("|"))
		separator += tableColor(fmt.Sprintf("-%s-+", strings.Repeat("-", maxTimeSpentLen)))
	}

	fmt.Fprintf(writer, "%s\n", separator)
	fmt.Fprintf(writer, "%s\n", header)
	fmt.Fprintf(writer, "%s\n", separator)

	for _, entry := range filteredWorkItems {
		//timeSpentColor := redColor
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
		if hasTimeSpent {
			row += fmt.Sprintf(" %s %s", timeSpentHeaderColor(fmt.Sprintf("%-*s", maxTimeSpentLen, entry.TimeSpent)), tableColor("|"))
		}
		fmt.Fprintf(writer, "%s\n", row)
	}
	fmt.Fprintf(writer, "%s\n", separator)
	writer.Flush()
}
