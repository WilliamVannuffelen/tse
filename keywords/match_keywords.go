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

func MatchParamToKeywords(keywords map[string]json.RawMessage, param string) (Keyword, error) {
	if keyword, exists := keywords[param]; exists {
		var kw Keyword
		err := json.Unmarshal(keyword, &kw)
		if err != nil {
			return Keyword{}, fmt.Errorf("%s %w", help.NewErrorStackTraceString("failed to unmarshal json to keyword struct"), err)
		}
		return kw, nil
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
		return Keyword{}, fmt.Errorf("%s %w", help.NewErrorStackTraceString(errorMessage), fmt.Errorf("keyword not found '%s'", param))
	}
	log.Debug("Found match for keyword: ", param)
	return matchedKeyword, nil
}
