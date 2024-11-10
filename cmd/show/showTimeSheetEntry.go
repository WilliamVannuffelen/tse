package show

import (
	//"os"
	"fmt"
	"github.com/spf13/cobra"
	logger "github.com/williamvannuffelen/go_zaplogger_iso8601"
	"github.com/williamvannuffelen/tse/config"
	help "github.com/williamvannuffelen/tse/helpers"
	"github.com/williamvannuffelen/tse/prettyprint"
	"github.com/williamvannuffelen/tse/workitem"
	"time"
)

var log logger.Logger

func SetLogger(l logger.Logger) {
	log = l
}

var appConfig = config.InitConfig()

var ShowTimeSheetEntryCmd = &cobra.Command{
	Use:           "show-timesheet-entry",
	Aliases:       []string{"show"},
	Short:         "Alias: show - Show timesheet entries",
	Long:          "Shows timesheet entries present in the timesheet.",
	SilenceUsage:  true,
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {
		var startOfWeek string
		var err error

		// Prep input values
		values := getFlagValues(cmd)
		setDefaultOutputFormat(values, appConfig)

		if values["date"] != "" {
			log.Debug("setting start of week")
			startOfWeek, err = help.GetStartOfWeek(values["date"].(string))
			if err != nil {
				log.Error(err)
				return
			}
			log.Debug("Start of week: ", startOfWeek)
		}

		err = setDefaultValues(values)
		if err != nil {
			log.Error(err)
			return
		}
		log.Debug("Date: ", values["date"])
		if startOfWeek == "" {
			log.Debug("Fetching start of current week")
			startOfWeek = help.GetCurrentWeekDate(time.Now)
		}

		log.Debug("Start of week: ", startOfWeek)
		// Get & process data
		workItems, err := getTimeSheetEntries(appConfig.File.TargetFilePath, startOfWeek)
		if err != nil {
			log.Error(err)
			return
		}
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

		aggregatedWorkItems, err := workitem.AggregateWorkItems(workItems)
		log.Debug(aggregatedWorkItems)
		if err != nil {
			log.Error(err)
			return
		}

		// Print output
		if values["day"] != "" {
			setDateIfDayProvided(values)
			prettyprint.PrintDayInSelectedFormat(values, timeSpentPerDay, startOfWeek, workItems, aggregatedWorkItems)
		}
		if values["no-week"] == false {
			prettyprint.PrintWeekInSelectedFormat(values, timeSpentPerDay, startOfWeek, workItems, aggregatedWorkItems)
		} else if values["day"] == "" {
			prettyprint.PrintDayInSelectedFormat(values, timeSpentPerDay, startOfWeek, workItems, aggregatedWorkItems)
		}
	},
}

func init() {
	ShowTimeSheetEntryCmd.Flags().BoolP("help", "h", false, "Display this help message")
	ShowTimeSheetEntryCmd.Flags().StringP("output", "o", appConfig.ShowTimeSheetEntry.DefaultOutputFormat, "Output format. Options: json, pp (pretty print).")
	ShowTimeSheetEntryCmd.Flags().StringP("date", "d", "", "Date to show timesheet entries for. Format: yyyy-MM-dd. i.e.: 2024-11-18")
	ShowTimeSheetEntryCmd.Flags().StringP("day", "D", "", "Day to show timesheet entries for: mon | tue | wed | thu | fri | sat | sun .")
	ShowTimeSheetEntryCmd.Flags().BoolP("no-week", "w", false, "Show only the provided date instead of the entire week")
	ShowTimeSheetEntryCmd.Flags().BoolP("hide-appref", "a", appConfig.ShowTimeSheetEntry.HideAppRef, "Hide the AppRef column")
	ShowTimeSheetEntryCmd.Flags().BoolP("hide-jiraref", "j", appConfig.ShowTimeSheetEntry.HideJiraRef, "Hide the JiraRef column")
	ShowTimeSheetEntryCmd.Flags().BoolP("hide-project", "p", appConfig.ShowTimeSheetEntry.HideProject, "Hide the Project column")
}

func setDateIfDayProvided(values map[string]interface{}) error {
	values["no-week"] = true
	date, err := help.GetDateFromDay(values["day"].(string), time.Now)
	if err != nil {
		return fmt.Errorf("%s %s", help.NewErrorStackTraceString("failed to process provided day value"), err)
	}
	values["date"] = date
	return nil
}
