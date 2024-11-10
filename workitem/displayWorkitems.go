package workitem

import (
	"fmt"
	help "github.com/williamvannuffelen/tse/helpers"
	"strconv"
)

type AggregatedWorkItem struct {
	Description string
	AppRef      string
	JiraRef     string
	Project     string
	TotalTime   float64
}

type TimeSpentPerDay struct {
	Day       string
	Date      string
	TimeSpent float64
}

func AggregateWorkItems(items []KiaraWorkItem) ([]AggregatedWorkItem, error) {
	aggregationMap := make(map[string]*AggregatedWorkItem)

	for _, item := range items {
		key := fmt.Sprintf("%s|%s|%s|%s", item.Description, item.AppRef, item.JiraRef, item.Project)
		if _, exists := aggregationMap[key]; !exists {
			aggregationMap[key] = &AggregatedWorkItem{
				Description: item.Description,
				AppRef:      item.AppRef,
				JiraRef:     item.JiraRef,
				Project:     item.Project,
				TotalTime:   0,
			}
		}
		timeSpentFloat, err := timeSpentToFloat(item.TimeSpent)
		if err != nil {
			return nil, fmt.Errorf("%s %w", help.NewErrorStackTraceString("failed to aggregate work items"), err)
		}
		aggregationMap[key].TotalTime += timeSpentFloat
	}

	var aggregatedItems []AggregatedWorkItem
	for _, aggItem := range aggregationMap {
		aggregatedItems = append(aggregatedItems, *aggItem)
	}

	return aggregatedItems, nil
}

func timeSpentToFloat(timeSpent string) (float64, error) {
	timeSpentFloat, err := strconv.ParseFloat(timeSpent, 64)
	if err != nil {
		return 0, fmt.Errorf("%s %w", help.NewErrorStackTraceString("failed to parse time spent"), err)
	}
	return timeSpentFloat, nil
}

func CalculateTotalTimeSpent(entries []KiaraWorkItem) (float64, error) {
	totalTimeSpent := 0.0
	for _, entry := range entries {
		timeSpentFloat, err := timeSpentToFloat(entry.TimeSpent)
		if err != nil {
			return 0, fmt.Errorf("%s %w", help.NewErrorStackTraceString("failed to calculate total time spent"), err)
		}
		totalTimeSpent += timeSpentFloat
	}
	return totalTimeSpent, nil
}

func CalculateTimeSpentPerDay(entries []KiaraWorkItem) ([]TimeSpentPerDay, error) {
	timeSpentMap := make(map[string]map[string]float64)

	for _, entry := range entries {
		timeSpentFloat, err := timeSpentToFloat(entry.TimeSpent)
		if err != nil {
			return nil, fmt.Errorf("%s %w", "failed to calculate time spent per day", err)
		}
		if _, exists := timeSpentMap[entry.Day]; !exists {
			timeSpentMap[entry.Day] = make(map[string]float64)
		}
		timeSpentMap[entry.Day][entry.Date] += timeSpentFloat
	}

	timeSpentPerDay := make([]TimeSpentPerDay, 0)
	for day, dateMap := range timeSpentMap {
		for date, totalTimeSpent := range dateMap {
			timeSpentPerDay = append(timeSpentPerDay, TimeSpentPerDay{Day: day, Date: date, TimeSpent: totalTimeSpent})
		}
	}

	return timeSpentPerDay, nil
}
