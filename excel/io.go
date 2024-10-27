package excel

import (
	"fmt"
  "github.com/xuri/excelize/v2"
)

func OpenExcelFile(fileName string) (*excelize.File, error) {
	log.Debug(fmt.Sprintf("Opening file '%s'", fileName))
	excelFile, err := excelize.OpenFile(fileName)
	if err != nil {
		return nil, err
	}
	return excelFile, nil
}

func SaveExcelFile (excelFile *excelize.File) () {
	newRow := []interface{}{"Mo", "2024-10-05", "new task", "OPS-999", 3, "New Project", "AppRef"}
	if err := AppendRow(excelFile, "2024-10-21", newRow); err != nil {
			fmt.Println("Error appending row:", err)
			return
	}
	// Save the file
	if err := excelFile.SaveAs("ebase.xlsx"); err != nil {
			fmt.Println("Error saving file:", err)
			return
	} else {
		log.Debug("Saved file successfully!", err)
	}
}