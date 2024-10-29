package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/williamvannuffelen/tse/config"
	"github.com/williamvannuffelen/tse/excel"
	help "github.com/williamvannuffelen/tse/helpers"
	"github.com/williamvannuffelen/tse/keywords"
	"github.com/williamvannuffelen/tse/workitem"
)

var appConfig config.Config

var addTimeSheetEntryCmd = &cobra.Command{
	Use:           "addTimeSheetEntry",
	Short:         "Add timesheet entry",
	Long:          "Adds a timesheet entry to the timesheet.",
	SilenceUsage:  true,
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {
		flags := []string{"date", "description", "jira-ref", "time", "project", "app-ref", "keyword"}
		values := make(map[string]string)

		for _, flag := range flags {
			value, _ := cmd.Flags().GetString(flag)
			values[flag] = value
			log.Debug(fmt.Sprintf("%s: %s", flag, value))
		}

		//TODO: add case here for keyword, fullkeyword
		date := values["date"]
		description := values["description"]
		jiraRef := values["jira-ref"]
		time := values["time"]
		project := values["project"]
		appRef := values["app-ref"]
		basicKeyword := values["basic-keyword"]
		keyword := values["keyword"]

		if keyword != "" {
			log.Debug("Keyword provided.")
			// TODO: source keywords paths from config?
			kw, err := keywords.MatchKeywords("./keywords/keywords_exact.json", keyword)
			if err != nil {
				fmt.Errorf("%s %w", help.NewErrorStackTraceString(fmt.Sprintf("failed to get info for keyword '%s'", keyword)), err)
			}
			//TODO: add appref to fullkeyword
			description = kw.(keywords.FullKeyword).Description
			jiraRef = kw.(keywords.FullKeyword).JiraRef
			project = kw.(keywords.FullKeyword).Project
		}

		if basicKeyword != "" {
			log.Debug("Basic Keyword provided.")
			// TODO: source keywords paths from config?
			kw, err := keywords.MatchKeywords("./keywords/keywords.json", basicKeyword)
			if err != nil {
				fmt.Errorf("%s %w", help.NewErrorStackTraceString(fmt.Sprintf("failed to get info for keyword '%s'", basicKeyword)), err)
			}
			jiraRef = kw.(keywords.BasicKeyword).JiraRef
			project = kw.(keywords.BasicKeyword).Project
		}

		//TODO: error if one of permutations is not reached

		workItem := CreateKiaraWorkItem(appConfig, date, description, jiraRef, time, project, appRef)
		log.Debug(fmt.Sprintf("WorkItem: %+v", workItem))

		targetFilePath := viper.Get("File.targetFilePath").(string)

		sheetName, ok := viper.Get("File.targetSheetName").(string)
		if !ok || sheetName == "" {
			log.Debug("No sheet name provided, using current week's sheet name.")
			sheetName = excel.GetCurrentWeekSheetName()
		}
		log.Debug("sheetName: ", sheetName)
		templateSheetName := viper.Get("File.templateSheetName").(string)
		err := WriteTimeSheetEntry(targetFilePath, sheetName, templateSheetName, workItem)
		if err != nil {
			log.Error(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(addTimeSheetEntryCmd)
	addTimeSheetEntryCmd.Flags().BoolP("help", "h", false, "Display this help message")
	addTimeSheetEntryCmd.Flags().StringP("date", "d", "", "Date of the timesheet entry in yyyy-MM-dd format. Will default to today if not provided.")
	addTimeSheetEntryCmd.Flags().StringP("description", "D", "", "Description of the timesheet entry")
	addTimeSheetEntryCmd.Flags().StringP("jira-ref", "j", "", "Jira reference of the timesheet entry. Will default to the value set in config.yaml")
	addTimeSheetEntryCmd.Flags().StringP("time", "t", "0", "Time spent, in hours, of the timesheet entry")
	addTimeSheetEntryCmd.Flags().StringP("project", "p", "", "Project of the timesheet entry. Will default to the value set in config.yaml")
	addTimeSheetEntryCmd.Flags().StringP("app-ref", "a", "", "App reference of the timesheet entry. Will default to the value set in config.yaml")
	addTimeSheetEntryCmd.Flags().StringP("keyword", "k", "", "Keyword of the timesheet entry. Used to source full description, project, jira-ref and app-ref for known tasks.")
	addTimeSheetEntryCmd.Flags().StringP("basic-keyword", "K", "", "Basic keyword of the timesheet entry. Used to source, project, jira-ref and app-ref for known tasks.")
	addTimeSheetEntryCmd.Flags().StringP("sheet", "s", "", "Sheet name to write the timesheet entry to. Will default to the current week's sheet name.")
}

func CreateKiaraWorkItem(appConfig config.Config, date string, description string, jiraRef string, time string, project string, appRef string) *workitem.KiaraWorkItem {
	log.Debug("Project: ", project)
	log.Debug("default projecT: ", appConfig.Project.DefaultProjectName)
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
