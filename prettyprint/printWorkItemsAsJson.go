package prettyprint

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"

	"github.com/williamvannuffelen/tse/workitem"
)

type WorkItemOutput struct {
	Week               string                   `json:"week"`
	TimeSpent          float64                  `json:"timespent"`
	TimeSpentMonday    float64                  `json:"timespentMonday"`
	TimeSpentTuesday   float64                  `json:"timespentTuesday"`
	TimeSpentWednesday float64                  `json:"timespentWednesday"`
	TimeSpentThursday  float64                  `json:"timespentThursday"`
	TimeSpentFriday    float64                  `json:"timespentFriday"`
	TimeSpentSaturday  float64                  `json:"timespentSaturday"`
	TimeSpentSunday    float64                  `json:"timespentSunday"`
	Items              []workitem.KiaraWorkItem `json:"items"`
}

func PrintWorkItemsAsJson(workItems []workitem.KiaraWorkItem, timeSpentPerDay []workitem.TimeSpentPerDay, firstDateOfWeek string) {
	sort.Slice(workItems, func(i, j int) bool {
		return workItems[i].Description < workItems[j].Description
	})

	totalTimeSpent := 0.0
	timeSpentMonday := 0.0
	timeSpentTuesday := 0.0
	timeSpentWednesday := 0.0
	timeSpentThursday := 0.0
	timeSpentFriday := 0.0
	timeSpentSaturday := 0.0
	timeSpentSunday := 0.0

	for _, item := range workItems {
		timeSpent, err := strconv.ParseFloat(item.TimeSpent, 64)
		if err == nil {
			totalTimeSpent += timeSpent
			switch item.Day {
			case "Mon":
				timeSpentMonday += timeSpent
			case "Tue":
				timeSpentTuesday += timeSpent
			case "Wed":
				timeSpentWednesday += timeSpent
			case "Thu":
				timeSpentThursday += timeSpent
			case "Fri":
				timeSpentFriday += timeSpent
			case "Sat":
				timeSpentSaturday += timeSpent
			case "Sun":
				timeSpentSunday += timeSpent
			}
		}
	}

	output := WorkItemOutput{
		Week:               firstDateOfWeek,
		TimeSpent:          totalTimeSpent,
		TimeSpentMonday:    timeSpentMonday,
		TimeSpentTuesday:   timeSpentTuesday,
		TimeSpentWednesday: timeSpentWednesday,
		TimeSpentThursday:  timeSpentThursday,
		TimeSpentFriday:    timeSpentFriday,
		TimeSpentSaturday:  timeSpentSaturday,
		TimeSpentSunday:    timeSpentSunday,
		Items:              workItems,
	}

	jsonOutput, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling output: ", err)
		return
	}
	fmt.Println(string(jsonOutput))
}
