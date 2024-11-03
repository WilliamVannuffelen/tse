package excel

import (
	"fmt"
	logger "github.com/williamvannuffelen/go_zaplogger_iso8601"
	help "github.com/williamvannuffelen/tse/helpers"
	"github.com/williamvannuffelen/tse/workitem"
	"github.com/xuri/excelize/v2"
)

var log logger.Logger

func SetLogger(l logger.Logger) {
	log = l
}

func SelectTargetSheet(fileName string, sheetName string) (*excelize.File, error) {
	excelFile, err := OpenExcelFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("%s %w", help.NewErrorStackTraceString("failed to set active sheet"), err)
	}
	sheetIndex, err := FindSheetIndex(excelFile, sheetName)
	if err != nil {
		return nil, fmt.Errorf("%s %w", help.NewErrorStackTraceString("failed to set active sheet"), err)
	}
	if sheetIndex == -1 {
		return nil, fmt.Errorf("%s %w", help.NewErrorStackTraceString("failed to set active sheet"), fmt.Errorf("sheet not found"))
	} else {
		SelectSheet(excelFile, sheetName, sheetIndex)
	}
	log.Debug(fmt.Sprintf("Selected target sheet %s at index %d in file %s", sheetName, sheetIndex, fileName))
	return excelFile, nil
}

// rename -> does more than setting target sheet
func SetTargetSheet(fileName string, sheetName string, templateSheetName string) (*excelize.File, error) {
	excelFile, err := OpenExcelFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("%s %w", help.NewErrorStackTraceString("failed to set active sheet"), err)
	}

	if sheetName == "" {
		log.Debug("No sheet name provided, using current week's sheet name.")
		sheetName = GetCurrentWeekSheetName()
	}

	sheetIndex, err := FindSheetIndex(excelFile, sheetName)
	if err != nil {
		return nil, fmt.Errorf("%s %w", help.NewErrorStackTraceString("failed to set active sheet"), err)
	}
	if sheetIndex == -1 {
		log.Info(fmt.Sprintf("Sheet '%s' not found in file '%s'", sheetName, fileName))
		newSheetIndex, err := MakeSheetFromTemplate(excelFile, sheetName, templateSheetName)
		if err != nil {
			return nil, fmt.Errorf("%s %w", help.NewErrorStackTraceString("failed to make target sheet"), err)
		}
		SelectSheet(excelFile, sheetName, newSheetIndex)
	} else {
		SelectSheet(excelFile, sheetName, sheetIndex)
	}
	return excelFile, nil
}

func AddNewTimesheetEntry(excelFile *excelize.File, sheet string, workItem *workitem.KiaraWorkItem) error {
	row := []interface{}{
		workItem.Day,
		workItem.Date,
		workItem.Description,
		workItem.JiraRef,
		workItem.TimeSpent,
		workItem.Project,
		workItem.AppRef,
	}
	err := AppendRow(excelFile, sheet, row)
	if err != nil {
		return fmt.Errorf("%s %w", help.NewErrorStackTraceString("failed to append row to sheet"), err)
	}
	err = SaveExcelFile(excelFile)
	if err != nil {
		return fmt.Errorf("%s %w", help.NewErrorStackTraceString("failed to save excel file"), err)
	}
	return nil
}
