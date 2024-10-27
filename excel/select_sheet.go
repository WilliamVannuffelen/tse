package excel

import (
	"fmt"
	"time"
  "github.com/xuri/excelize/v2"
  //logger "github.com/williamvannuffelen/go_zaplogger_iso8601"
)



func FindSheetIndex(excelFile *excelize.File, sheetName string) (int, error) {
	log.Debug(fmt.Sprintf("Finding sheet index for sheet '%s'", sheetName))
	index, err := excelFile.GetSheetIndex(sheetName)
	if err != nil {
			return -1, err
	}

	return index, nil
}

func FindLastSheetIndex(excelFile *excelize.File) (int, error) {
	log.Debug("Finding last sheet index")
	sheetList := excelFile.GetSheetList()
	lastSheetIndex := len(sheetList) - 1
	return lastSheetIndex, nil
}

func SelectSheet(excelFile *excelize.File, sheetName string, sheetIndex int) error {
	excelFile.SetActiveSheet(sheetIndex)
	log.Debug(fmt.Sprintf("Selected sheet '%s' at index '%d'", sheetName, sheetIndex))
	return nil
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


// perhaps not in the right place, figure out where to put this later
func GetCurrentWeekSheetName() string {
	now := time.Now()
	today := time.Now().Weekday()

	offset := int(time.Monday - today)
	if offset > 0 {
		offset = -6
	}

	monday := now.AddDate(0, 0, offset).Format("2006-01-02")

	log.Debug(fmt.Sprintf("This week's Monday: %s", monday))
	return monday
}