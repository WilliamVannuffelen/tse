package keywords

import (
	"fmt"
	"github.com/williamvannuffelen/tse/config"
)

var appConfig = config.InitConfig()

func GetKeywordType(values map[string]string) string {
	log.Debug("Getting keyword type")
	if values["keyword"] != "" {
		log.Debug("Keyword provided: ", values["keyword"])
		return "keyword"
	} else if values["basic-keyword"] != "" {
		log.Debug("Basic-keyword provided: ", values["basic-keyword"])
		return "basic-keyword"
	}
	return ""
}

func ValidateDescriptionSetForKeywordType(keywordType string, values map[string]string) error {
	if keywordType == "keyword" && values["description"] == "" {
		return fmt.Errorf("no description provided for full keyword. Provide one using -d or use a basic-keyword with -K")
	}
	if keywordType == "basic-keyword" && values["description"] != "" {
		log.Warn("Description provided but using a basic-keyword. Description will be ignored.")
	}
	return nil
}

func ValidateFlags(values map[string]string) error {
	keywordType := GetKeywordType(values)
	if keywordType == "" {
		return fmt.Errorf("no keyword provided. Provide one using -k or -K")
	}
	err := ValidateDescriptionSetForKeywordType(keywordType, values)
	return err
}

func SetDefaultValues(values map[string]string) {
	values["jira-ref"] = SetDefaultValue("JiraRef", appConfig.JiraRef.DefaultValue, appConfig.JiraRef.SetDefaultValueForNewKeywords, values["jira-ref"])
	values["project"] = SetDefaultValue("Project", appConfig.Project.DefaultProjectName, true, values["project"])
	values["app-ref"] = SetDefaultValue("AppRef", appConfig.AppRef.DefaultValue, appConfig.AppRef.SetDefaultValueForNewKeywords, values["app-ref"])
}

func SetDefaultValue(valueType string, defaultValue string, setDefaultValue bool, value string) string {
	if value == "" {
		log.Debug(fmt.Sprintf("value %s is empty", valueType))
		if setDefaultValue {
			value = defaultValue
			log.Debug(fmt.Sprintf("No %s provided and setDefaultValue is true, using default: '%s'", valueType, defaultValue))
		} else {
			log.Debug(fmt.Sprintf("No %s provided and setDefaultValue is false. Keeping empty string.", valueType))
		}
	} else {
		log.Debug(fmt.Sprintf("Provided %s: '%s'", valueType, value))
	}
	return value
}
