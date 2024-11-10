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

func (kwi *KiaraWorkItem) SetDate(date string) {
	if date == "" {
		date = time.Now().Format("2006-01-02")
		log.Debug(fmt.Sprintf("No date provided, using current date. '%s'", date))
	}
	kwi.Date = date
}

func (kwi *KiaraWorkItem) SetDay(date string) {
	parsedDate, _ := time.Parse("2006-01-02", date)
	dayofWeek := parsedDate.Weekday().String()[:3]
	kwi.Day = dayofWeek
}

func (kwi *KiaraWorkItem) SetDefaultValue(valueType string, defaultValue string, setDefaultValue bool, value string) {
	if value == "" {
		log.Debug("value is empty")
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
	SetDay(date string)
}

type KiaraWorkItem struct {
	Day         string
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
	timeSpent string,
	project string,
	appRef string,
) *KiaraWorkItem {
	kwi := &KiaraWorkItem{
		Description: description,
		JiraRef:     jiraRef,
		TimeSpent:   timeSpent,
		Project:     project,
		AppRef:      appRef,
	}
	// hardcoded true for setDefaultValue since everyone logically wants a default project
	kwi.SetDate(date)
	kwi.SetDay(kwi.Date)
	kwi.SetDefaultValue("JiraRef", appConfig.JiraRef.DefaultValue, appConfig.JiraRef.SetDefaultValue, jiraRef)
	kwi.SetDefaultValue("Project", appConfig.Project.DefaultProjectName, true, project)
	kwi.SetDefaultValue("AppRef", appConfig.AppRef.DefaultValue, appConfig.AppRef.SetDefaultValue, appRef)
	return kwi
}

func Run() {
}
