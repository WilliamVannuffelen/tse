package excel

import (
	"fmt"
	help "github.com/williamvannuffelen/tse/helpers"
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
