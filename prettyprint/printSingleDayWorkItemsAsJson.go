package prettyprint

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"

	"github.com/williamvannuffelen/tse/workitem"
)

type Output struct {
	Day       string                   `json:"day"`
	Date      string                   `json:"date"`
	TimeSpent float64                  `json:"timespent"`
	Items     []workitem.KiaraWorkItem `json:"items"`
}

func PrintSingleDayWorkItemsAsJson(workItems []workitem.KiaraWorkItem, date string, timeSpentPerDay []workitem.TimeSpentPerDay) {
	filteredWorkItems := []workitem.KiaraWorkItem{}
	for _, item := range workItems {
		if item.Date == date {
			filteredWorkItems = append(filteredWorkItems, item)
		}
	}

	sort.Slice(filteredWorkItems, func(i, j int) bool {
		return filteredWorkItems[i].Description < filteredWorkItems[j].Description
	})

	totalTimeSpent := 0.0
	for _, item := range filteredWorkItems {
		timeSpent, err := strconv.ParseFloat(item.TimeSpent, 64)
		if err == nil {
			totalTimeSpent += timeSpent
		}
	}
	day := ""
	if len(filteredWorkItems) > 0 {
		day = filteredWorkItems[0].Day
	}
	output := Output{
		Day:       day,
		Date:      date,
		TimeSpent: totalTimeSpent,
		Items:     filteredWorkItems,
	}
	jsonOutput, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling output: ", err)
		return
	}
	fmt.Println(string(jsonOutput))
}
