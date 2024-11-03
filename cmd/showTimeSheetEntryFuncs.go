package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/williamvannuffelen/tse/config"
	"github.com/williamvannuffelen/tse/excel"
	help "github.com/williamvannuffelen/tse/helpers"
	"github.com/williamvannuffelen/tse/workitem"
	"github.com/xuri/excelize/v2"
)

// todo: generalize and centralize reused function
func getFlagValues(cmd *cobra.Command) map[string]interface{} {
	boolFlags := []string{"no-week"}
	values := make(map[string]interface{})
	for _, flag := range boolFlags {
		value, _ := cmd.Flags().GetBool(flag)
		values[flag] = value
		log.Debug(fmt.Sprintf("%s: %t", flag, value))
	}

	stringFlags := []string{"output", "date", "day"}
	for _, flag := range stringFlags {
		value, _ := cmd.Flags().GetString(flag)
		values[flag] = value
		log.Debug(fmt.Sprintf("%s: %s", flag, value))
	}

	return values
}

// todo: generalize and centralize reused function
func setDefaultOutputFormat(values map[string]interface{}, appConfig config.Config) {
	if values["output"] == "" {
		values["output"] = appConfig.Keywords.DefaultOutputFormat
	}
}

func setDefaultValues(values map[string]interface{}) error {
	err := setDefaultDate(values)
	if err != nil {
		return fmt.Errorf("%s %w", help.NewErrorStackTraceString("failed to set default values"), err)
	}
	return nil
}

func setDefaultDate(values map[string]interface{}) error {
	if values["day"] != "" {
		dateVal, err := help.GetDateFromDay(values["day"].(string))
		if err != nil {
			return fmt.Errorf("%s %w", help.NewErrorStackTraceString("failed to set default date"), err)
		}
		values["date"] = dateVal
	} else {
		if values["date"] == "" {
			values["date"] = help.GetCurrentWeekDate()
		}
	}
	return nil
}

// better moved to excel package
func getTimeSheetEntries(fileName string, sheetName string) ([]workitem.KiaraWorkItem, error) {
	excelFile, err := excel.SelectTargetSheet(fileName, sheetName)
	if err != nil {
		return nil, fmt.Errorf("%s %w", help.NewErrorStackTraceString("failed to select target sheet"), err)
	}
	entries, err := readTimeSheetEntries(excelFile, sheetName)
	if err != nil {
		return nil, fmt.Errorf("%s %w", help.NewErrorStackTraceString("failed to read timesheet entries"), err)
	}
	return entries, nil
}

func readTimeSheetEntries(excelFile *excelize.File, sheetName string) ([]workitem.KiaraWorkItem, error) {
	rows, err := excelFile.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("%s %w", help.NewErrorStackTraceString("failed to read rows"), err)
	}

	var entries []workitem.KiaraWorkItem
	for _, row := range rows[1:] {
		entry := workitem.KiaraWorkItem{
			Day:         getValue(row, 0),
			Date:        getValue(row, 1),
			Description: getValue(row, 2),
			JiraRef:     getValue(row, 3),
			TimeSpent:   getValue(row, 4),
			Project:     getValue(row, 5),
			AppRef:      getValue(row, 6),
		}
		entries = append(entries, entry)
	}
	return entries, nil
}

func getValue(row []string, index int) string {
	if index < len(row) {
		return row[index]
	}
	return ""
}
