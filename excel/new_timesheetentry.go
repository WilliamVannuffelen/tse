package excel

import (
	"github.com/xuri/excelize/v2"
	"strconv"
)

func AppendRow(excelFile *excelize.File, sheet string, row []interface{}) error {
	rows, err := excelFile.GetRows(sheet)
	if err != nil {
		return err
	}

	nextRow := len(rows) + 1
	axis := "A" + strconv.Itoa(nextRow)
	return excelFile.SetSheetRow(sheet, axis, &row)
}
