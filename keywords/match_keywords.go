package keywords

import (
	_ "encoding/json"
	"fmt"
	logger "github.com/williamvannuffelen/go_zaplogger_iso8601"
	_ "github.com/williamvannuffelen/tse/config"
	help "github.com/williamvannuffelen/tse/helpers"
)

var log logger.Logger

func SetLogger(l logger.Logger) {
	log = l
}

// TODO: move to keywords package - gets sourced by addKeyword as well, makes no sense here
func MatchAndExtractKeywords(filePath string, keyword string, callerType string) (map[string]string, error) {
	log.Debug("Matching keywords for keyword: ", keyword)
	kw, err := MatchKeywords(filePath, keyword)
	if err != nil {
		log.Debug("err in matchkeywords")
		return nil, fmt.Errorf("%s %w", help.NewErrorStackTraceString(fmt.Sprintf("failed to get info for provided keyword '%s'", keyword)), err)
	}
	if kw == (Keyword{}) {
		if callerType == "addTimeSheetEntry" {
			log.Debug("No match found for keyword: ", keyword)
			return nil, fmt.Errorf("%s %w", help.NewErrorStackTraceString(fmt.Sprintf("no match found for keyword '%s'", keyword)), fmt.Errorf("keyword not found"))
		}
		if callerType == "addKeyword" {
			log.Debug("No match found for keyword: ", keyword)
			return nil, nil
		}
	}
	keywordValues := map[string]string{
		"description": kw.Description,
		"jira-ref":    kw.JiraRef,
		"project":     kw.Project,
		"app-ref":     kw.AppRef,
	}
	log.Debug(fmt.Sprintf("Got values from keyword: Description: '%s', JiraRef: '%s', Project: '%s', AppRef: '%s'",
		keywordValues["description"],
		keywordValues["jira-ref"],
		keywordValues["project"],
		keywordValues["app-ref"]))
	return keywordValues, nil
}

func MatchParamToKeywords(keywords map[string]Keyword, param string) (Keyword, error) {
	if keyword, exists := keywords[param]; exists {
		return keyword, nil
	}
	return Keyword{}, nil
}

func MatchKeywords(fileName string, param string) (Keyword, error) {
	errorMessage := "failed to match keywords"
	keywords, err := UnmarshalToKeywords(fileName)
	if err != nil {
		return Keyword{}, fmt.Errorf("%s %w", help.NewErrorStackTraceString(errorMessage), err)
	}
	matchedKeyword, err := MatchParamToKeywords(keywords, param)
	if err != nil {
		return Keyword{}, fmt.Errorf("%s %w", help.NewErrorStackTraceString(errorMessage), err)
	}
	if matchedKeyword == (Keyword{}) {
		log.Debug("No match found for keyword: ", param)
		return Keyword{}, nil
	}
	log.Debug("Found match for keyword: ", param)
	return matchedKeyword, nil
}
