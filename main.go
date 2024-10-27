/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	_ "fmt"
	logger "github.com/williamvannuffelen/go_zaplogger_iso8601"
	"github.com/williamvannuffelen/tse/cmd"
	"github.com/williamvannuffelen/tse/config"
	"github.com/williamvannuffelen/tse/excel"
	"github.com/williamvannuffelen/tse/workitem"
	_ "runtime"
)

// UNUSED allows unused variables to be included in Go programs
func UNUSED(x ...interface{}) {}

// func logUnwrappedError (log logger.Logger, err error) {
// 	log.Info(err)
// }

// func logErrorChain(log logger.Logger, err error) {
// 	for err != nil {
// 		logUnwrappedError(log, err)
// 		err = errors.Unwrap(err)
// 	}
// }

func main() {
	appConfig := config.InitConfig()
	log, err := logger.InitLogger("log.txt", appConfig.General.LogLevel)
	if err != nil {
		log.Warn(err)
	}
	cmd.SetLogger(log)
	workitem.SetLogger(log)
	excel.SetLogger(log)

	cmd.Execute(appConfig)

	//workitem.Run()

	// excel.SetLogger(log)
	// excelFile, err := excel.SetTargetSheet("ebase.xlsx", "", config.File.TemplateSheetName)
	// if err != nil {
	// 	//log.Warn(err)
	// 	log.Warn(fmt.Sprintf("Error: %+v", err))
	// } else {
	// 	log.Info("Selected sheet ", excelFile.Path)
	// }
}
