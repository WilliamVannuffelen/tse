package workitem

import (
	"fmt"
	_ "github.com/spf13/viper"
	logger "github.com/williamvannuffelen/go_zaplogger_iso8601"
	"github.com/williamvannuffelen/tse/config"
	"time"
)

var log logger.Logger

func SetLogger(l logger.Logger) {
	log = l
}

func GetCurrentDayDateString() string {
	currentTime := time.Now()
	return currentTime.Format("2006-01-02")
}

func (kwi *KiaraWorkItem) SetDate(date string) {
	if date == "" {
		date = GetCurrentDayDateString()
		log.Debug(fmt.Sprintf("No date provided, using current date. '%s'", date))
	}
	kwi.Date = date
}

func (kwi *KiaraWorkItem) SetDefaultValue(valueType string, defaultValue string, setDefaultValue bool, value string) {
	if value == "" {
		if setDefaultValue {
			value = defaultValue
			log.Debug(fmt.Sprintf("No %s provided and setDefaultValue is true, using default: '%s'", valueType, defaultValue))
		} else {
			log.Debug(fmt.Sprintf("No %s provided and setDefaultValue is false. Keeping empty string.", valueType))
		}
	} else {
		log.Debug(fmt.Sprintf("Provided %s: '%s'", valueType, value))
	}
	switch valueType {
	case "JiraRef":
		kwi.JiraRef = value
	case "Project":
		kwi.Project = value
	case "AppRef":
		kwi.AppRef = value
	}
}

type KiaraWorkItemGenerator interface {
	SetDate(date string)
	SetDefaultValue(valueType string, defaultValue string, setDefaultValue bool, value string)
}

type KiaraWorkItem struct {
	Date        string
	Description string
	JiraRef     string
	TimeSpent   string
	Project     string
	AppRef      string
}

func NewKiaraWorkItem(
	appConfig config.Config,
	date string,
	description string,
	jiraRef string,
	project string,
	appRef string,
	timeSpent string,
) *KiaraWorkItem {
	kwi := &KiaraWorkItem{
		Description: description,
		JiraRef:     jiraRef,
		TimeSpent:   timeSpent,
		Project:     project,
		AppRef:      appRef,
	}
	//log.Debug("KiaraWorkItem created.")
	// hardcoded true for setDefaultValue since everyone logically wants a default project
	kwi.SetDate(date)
	kwi.SetDefaultValue("JiraRef", appConfig.JiraRef.DefaultValue, appConfig.JiraRef.SetDefaultValue, jiraRef)
	kwi.SetDefaultValue("Project", appConfig.Project.DefaultProjectName, true, project)
	kwi.SetDefaultValue("AppRef", appConfig.AppRef.DefaultValue, appConfig.AppRef.SetDefaultValue, appRef)
	return kwi
}

func Run() {
}
