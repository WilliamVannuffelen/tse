package workitem

import (
	"fmt"
	logger "github.com/williamvannuffelen/go_zaplogger_iso8601"
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
	kwi.SetDate(date)
	kwi.SetDefaultValue("JiraRef", "OPS-305", false, jiraRef)
	kwi.SetDefaultValue("Project", "CS0126444 - Wonen Cloudzone - dedicated operationeel projectteam", true, project)
	kwi.SetDefaultValue("AppRef", "99999", true, appRef)
	return kwi
}

func Run() {
	var date string = ""
	var description string = "Test work item"
	var jiraRef string = ""
	var timeSpent string = "1.0"
	var project string = ""
	var appRef string = "Test app ref"

	kwi := NewKiaraWorkItem(
		date,
		description,
		jiraRef,
		project,
		appRef,
		timeSpent,
	)
	log.Debug(fmt.Sprintf(
		"KiaraWorkItem:\nDate: %s\nDescription: %s\nJiraRef: %s\nTimeSpent: %s\nProject: %s\nAppRef: %s",
		kwi.Date,
		kwi.Description,
		kwi.JiraRef,
		kwi.TimeSpent,
		kwi.Project,
		kwi.AppRef,
	))
}
