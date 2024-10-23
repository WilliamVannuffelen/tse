/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	logger "github.com/williamvannuffelen/go_zaplogger_iso8601"
	"github.com/williamvannuffelen/tse/cmd"
	"github.com/williamvannuffelen/tse/workitem"
)

// UNUSED allows unused variables to be included in Go programs
func UNUSED(x ...interface{}) {}

func main() {
	log, err := logger.InitLogger("log.txt", "debug")
	if err != nil {
		log.Warn(err)
	}
	log.Info("Foo!")
	cmd.SetLogger(log)
	cmd.Execute()
	workitem.SetLogger(log)
	workitem.Run()
}
