package excel

import (
	"fmt"
	help "github.com/williamvannuffelen/tse/helpers"
	"github.com/xuri/excelize/v2"
)

func OpenExcelFile(fileName string) (*excelize.File, error) {
	log.Debug(fmt.Sprintf("Opening file '%s'", fileName))
	excelFile, err := excelize.OpenFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("%s %w", help.NewErrorStackTraceString("failed to open excel file"), err)
	}
	return excelFile, nil
}

func SaveExcelFile(excelFile *excelize.File) error {
	newRow := []interface{}{"Mo", "2024-10-05", "new task", "OPS-999", 3, "New Project", "AppRef"}
	if err := AppendRow(excelFile, "2024-10-21", newRow); err != nil {
		return fmt.Errorf("%s %w", help.NewErrorStackTraceString("failed to save excel file"), err)
	}
	if err := excelFile.SaveAs("ebase.xlsx"); err != nil {
		fmt.Println("Error saving file:", err)
		return fmt.Errorf("%s %w", help.NewErrorStackTraceString("failed to save excel file"), err)
	}
	return nil
}
