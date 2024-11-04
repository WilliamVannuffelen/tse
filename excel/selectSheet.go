package excel

import (
	"fmt"
	help "github.com/williamvannuffelen/tse/helpers"
	"github.com/xuri/excelize/v2"
)

func FindSheetIndex(excelFile *excelize.File, sheetName string) (int, error) {
	log.Debug(fmt.Sprintf("Finding sheet index for sheet '%s'", sheetName))
	index, err := excelFile.GetSheetIndex(sheetName)
	if err != nil {
		return -1, fmt.Errorf("%s %w", help.NewErrorStackTraceString(fmt.Sprintf("failed to grab sheet index for %s", sheetName)), err)
	}

	return index, nil
}

func FindLastSheetIndex(excelFile *excelize.File) (int, error) {
	log.Debug("Finding last sheet index")
	sheetList := excelFile.GetSheetList()
	lastSheetIndex := len(sheetList) - 1
	if lastSheetIndex < 0 {
		return -1, fmt.Errorf("%s %w", help.NewErrorStackTraceString("failed to find last sheet index"), fmt.Errorf("no sheets found"))
	}
	return lastSheetIndex, nil
}

func FindTemplateSheet(excelFile *excelize.File, templateSheetName string) (int, error) {
	if templateSheetName == "" {
		log.Debug(("Template sheet name not provided. Assuming it's the first sheet."))
		return 0, nil
	}
	log.Debug(fmt.Sprintf("Finding template sheet index for sheet '%s'", templateSheetName))
	index, err := excelFile.GetSheetIndex(templateSheetName)
	if err != nil {
		return -1, err
	}
	return index, nil
}

func SelectSheet(excelFile *excelize.File, sheetName string, sheetIndex int) error {
	excelFile.SetActiveSheet(sheetIndex)
	log.Debug(fmt.Sprintf("Selected sheet '%s' at index '%d'", sheetName, sheetIndex))
	return nil
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
