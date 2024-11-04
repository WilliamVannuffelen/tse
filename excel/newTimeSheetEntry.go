package excel

import (
	"fmt"
	help "github.com/williamvannuffelen/tse/helpers"
	"github.com/williamvannuffelen/tse/workitem"
	"github.com/xuri/excelize/v2"
	"strconv"
)

func AppendRow(excelFile *excelize.File, sheet string, row []interface{}) error {
	rows, err := excelFile.GetRows(sheet)
	if err != nil {
		return fmt.Errorf("%s %w", help.NewErrorStackTraceString("failed to get rows to append to"), err)
	}

	nextRow := len(rows) + 1
	axis := "A" + strconv.Itoa(nextRow)

	err = excelFile.SetSheetRow(sheet, axis, &row)
	if err != nil {
		return fmt.Errorf("%s %w", help.NewErrorStackTraceString(fmt.Sprintf("failed to write entry to file: %s", row)), err)
	}
	return err
}

func AddNewTimesheetEntry(excelFile *excelize.File, sheet string, workItem *workitem.KiaraWorkItem, fileName string) error {
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
	err = SaveExcelFile(excelFile, fileName)
	if err != nil {
		return fmt.Errorf("%s %w", help.NewErrorStackTraceString("failed to save excel file"), err)
	}
	return nil
}
