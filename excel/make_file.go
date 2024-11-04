package excel

import (
	"fmt"
	help "github.com/williamvannuffelen/tse/helpers"
	"github.com/xuri/excelize/v2"
)

// func CreateSheet (file *excelize.File, date string) (int, error) {
// 	index, err := file.NewSheet(date)
// 	if err != nil {
// 		return -1, fmt.Errorf("%s %w", help.NewErrorStackTraceString("failed to create sheet"), err)
// 	}
// 	columns := []string{"A", "B", "C", "D", "E", "F", "G"}
// 	columnNames := []string{"Day", "Date", "Description", "JiraRef", "TimeSpent", "Project", "AppRef"}
// 	for i, column := range columns {
// 		err = file.SetCellValue(date, column + "1", columnNames[i])
// 		if err != nil {
// 			return -1, fmt.Errorf("%s %w", help.NewErrorStackTraceString("failed to set cell value"), err)
// 		}
// 	}
// 	file.SetActiveSheet(index)
// 	return index, nil
// }

func CreateExcelFile(filePath string, date string) error {
	file := excelize.NewFile()
	defer file.Close()

	index, err := MakeSheetFromScratch(file, date)
	if err != nil {
		return fmt.Errorf("%s %w", help.NewErrorStackTraceString("failed to create sheet"), err)
	}

	file.SetActiveSheet(index)

	if err := file.SaveAs(filePath); err != nil {
		return fmt.Errorf("%s %w", help.NewErrorStackTraceString("failed to save excel file"), err)
	}
	log.Debug(fmt.Sprintf("Created new excel file %s and sheet %s", filePath, date))
	return nil
}
