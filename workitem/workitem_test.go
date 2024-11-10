package workitem

import (
	logger "github.com/williamvannuffelen/go_zaplogger_iso8601"
	"github.com/williamvannuffelen/tse/config"
	"testing"
	"time"
)

func init() {
	log, err := logger.InitLogger("log.txt", "error")
	if err != nil {
		log.Warn(err)
	}
	SetLogger(log)
}

func TestSetDate(t *testing.T) {
	testCases := []struct {
		name     string
		date     string
		expected string
	}{
		{
			name:     "Empty date string",
			date:     "",
			expected: time.Now().Format("2006-01-02"),
		},
		{
			name:     "Specific date string",
			date:     "2023-09-30",
			expected: "2023-09-30",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			kwi := &KiaraWorkItem{}
			kwi.SetDate(tc.date)

			Equal(t, kwi.Date, tc.expected)
		})
	}
}

func TestSetDay(t *testing.T) {
	testCases := []struct {
		name     string
		date     string
		expected string
	}{
		{
			name:     "valid monday",
			date:     "2024-09-30",
			expected: "Mon",
		},
		{
			name:     "valid wednesday",
			date:     "2024-10-02",
			expected: "Wed",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			kwi := &KiaraWorkItem{}
			kwi.SetDay(tc.date)
			Equal(t, kwi.Day, tc.expected)
		})
	}
}

func TestSetDefaultValue(t *testing.T) {
	testCases := []struct {
		name            string
		valueType       string
		defaultValue    string
		setDefaultValue bool
		value           string
		expected        string
	}{
		{
			name:            "Empty JiraRef value and setDefaultValue is true",
			valueType:       "JiraRef",
			defaultValue:    "OPS-305",
			setDefaultValue: true,
			value:           "",
			expected:        "OPS-305",
		},
		{
			name:            "Specific value",
			valueType:       "JiraRef",
			defaultValue:    "OPS-305",
			setDefaultValue: true,
			value:           "OPS-306",
			expected:        "OPS-306",
		},
		{
			name:            "Empty value and setDefaultValue is false",
			valueType:       "JiraRef",
			defaultValue:    "OPS-305",
			setDefaultValue: false,
			value:           "",
			expected:        "",
		},
		{
			name:            "Empty Project value and setDefaultValue is true",
			valueType:       "Project",
			defaultValue:    "PROJ-123",
			setDefaultValue: true,
			value:           "",
			expected:        "PROJ-123",
		},
		{
			name:            "Specific Project value",
			valueType:       "Project",
			defaultValue:    "PROJ-123",
			setDefaultValue: true,
			value:           "PROJ-124",
			expected:        "PROJ-124",
		},
		{
			name:            "Empty Project value and setDefaultValue is false",
			valueType:       "Project",
			defaultValue:    "PROJ-123",
			setDefaultValue: false,
			value:           "",
			expected:        "",
		},
		{
			name:            "Empty AppRef value and setDefaultValue is true",
			valueType:       "AppRef",
			defaultValue:    "APP-789",
			setDefaultValue: true,
			value:           "",
			expected:        "APP-789",
		},
		{
			name:            "Specific AppRef value",
			valueType:       "AppRef",
			defaultValue:    "APP-789",
			setDefaultValue: true,
			value:           "APP-790",
			expected:        "APP-790",
		},
		{
			name:            "Empty AppRef value and setDefaultValue is false",
			valueType:       "AppRef",
			defaultValue:    "APP-789",
			setDefaultValue: false,
			value:           "",
			expected:        "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			kwi := &KiaraWorkItem{}
			kwi.SetDefaultValue(tc.valueType, tc.defaultValue, tc.setDefaultValue, tc.value)
			switch tc.valueType {
			case "JiraRef":
				Equal(t, kwi.JiraRef, tc.expected)
			case "Project":
				Equal(t, kwi.Project, tc.expected)
			case "AppRef":
				Equal(t, kwi.AppRef, tc.expected)
			}
		})
	}
}

func TestNewKiaraWorkItem(t *testing.T) {
	appConfig := config.Config{
		General: config.GeneralConfig{
			DebugEnabled:          true,
			LogLevel:              "debug",
			SilenceConfigWarnings: false,
		},
		File: config.FileConfig{
			TargetFilePath:    "/path/to/target/file",
			UseTemplateSheet:  true,
			TargetSheetName:   "TargetSheet",
			TemplateSheetName: "TemplateSheet",
		},
		Project: config.ProjectConfig{
			DefaultProjectName: "DefaultProject",
		},
		JiraRef: config.JiraRefConfig{
			DefaultValue:                  "JIRA-123",
			SetDefaultValue:               true,
			SetDefaultValueForNewKeywords: false,
		},
		AppRef: config.AppRefConfig{
			DefaultValue:                  "APP-456",
			SetDefaultValue:               true,
			SetDefaultValueForNewKeywords: false,
		},
		Keywords: config.Keywords{
			DefaultOutputFormat: "json",
		},
		ShowTimeSheetEntry: config.ShowTimeSheetEntry{
			DefaultOutputFormat: "json",
			HideProject:         false,
			HideAppRef:          false,
			HideJiraRef:         false,
		},
	}

	testCases := []struct {
		appConfig       config.Config
		name            string
		date            string
		description     string
		jiraRef         string
		timeSpent       string
		project         string
		appRef          string
		expectedDate    string
		expectedDay     string
		expectedJiraRef string
		expectedProject string
		expectedAppRef  string
	}{
		{
			appConfig:       appConfig,
			name:            "Empty values where possible",
			date:            "",
			description:     "Work on project",
			jiraRef:         "",
			timeSpent:       "",
			project:         "",
			appRef:          "",
			expectedDate:    time.Now().Format("2006-01-02"),
			expectedDay:     time.Now().Format("Mon"),
			expectedJiraRef: "JIRA-123",
			expectedProject: "DefaultProject",
			expectedAppRef:  "APP-456",
		},
		{
			appConfig:       appConfig,
			name:            "Specific values",
			date:            "2024-09-30",
			description:     "Work on project",
			jiraRef:         "OPS-305",
			timeSpent:       "8",
			project:         "PROJ-123",
			appRef:          "APP-789",
			expectedDate:    "2024-09-30",
			expectedDay:     "Mon",
			expectedJiraRef: "OPS-305",
			expectedProject: "PROJ-123",
			expectedAppRef:  "APP-789",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			kwi := NewKiaraWorkItem(appConfig, tc.date, tc.description, tc.jiraRef, tc.timeSpent, tc.project, tc.appRef)

			Equal(t, kwi.Date, tc.expectedDate)
			Equal(t, kwi.Day, tc.expectedDay)
			Equal(t, kwi.JiraRef, tc.expectedJiraRef)
			Equal(t, kwi.Project, tc.expectedProject)
			Equal(t, kwi.AppRef, tc.expectedAppRef)
		})
	}
}
