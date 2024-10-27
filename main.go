/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	logger "github.com/williamvannuffelen/go_zaplogger_iso8601"
	_ "github.com/williamvannuffelen/tse/cmd"
	_ "github.com/williamvannuffelen/tse/workitem"
	"github.com/williamvannuffelen/tse/excel"
	"github.com/williamvannuffelen/tse/config"
)

// UNUSED allows unused variables to be included in Go programs
func UNUSED(x ...interface{}) {}

func main() {
	config := config.InitConfig()
	log, err := logger.InitLogger("log.txt", config.General.LogLevel)
	if err != nil {
		log.Warn(err)
	}
	//cmd.SetLogger(log)
	//cmd.Execute()
	//workitem.SetLogger(log)
	//workitem.Run()
	excel.SetLogger(log)
	excelFile, err := excel.SetTargetSheet("ebase.xlsx", "", config.File.TemplateSheetName)
	if err != nil {
		log.Warn(err)
	} else {
		log.Info("Selected sheet ", excelFile.Path)
	}
}
