package excel

import (
	"fmt"
  "github.com/xuri/excelize/v2"
)

func MakeSheetFromTemplate(excelFile *excelize.File, sheetName string, templateSheetName string) (int, error) {
	newSheetIndex, err := MakeSheet(excelFile, sheetName)
	if err != nil {
		return 0, err // return -1?
	}
	templateSheetIndex, err := FindTemplateSheet(excelFile, templateSheetName)
	if err != nil {
		return 0, err
	}
	err = CopySheet(excelFile, templateSheetIndex, newSheetIndex)
	if err != nil {
		return 0, err
	}
	return newSheetIndex, nil
}


func MakeSheet(excelFile *excelize.File, sheetName string) (int, error) {
	newSheetIndex, err := excelFile.NewSheet(sheetName)
	if err != nil {
		log.Debug("Error creating new sheet")
		return 0, err
	}
	log.Debug(fmt.Sprintf("Sheet created successfully: '%s' at index '%d'", sheetName, newSheetIndex))
	return newSheetIndex, err
}

func CopySheet(excelFile *excelize.File, sourceSheetIndex int, targetSheetIndex int) error {
	if err := excelFile.CopySheet(sourceSheetIndex, targetSheetIndex); err != nil {
		return err
	}
	log.Debug("Sheet copied successfully")
	log.Debug(excelFile.GetSheetName(targetSheetIndex))
	return nil
}