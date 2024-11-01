package keywords

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/williamvannuffelen/tse/config"
	help "github.com/williamvannuffelen/tse/helpers"
	"os"
)

var appConfig = config.InitConfig()

func ValidateFlags(values map[string]string) error {
	if values["keyword"] == "" {
		return fmt.Errorf("%s %w", help.NewErrorStackTraceString("no keyword provided. Provide one using -k"), errors.New("keyword value is nil"))
	}
	return nil
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

func AddNewKeyword(values map[string]string, keywordsMap map[string]Keyword) ([]byte, error) {
	keywordKey := values["keyword"]
	keyword := Keyword{
		Description: values["description"],
		JiraRef:     values["jira-ref"],
		Project:     values["project"],
		AppRef:      values["app-ref"],
	}
	log.Debug(fmt.Sprintf("Adding new keyword: '%s'. Description: '%s' Jiraref: '%s' Project: '%s' AppRef: '%s'",
		keywordKey,
		values["description"],
		values["jira-ref"],
		values["project"],
		values["app-ref"]))

	keywordsMap[keywordKey] = keyword
	updatedKeywords, err := json.MarshalIndent(keywordsMap, "", " ")
	if err != nil {
		return nil, fmt.Errorf("%s %w", help.NewErrorStackTraceString("failed to marshal new keyword to json"), err)
	}

	// uncomment if required; noisy and expensive even by debug standards
	// var updatedKwMap map[string]Keyword
	// err = json.Unmarshal(updatedKeywords, &updatedKwMap)
	// if err != nil {
	// 	return fmt.Errorf("%s %w", help.NewErrorStackTraceString("failed to unmarshal json to keyword map"), err)
	// }
	// log.Debug("Unmarshalled updated keywords: ", updatedKwMap)
	return updatedKeywords, nil
}

func WriteKeywordsToFile(keywordsFilePath string, updatedKeywords []byte) error {
	file, err := os.OpenFile(keywordsFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("%s %w", help.NewErrorStackTraceString("failed to open keywords.json"), err)
	}
	defer file.Close()

	_, err = file.Write(updatedKeywords)
	if err != nil {
		return fmt.Errorf("%s %w", help.NewErrorStackTraceString("failed to write to keywords.json"), err)
	}
	return nil
}
