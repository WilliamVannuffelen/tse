package add

import (
	logger "github.com/williamvannuffelen/go_zaplogger_iso8601"
	"github.com/williamvannuffelen/tse/config"
	"testing"
)

func init() {
	log, err := logger.InitLogger("log.txt", appConfig.General.LogLevel) // TODO: add log path to config
	if err != nil {
		log.Warn(err)
	}
	appConfig = config.InitConfig()
	SetLogger(log)
}

func Equal[v comparable](t *testing.T, got, expected v) {
	t.Helper()

	if expected != got {
		t.Errorf(`assert.Equal(
		t,
		got:
		%v,
		expected:
		%v
		)`, got, expected)
	}
}

func TestValidateInputValues(t *testing.T) {
	testCases := []struct {
		name          string
		input         map[string]string
		expectedError string
	}{
		{
			name: "Valid input",
			input: map[string]string{
				"description": "Work on project",
				"date":        "2023-10-01",
				"timespent":   "8",
			},
			expectedError: "",
		},
		{
			name: "Missing description",
			input: map[string]string{
				"date":      "2023-10-01",
				"timespent": "8",
			},
			expectedError: "no description provided. Provide one using -d or use a keyword with -k or -K",
		},
		{
			name: "Invalid date format",
			input: map[string]string{
				"description": "Work on project",
				"date":        "01-10-2023",
				"timespent":   "8",
			},
			expectedError: "invalid date format. Please use yyyy-MM-dd. e.g. 2024-09-31",
		},
		{
			name: "Missing time",
			input: map[string]string{
				"description": "Work on project",
				"date":        "2023-10-01",
			},
			expectedError: "no timespent provided. Provide one using -t",
		},
		{
			name: "Invalid time format",
			input: map[string]string{
				"description": "Work on project",
				"date":        "2023-10-01",
				"timespent":   "8h",
			},
			expectedError: "invalid timespent format. Please use a number. e.g. 8",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateInputValues(tc.input)
			if tc.expectedError != "" {
				Equal(t, err.Error(), tc.expectedError)
			} else {
				Equal(t, err, nil)
			}
		})
	}
}
