package excel

import (
	"fmt"
	help "github.com/williamvannuffelen/tse/helpers"
	"github.com/xuri/excelize/v2"
)

func MakeSheetFromTemplate(excelFile *excelize.File, sheetName string, templateSheetName string) (int, error) {
	newSheetIndex, err := MakeSheet(excelFile, sheetName)
	if err != nil {
		return 0, fmt.Errorf("%s %w", help.NewErrorStackTraceString("failed to make target sheet"), err)
	}
	templateSheetIndex, err := FindTemplateSheet(excelFile, templateSheetName)
	if err != nil {
		return 0, fmt.Errorf("%s %w", help.NewErrorStackTraceString("failed to find template sheet"), err)
	}
	err = CopySheet(excelFile, templateSheetIndex, newSheetIndex)
	if err != nil {
		return 0, fmt.Errorf("%s %w", help.NewErrorStackTraceString("failed to copy sheet"), err)
	}
	return newSheetIndex, nil
}

func MakeSheet(excelFile *excelize.File, sheetName string) (int, error) {
	newSheetIndex, err := excelFile.NewSheet(sheetName)
	if err != nil {
		return 0, fmt.Errorf("%s %w", help.NewErrorStackTraceString("failed to create new sheet"), err)
	}
	log.Debug(fmt.Sprintf("Sheet created successfully: '%s' at index '%d'", sheetName, newSheetIndex))
	return newSheetIndex, nil
}

func MakeSheetFromScratch(file *excelize.File, date string) (int, error) {
	index, err := file.NewSheet(date)
	if err != nil {
		return -1, fmt.Errorf("%s %w", help.NewErrorStackTraceString("failed to create sheet"), err)
	}
	columns := []string{"A", "B", "C", "D", "E", "F", "G"}
	columnNames := []string{"Day", "Date", "Description", "JiraRef", "TimeSpent", "Project", "AppRef"}
	for i, column := range columns {
		err = file.SetCellValue(date, column+"1", columnNames[i])
		if err != nil {
			return -1, fmt.Errorf("%s %w", help.NewErrorStackTraceString("failed to set cell value"), err)
		}
	}
	file.SetActiveSheet(index)
	return index, nil
}

func CopySheet(excelFile *excelize.File, sourceSheetIndex int, targetSheetIndex int) error {
	if err := excelFile.CopySheet(sourceSheetIndex, targetSheetIndex); err != nil {
		return fmt.Errorf("%s %w", help.NewErrorStackTraceString("failed to copy sheet contents from template"), err)
	}
	log.Debug("Sheet copied successfully")
	log.Debug(excelFile.GetSheetName(targetSheetIndex))
	return nil
}
