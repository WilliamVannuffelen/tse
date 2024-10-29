/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	// _ "fmt"
	// logger "github.com/williamvannuffelen/go_zaplogger_iso8601"
	"github.com/williamvannuffelen/tse/cmd"
	// "github.com/williamvannuffelen/tse/config"
	// "github.com/williamvannuffelen/tse/excel"
	// "github.com/williamvannuffelen/tse/keywords"
	// "github.com/williamvannuffelen/tse/workitem"
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
	// appConfig := config.InitConfig()
	// log, err := logger.InitLogger("log.txt", appConfig.General.LogLevel)
	// if err != nil {
	// 	log.Warn(err)
	// }
	// cmd.SetLogger(log)
	// workitem.SetLogger(log)
	// excel.SetLogger(log)
	// keywords.SetLogger(log)

	// cmd.Execute(appConfig)

	cmd.Execute()

	// param := "crma"
	// matchedKeyword, err := keywords.MatchKeywords("./keywords/keywords.json", param)
	// if err != nil {
	// 	log.Warn(err)
	// }
	// log.Info(fmt.Sprintf("Matched keyword for %s", param), matchedKeyword)
}
