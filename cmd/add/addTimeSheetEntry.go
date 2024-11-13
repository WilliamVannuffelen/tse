package add

import (
	"fmt"
	"github.com/spf13/cobra"
	logger "github.com/williamvannuffelen/go_zaplogger_iso8601"
	"github.com/williamvannuffelen/tse/config"
	help "github.com/williamvannuffelen/tse/helpers"
	"github.com/williamvannuffelen/tse/prettyprint"
	"github.com/williamvannuffelen/tse/workitem"
	"os"
)

var log logger.Logger

func SetLogger(l logger.Logger) {
	log = l
}

var appConfig = config.InitConfig()

var AddTimeSheetEntryCmd = &cobra.Command{
	Use:           "add-timesheet-entry",
	Aliases:       []string{"add"},
	Short:         "Alias: add - Add timesheet entry",
	Long:          "Adds a timesheet entry to the timesheet.",
	SilenceUsage:  true,
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {
		log.Debug("processing flags")
		flags := []string{"date", "description", "jira-ref", "timespent", "project", "app-ref", "keyword"}
		values := make(map[string]string)
		for _, flag := range flags {
			value, _ := cmd.Flags().GetString(flag)
			values[flag] = value
			log.Debug(fmt.Sprintf("%s: %s", flag, value))
		}
		log.Debug("Processing values")
		processedValues, err := ProcessKeywords(appConfig, values)
		if err != nil {
			log.Error(fmt.Errorf("%s %w", help.NewErrorStackTraceString("failed to process keyword info"), err))
			return
		}
		log.Debug("Processed values: ", processedValues)

		err = ValidateInputValues(processedValues)
		if err != nil {
			log.Error(fmt.Errorf("%s %w", help.NewErrorStackTraceString("failed to validate provided input values"), err))
			return
		}
		log.Debug("Input values validated.")

		workItem := workitem.NewKiaraWorkItem(
			appConfig,
			processedValues["date"],
			processedValues["description"],
			processedValues["jira-ref"],
			processedValues["timespent"],
			processedValues["project"],
			processedValues["app-ref"],
		)
		log.Debug(fmt.Sprintf("WorkItem: %+v", workItem))

		targetFilePath := appConfig.File.TargetFilePath
		sheetName := appConfig.File.TargetSheetName
		if sheetName == "" {
			//sheetName = help.GetCurrentWeekDate()
			sheetName, err = help.GetStartOfWeek(workItem.Date)
			if err != nil {
				log.Error(err)
				return
			}
		}
		templateSheetName := appConfig.File.TemplateSheetName
		err = WriteTimeSheetEntry(targetFilePath, sheetName, templateSheetName, workItem)
		if err != nil {
			log.Error(err)
			return
		}
		prettyprint.PrintWorkItem(os.Stdout, workItem)
	},
}

func init() {
	AddTimeSheetEntryCmd.Flags().BoolP("help", "h", false, "Display this help message")
	AddTimeSheetEntryCmd.Flags().StringP("date", "D", "", "Date of the timesheet entry in yyyy-MM-dd format. Will default to today if not provided.")
	AddTimeSheetEntryCmd.Flags().StringP("description", "d", "", "Description of the timesheet entry")
	AddTimeSheetEntryCmd.Flags().StringP("jira-ref", "j", "", "Jira reference of the timesheet entry. Will default to the value set in config.yaml")
	AddTimeSheetEntryCmd.Flags().StringP("timespent", "t", "0", "Time spent, in hours, of the timesheet entry")
	AddTimeSheetEntryCmd.Flags().StringP("project", "p", "", "Project of the timesheet entry. Will default to the value set in config.yaml")
	AddTimeSheetEntryCmd.Flags().StringP("app-ref", "a", "", "App reference of the timesheet entry. Will default to the value set in config.yaml")
	AddTimeSheetEntryCmd.Flags().StringP("keyword", "k", "", "Keyword of the timesheet entry. Used to source full description, project, jira-ref and app-ref for known tasks.")
	AddTimeSheetEntryCmd.Flags().StringP("sheet", "s", "", "Sheet name to write the timesheet entry to. Will default to the current week's sheet name.")
}
