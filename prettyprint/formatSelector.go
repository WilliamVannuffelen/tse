package prettyprint

import (
	"fmt"
	"github.com/williamvannuffelen/tse/workitem"
	"os"
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
		PrintTimeSpentPerDayTable(os.Stdout, timeSpentPerDay, values["date"].(string))
		PrintSingleDayWorkItemTable(os.Stdout, workItems, values["date"].(string), !(values["hide-project"].(bool)), !(values["hide-appref"].(bool)), !(values["hide-jiraref"].(bool)))
	} else {
		PrintSingleDayWorkItemsAsJson(os.Stdout, workItems, values["date"].(string))
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
		PrintTimeSpentPerDayTable(os.Stdout, timeSpentPerDay, "")
		PrintAggregatedWorkItemTable(os.Stdout, aggregatedWorkItems, !(values["hide-project"].(bool)), !(values["hide-appref"].(bool)), !(values["hide-jiraref"].(bool)))
	} else {
		PrintWorkItemsAsJson(os.Stdout, workItems, timeSpentPerDay, startOfWeek)
	}
}
