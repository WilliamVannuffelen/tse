package prettyprint

import (
	"encoding/json"
	"fmt"
	"io"
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

func PrintSingleDayWorkItemsAsJson(w io.Writer, workItems []workitem.KiaraWorkItem, date string) {
	filteredWorkItems := []workitem.KiaraWorkItem{}
	for _, item := range workItems {
		if item.Date == date {
			fmt.Println("item.Date: ", item.Date)
			filteredWorkItems = append(filteredWorkItems, item)
		}
	}

	sort.Slice(filteredWorkItems, func(i, j int) bool {
		return filteredWorkItems[i].Description < filteredWorkItems[j].Description
	})

	totalTimeSpent := 0.0
	for _, item := range filteredWorkItems {
		fmt.Println("Processing item: ", item)
		timeSpent, err := strconv.ParseFloat(item.TimeSpent, 64)
		if err == nil {
			fmt.Println("timeSpent incremented with:  ", timeSpent)
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
	fmt.Println(output)
	jsonOutput, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		fmt.Fprintln(w, "Error marshalling output: ", err)
		return
	}
	fmt.Fprintln(w, string(jsonOutput))
}
