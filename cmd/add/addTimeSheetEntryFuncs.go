package add

import (
	"fmt"
	"github.com/williamvannuffelen/tse/config"
	"github.com/williamvannuffelen/tse/excel"
	help "github.com/williamvannuffelen/tse/helpers"
	"github.com/williamvannuffelen/tse/keywords"
	"github.com/williamvannuffelen/tse/workitem"
	"strconv"
	"time"
)

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

func WriteTimeSheetEntry(fileName string, sheetName string, templateSheetName string, workItem *workitem.KiaraWorkItem) error {
	excelFile, err := excel.MakeOrSelectTargetSheet(fileName, sheetName, templateSheetName)
	if err != nil {
		log.Warn(fmt.Sprintf("failed to open file %s", fileName))
		log.Info("Creating new file: ", fileName)
		err = nil
		excelFile, err = excel.CreateExcelFile(fileName, sheetName)
		if err != nil {
			return fmt.Errorf("%s %w", help.NewErrorStackTraceString("failed to create new excel file"), err)
		}
	}
	log.Debug("sheetName: ", sheetName)
	err = excel.AddNewTimesheetEntry(excelFile, sheetName, workItem, fileName)
	if err != nil {
		return fmt.Errorf("%s %w", help.NewErrorStackTraceString("failed to add timesheet entry"), err)
	}
	return nil
}
