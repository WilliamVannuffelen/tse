package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	//logger "github.com/williamvannuffelen/go_zaplogger_iso8601"
	"github.com/williamvannuffelen/tse/config"
	"github.com/williamvannuffelen/tse/excel"
	help "github.com/williamvannuffelen/tse/helpers"
	"github.com/williamvannuffelen/tse/keywords"
	"github.com/williamvannuffelen/tse/workitem"
)

//var log logger.Logger
// var err error

// func SetLogger(l logger.Logger) {
// 	log = l
// }

func MatchAndExtractKeywords(filePath string, keyword string) (map[string]string, error) {
	log.Debug("Matching keywords for keyword: ", keyword)
	kw, err := keywords.MatchKeywords(filePath, keyword)
	if err != nil {
		return nil, fmt.Errorf("%s %w", help.NewErrorStackTraceString(fmt.Sprintf("failed to get info for provided keyword '%s'", keyword)), err)
	}
	keywordValues := map[string]string{
		"description": kw.Description,
		"jira-ref":    kw.JiraRef,
		"project":     kw.Project,
	}
	log.Debug("Got values from keyword: ", keywordValues["description"], keywordValues["jiraRef"], keywordValues["project"])
	return keywordValues, nil
}

func ProcessKeywords(appConfig config.Config, values map[string]string) (map[string]string, error) {
	log.Debug("inside processkw")
	if values["keyword"] != "" {
		log.Debug("Keyword provided. Fetching values.")
		keywordValues, err := MatchAndExtractKeywords("./keywords/keywords_exact.json", values["keyword"])
		if err != nil {
			return nil, fmt.Errorf("%s %w", help.NewErrorStackTraceString(fmt.Sprintf("failed to get info for keyword '%s'", values["keyword"])), err)
		}
		values["description"] = keywordValues["description"]
		values["jira-ref"] = keywordValues["jira-ref"]
		values["project"] = keywordValues["project"]
	}
	if values["basic-keyword"] != "" {
		log.Debug("Basic keyword provided. Fetching values.")
		keywordValues, err := MatchAndExtractKeywords("./keywords/keywords.json", values["basic-keyword"])
		if err != nil {
			return nil, fmt.Errorf("%s %w", help.NewErrorStackTraceString(fmt.Sprintf("failed to get info for keyword '%s'", values["basic-keyword"])), err)
		}
		values["jira-ref"] = keywordValues["jira-ref"]
		values["project"] = keywordValues["project"]
	}
	// if values["project"] == "" {
	// 	log.Debug("No project provided. Using default project.")
	// 	values["project"] = appConfig.Project.DefaultProjectName
	// }
	return values, nil
}

var addTimeSheetEntryCmd = &cobra.Command{
	Use:           "addTimeSheetEntry",
	Short:         "Add timesheet entry",
	Long:          "Adds a timesheet entry to the timesheet.",
	SilenceUsage:  true,
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {
		// appConfig := config.InitConfig()
		// log, err := logger.InitLogger("log.txt", appConfig.General.LogLevel)
		// if err != nil {
		// 	log.Warn(err)
		// }
		// SetLogger(log) // not redundant: required for funcs inside cmd package, but outside of Run
		// workitem.SetLogger(log)
		// excel.SetLogger(log)
		// keywords.SetLogger(log)

		// TODO: Add validation for date fmt = yyyy-MM-dd
		// TODO: add validation for time fmt: 1 or 1.5 - not 1,5
		log.Debug("processing flags")
		flags := []string{"date", "description", "jira-ref", "time", "project", "app-ref", "keyword", "basic-keyword"}
		values := make(map[string]string)
		for _, flag := range flags {
			value, _ := cmd.Flags().GetString(flag)
			values[flag] = value
			log.Debug(fmt.Sprintf("%s: %s", flag, value))
		}
		processedValues, err := ProcessKeywords(appConfig, values)
		if err != nil {
			fmt.Errorf("%s %w", help.NewErrorStackTraceString(fmt.Sprintf("failed to process keyword info")), err)
		}
		log.Debug("Processed values: ", processedValues)

		if processedValues["description"] == "" {
			log.Error("No description provided. Provide one using -d or use a keyword with -k or -K.")
			return
		}

		workItem := CreateKiaraWorkItem(
			appConfig,
			processedValues["date"],
			processedValues["description"],
			processedValues["jiraRef"],
			processedValues["time"],
			processedValues["project"],
			processedValues["appRef"],
		)
		log.Debug(fmt.Sprintf("WorkItem: %+v", workItem))

		targetFilePath := appConfig.File.TargetFilePath
		sheetName := appConfig.File.TargetSheetName
		if sheetName == "" {
			sheetName = excel.GetCurrentWeekSheetName()
		}
		templateSheetName := appConfig.File.TemplateSheetName
		err = WriteTimeSheetEntry(targetFilePath, sheetName, templateSheetName, workItem)
		if err != nil {
			log.Error(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(addTimeSheetEntryCmd)
	addTimeSheetEntryCmd.Flags().BoolP("help", "h", false, "Display this help message")
	addTimeSheetEntryCmd.Flags().StringP("date", "D", "", "Date of the timesheet entry in yyyy-MM-dd format. Will default to today if not provided.")
	addTimeSheetEntryCmd.Flags().StringP("description", "d", "", "Description of the timesheet entry")
	addTimeSheetEntryCmd.Flags().StringP("jira-ref", "j", "", "Jira reference of the timesheet entry. Will default to the value set in config.yaml")
	addTimeSheetEntryCmd.Flags().StringP("time", "t", "0", "Time spent, in hours, of the timesheet entry")
	addTimeSheetEntryCmd.Flags().StringP("project", "p", "", "Project of the timesheet entry. Will default to the value set in config.yaml")
	addTimeSheetEntryCmd.Flags().StringP("app-ref", "a", "", "App reference of the timesheet entry. Will default to the value set in config.yaml")
	addTimeSheetEntryCmd.Flags().StringP("keyword", "k", "", "Keyword of the timesheet entry. Used to source full description, project, jira-ref and app-ref for known tasks.")
	addTimeSheetEntryCmd.Flags().StringP("basic-keyword", "K", "", "Basic keyword of the timesheet entry. Used to source, project, jira-ref and app-ref for known tasks.")
	addTimeSheetEntryCmd.Flags().StringP("sheet", "s", "", "Sheet name to write the timesheet entry to. Will default to the current week's sheet name.")
}

func CreateKiaraWorkItem(appConfig config.Config, date string, description string, jiraRef string, time string, project string, appRef string) *workitem.KiaraWorkItem {
	workItem := workitem.NewKiaraWorkItem(appConfig, date, description, jiraRef, time, project, appRef)
	return workItem
}

func WriteTimeSheetEntry(fileName string, sheetName string, templateSheetName string, workItem *workitem.KiaraWorkItem) error {
	excelFile, err := excel.SetTargetSheet(fileName, sheetName, templateSheetName)
	if err != nil {
		return fmt.Errorf("%s %w", help.NewErrorStackTraceString("failed to set active sheet"), err)
	}
	log.Debug("sheetName: ", sheetName)
	err = excel.AddNewTimesheetEntry(excelFile, sheetName, workItem)
	if err != nil {
		return fmt.Errorf("%s %w", help.NewErrorStackTraceString("failed to add timesheet entry"), err)
	}
	return nil
}
