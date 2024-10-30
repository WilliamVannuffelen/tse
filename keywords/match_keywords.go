package keywords

import (
	"encoding/json"
	"fmt"
	logger "github.com/williamvannuffelen/go_zaplogger_iso8601"
	_ "github.com/williamvannuffelen/tse/config"
	help "github.com/williamvannuffelen/tse/helpers"
)

var log logger.Logger

func SetLogger(l logger.Logger) {
	log = l
}

func MatchParamToKeywords(keywords map[string]json.RawMessage, param string) (FullKeyword, error) {
	if keyword, exists := keywords[param]; exists {
		var kw FullKeyword
		err := json.Unmarshal(keyword, &kw)
		if err != nil {
			return FullKeyword{}, fmt.Errorf("%s %w", help.NewErrorStackTraceString("failed to unmarshal json to keyword struct"), err)
		}
		log.Debug("Unmarshalled keyword: ", kw)
		return kw, nil
	}
	return FullKeyword{}, nil
}

func MatchKeywords(fileName string, param string) (FullKeyword, error) {
	errorMessage := "failed to match keywords"
	keywords, err := UnmarshalToKeywords(fileName)
	if err != nil {
		return FullKeyword{}, fmt.Errorf("%s %w", help.NewErrorStackTraceString(errorMessage), err)
	}
	log.Debug("Starting MatchParamtoKeywords")
	matchedKeyword, err := MatchParamToKeywords(keywords, param)
	if err != nil {
		return FullKeyword{}, fmt.Errorf("%s %w", help.NewErrorStackTraceString(errorMessage), err)
	}
	if matchedKeyword == (FullKeyword{}) {
		return FullKeyword{}, fmt.Errorf("%s %w", help.NewErrorStackTraceString(errorMessage), fmt.Errorf("keyword not found '%s'", param))
	}
	log.Debug("Matched keyword: ", matchedKeyword)
	return matchedKeyword, nil
}
