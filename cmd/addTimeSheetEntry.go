package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"strconv"
	"time"
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
		"app-ref":     kw.AppRef,
	}
	log.Debug(fmt.Sprintf("Got values from keyword: Description: '%s', JiraRef: '%s', Project: '%s', AppRef: '%s'",
		keywordValues["description"],
		keywordValues["jira-ref"],
		keywordValues["project"],
		keywordValues["app-ref"]))
	return keywordValues, nil
}

func ProcessKeywords(appConfig config.Config, values map[string]string) (map[string]string, error) {
	log.Debug("inside processkw")
	if values["keyword"] != "" {
		log.Debug("Keyword provided. Fetching values.")
		keywordValues, err := MatchAndExtractKeywords("./keywords/keywords_temp.json", values["keyword"])
		if err != nil {
			return nil, fmt.Errorf("%s %w", help.NewErrorStackTraceString(fmt.Sprintf("failed to get info for keyword '%s'", values["keyword"])), err)
		}

		// if flag was empty for property, use value from keyword
		for _, key := range []string{"description", "jira-ref", "project", "app-ref"} {
			if values[key] == "" {
				values[key] = keywordValues[key]
			}
		}
	}
	return values, nil
}

func ValidateInputValues(processedValues map[string]string) error {
	if processedValues["description"] == "" {
		return fmt.Errorf("no description provided. Provide one using -d or use a keyword with -k or -K")
	}
	if processedValues["date"] != "" {
		_, err := time.Parse("2006-01-02", processedValues["date"])
		if err != nil {
			return fmt.Errorf("invalid date format. Please use yyyy-MM-dd. e.g. 2024-09-31")
		}
	}
	if processedValues["time"] == "" {
		return fmt.Errorf("no time provided. Provide one using -t")
	}
	if processedValues["time"] != "" {
		_, err := strconv.ParseFloat(processedValues["time"], 64)
		if err != nil {
			return fmt.Errorf("invalid time format. Please use a number. e.g. 8")
		}
	}
	// other values are optional
	return nil
}

var addTimeSheetEntryCmd = &cobra.Command{
	Use:           "addTimeSheetEntry",
	Short:         "Add timesheet entry",
	Long:          "Adds a timesheet entry to the timesheet.",
	SilenceUsage:  true,
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {
		log.Debug("processing flags")
		flags := []string{"date", "description", "jira-ref", "time", "project", "app-ref", "keyword"}
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

		workItem := CreateKiaraWorkItem(
			appConfig,
			processedValues["date"],
			processedValues["description"],
			processedValues["jira-ref"],
			processedValues["time"],
			processedValues["project"],
			processedValues["app-ref"],
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
