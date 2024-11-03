package add

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	logger "github.com/williamvannuffelen/go_zaplogger_iso8601"
	"github.com/williamvannuffelen/tse/config"
	"github.com/williamvannuffelen/tse/excel"
	help "github.com/williamvannuffelen/tse/helpers"
	"github.com/williamvannuffelen/tse/keywords"
	"github.com/williamvannuffelen/tse/workitem"
	"strconv"
	"time"
)

var log logger.Logger

func SetLogger(l logger.Logger) {
	log = l
}

var appConfig = config.InitConfig()

func ProcessKeywords(appConfig config.Config, values map[string]string) (map[string]string, error) {
	log.Debug("inside processkw")
	if values["keyword"] != "" {
		log.Debug("Keyword provided. Fetching values.")
		keywordValues, err := keywords.MatchAndExtractKeywords("./keywords/keywords.json", values["keyword"], "addTimeSheetEntry")
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
	log.Debug("Validating input")
	if processedValues["description"] == "" {
		return fmt.Errorf("no description provided. Provide one using -d or use a keyword with -k or -K")
	}
	if processedValues["date"] != "" {
		_, err := time.Parse("2006-01-02", processedValues["date"])
		if err != nil {
			return fmt.Errorf("invalid date format. Please use yyyy-MM-dd. e.g. 2024-09-31")
		}
	}
	if processedValues["timespent"] == "" {
		return fmt.Errorf("no timespent provided. Provide one using -t")
	}
	if processedValues["timespent"] != "" {
		_, err := strconv.ParseFloat(processedValues["timespent"], 64)
		if err != nil {
			return fmt.Errorf("invalid timespent format. Please use a number. e.g. 8")
		}
		if processedValues["timespent"] == "0" {
			return fmt.Errorf("timespent cannot be 0. Please provide a valid time with flag -t")
		}
	}
	// other values are optional
	return nil
}

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

		workItem := CreateKiaraWorkItem(
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
			sheetName = excel.GetCurrentWeekSheetName()
		}
		templateSheetName := appConfig.File.TemplateSheetName
		err = WriteTimeSheetEntry(targetFilePath, sheetName, templateSheetName, workItem)
		if err != nil {
			log.Error(err)
		}
		PrintWorkItem(workItem)
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

func CreateKiaraWorkItem(appConfig config.Config, date string, description string, jiraRef string, timespent string, project string, appRef string) *workitem.KiaraWorkItem {
	workItem := workitem.NewKiaraWorkItem(appConfig, date, description, jiraRef, timespent, project, appRef)
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

func PrintWorkItem(workItem *workitem.KiaraWorkItem) {
	keyColor := color.New(color.FgCyan).SprintFunc()
	valueColor := color.New(color.FgYellow).SprintFunc()
	fmt.Println("Added timesheet entry:")
	fmt.Println(keyColor("  \"Day\":"), valueColor(workItem.Day))
	fmt.Println(keyColor("  \"Date\":"), valueColor(workItem.Date))
	fmt.Println(keyColor("  \"Description\":"), valueColor(workItem.Description))
	fmt.Println(keyColor("  \"JiraRef\":"), valueColor(workItem.JiraRef))
	fmt.Println(keyColor("  \"TimeSpent\":"), valueColor(workItem.TimeSpent))
	fmt.Println(keyColor("  \"Project\":"), valueColor(workItem.Project))
	fmt.Println(keyColor("  \"AppRef\":"), valueColor(workItem.AppRef))
}
