package prettyprint

import (
	"fmt"
	"github.com/williamvannuffelen/tse/workitem"
)

func PrintDayInSelectedFormat(
	values map[string]interface{},
	timeSpentPerDay []workitem.TimeSpentPerDay,
	startOfWeek string,
	workItems []workitem.KiaraWorkItem,
	aggregatedWorkItems []workitem.AggregatedWorkItem,
) {
	if values["output"] == "pp" {
		fmt.Println("Showing only the selected day:", values["day"], values["date"])
		PrintTimeSpentPerDayTable(timeSpentPerDay, values["date"].(string))
		PrintSingleDayWorkItemTable(workItems, values["date"].(string), !(values["hide-project"].(bool)), !(values["hide-appref"].(bool)), !(values["hide-jiraref"].(bool)))
	} else {
		PrintSingleDayWorkItemsAsJson(workItems, values["date"].(string), timeSpentPerDay)
	}
}

func PrintWeekInSelectedFormat(
	values map[string]interface{},
	timeSpentPerDay []workitem.TimeSpentPerDay,
	startOfWeek string,
	workItems []workitem.KiaraWorkItem,
	aggregatedWorkItems []workitem.AggregatedWorkItem,
) {
	if values["output"] == "pp" {
		fmt.Println("Showing entire week starting on ", startOfWeek)
		PrintTimeSpentPerDayTable(timeSpentPerDay, "")
		PrintAggregatedWorkItemTable(aggregatedWorkItems, !(values["hide-project"].(bool)), !(values["hide-appref"].(bool)), !(values["hide-jiraref"].(bool)))
	} else {
		PrintWorkItemsAsJson(workItems, timeSpentPerDay, startOfWeek)
	}
}
