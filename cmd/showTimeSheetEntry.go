package cmd

import (
	//"os"
	"fmt"
	"github.com/spf13/cobra"
	help "github.com/williamvannuffelen/tse/helpers"
	"github.com/williamvannuffelen/tse/prettyprint"
	"github.com/williamvannuffelen/tse/workitem"
)

var showTimeSheetEntryCmd = &cobra.Command{
	Use:           "show-timesheet-entry",
	Aliases:       []string{"show"},
	Short:         "Alias: show - Show timesheet entries",
	Long:          "Shows timesheet entries present in the timesheet.",
	SilenceUsage:  true,
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {
		values := getFlagValues(cmd)
		setDefaultOutputFormat(values, appConfig)

		if values["date"] != "" {
			startOfWeek, err := help.GetStartOfWeek(values["date"].(string))
			if err != nil {
				log.Error(err)
				return
			}
			log.Debug("Start of week: ", startOfWeek)
		}

		err := setDefaultValues(values)
		if err != nil {
			log.Error(err)
			return
		}
		log.Debug("Date: ", values["date"])

		workItems, err := getTimeSheetEntries(appConfig.File.TargetFilePath, values["date"].(string))
		if err != nil {
			log.Error(err)
			return
		}
		log.Debug("WorkItems: ", workItems)

		totalTimeSpent, err := workitem.CalculateTotalTimeSpent(workItems)
		if err != nil {
			log.Error(err)
			return
		}
		log.Debug("Total time spent: ", totalTimeSpent)
		timeSpentPerDay, err := workitem.CalculateTimeSpentPerDay(workItems)
		if err != nil {
			log.Error(err)
			return
		}
		log.Debug("Time spent per day: ", timeSpentPerDay)

		timeSpentPerTask, err := workitem.CalculateTimeSpentPerTask(workItems)
		if err != nil {
			log.Error(err)
			return
		}
		log.Debug("Time spent per task: ", timeSpentPerTask)

		aggregatedWorkItems, err := workitem.AggregateWorkItems(workItems)
		if err != nil {
			log.Error(err)
			return
		}
		log.Debug("Aggregated work items: ", aggregatedWorkItems)
		log.Debug(len(aggregatedWorkItems))

		for _, item := range aggregatedWorkItems {
			prettyprint.PrintStructFields(item)
		}

		prettyprint.PrintTimeSpentPerDayTable(timeSpentPerDay)
		prettyprint.PrintTimeSpentWeekTotal(totalTimeSpent)
		prettyprint.PrintAggregatedWorkItemTable(aggregatedWorkItems, false, false, true)

		if values["no-week"] == false {
			fmt.Println("Showing entire week")
		} else {
			fmt.Println("Showing only the date")
		}

	},
}

func init() {
	showTimeSheetEntryCmd.Flags().BoolP("help", "h", false, "Display this help message")
	showTimeSheetEntryCmd.Flags().StringP("output", "o", "", "Output format. Options: json, pp (pretty print).")
	showTimeSheetEntryCmd.Flags().StringP("date", "d", "", "Date to show timesheet entries for. Format: yyyy-MM-dd. i.e.: 2024-11-18")
	showTimeSheetEntryCmd.Flags().StringP("day", "D", "", "Day to show timesheet entries for: mon | tue | wed | thu | fri | sat | sun .")
	showTimeSheetEntryCmd.Flags().BoolP("no-week", "w", false, "Show only the provided date instead of the entire week")
}
