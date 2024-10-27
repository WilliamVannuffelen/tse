package helpers

import (
	"fmt"
	"runtime"
	"strings"
)

func NewErrorStackTraceString(message string) string {
	pc, file, line, _ := runtime.Caller(1)

	callerName := runtime.FuncForPC(pc).Name()
	callerParts := strings.Split(callerName, "/")
	callerName = callerParts[len(callerParts)-1]

	parts := strings.Split(file, "/")
	fileName := strings.Join(parts[len(parts)-2:], "/")

	return fmt.Sprintf("%s:%d - %s - %s \n from", fileName, line, callerName, message)
}
