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
	log.Info("Foo!")
	log.Info(config.General.DebugEnabled)
	//cmd.SetLogger(log)
	//cmd.Execute()
	//workitem.SetLogger(log)
	//workitem.Run()
	excel.SetLogger(log)
	err = excel.MakeSheet("t_upload2.xlsx", "2024-10-21", "")
	if err != nil {
		log.Warn(err)
	} else {
		log.Info("Sheet created")
	}
}
