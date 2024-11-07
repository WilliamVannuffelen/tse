package excel

import (
	"fmt"

	logger "github.com/williamvannuffelen/go_zaplogger_iso8601"
	"github.com/williamvannuffelen/tse/config"
	help "github.com/williamvannuffelen/tse/helpers"
	"github.com/xuri/excelize/v2"
	"time"
)

var log logger.Logger

func SetLogger(l logger.Logger) {
	log = l
}

var appConfig = config.InitConfig()

// rename -> does more than setting target sheet
func MakeOrSelectTargetSheet(fileName string, sheetName string, templateSheetName string) (*excelize.File, error) {
	excelFile, err := OpenExcelFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("%s %w", help.NewErrorStackTraceString("failed to set active sheet"), err)
	}

	if sheetName == "" {
		log.Debug("No sheet name provided, using current week's sheet name.")
		sheetName = help.GetCurrentWeekDate(time.Now)
	}

	sheetIndex, err := FindSheetIndex(excelFile, sheetName)
	if err != nil {
		return nil, fmt.Errorf("%s %w", help.NewErrorStackTraceString("failed to set active sheet"), err)
	}
	if sheetIndex == -1 {
		log.Info(fmt.Sprintf("Sheet '%s' not found in file '%s'", sheetName, fileName))
		var newSheetIndex int
		if appConfig.File.UseTemplateSheet {
			newSheetIndex, err = MakeSheetFromTemplate(excelFile, sheetName, templateSheetName)
			if err != nil {
				return nil, fmt.Errorf("%s %w", help.NewErrorStackTraceString("failed to make target sheet"), err)
			}
		} else {
			newSheetIndex, err = MakeSheetFromScratch(excelFile, sheetName)
			if err != nil {
				return nil, fmt.Errorf("%s %w", help.NewErrorStackTraceString("failed to make target sheet"), err)
			}
		}
		SelectSheet(excelFile, sheetName, newSheetIndex)
	} else {
		SelectSheet(excelFile, sheetName, sheetIndex)
	}
	return excelFile, nil
}
