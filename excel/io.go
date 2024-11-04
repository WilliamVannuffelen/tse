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
	defer excelFile.Close()

	return excelFile, nil
}

func CreateExcelFile(filePath string, date string) (*excelize.File, error) {
	file := excelize.NewFile()

	index, err := MakeSheetFromScratch(file, date)
	if err != nil {
		return nil, fmt.Errorf("%s %w", help.NewErrorStackTraceString("failed to create sheet"), err)
	}

	file.SetActiveSheet(index)
	return file, nil
}

func SaveExcelFile(excelFile *excelize.File, fileName string) error {
	if err := excelFile.SaveAs(fileName); err != nil {
		return fmt.Errorf("%s %w", help.NewErrorStackTraceString("failed to save excel file"), err)
	}
	return nil
}
