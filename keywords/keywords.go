package keywords

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	_ "github.com/williamvannuffelen/tse/config"
	help "github.com/williamvannuffelen/tse/helpers"
)

type Keyword interface{}

type BasicKeyword struct {
	JiraRef string `json:"jiraRef"`
	Project string `json:"project"`
}

type FullKeyword struct {
	JiraRef     string `json:"jiraRef"`
	Project     string `json:"project"`
	Description string `json:"description"`
}

func OpenKeywordsFile(fileName string) (*os.File, error) {
	log.Debug(fmt.Sprintf("Opening file '%s'", fileName))
	file, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("%s %w", help.NewErrorStackTraceString("failed to open keywords_exact json"), err)
	}
	return file, nil
}

func ReadFileBytes(file *os.File) ([]byte, error) {
	log.Debug("Reading file as bytes.")
	byteValue, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("%s %w", help.NewErrorStackTraceString("failed to read file bytes"), err)
	}
	return byteValue, nil
}

func UnmarshalJson(byteValue []byte) (map[string]json.RawMessage, error) {
	log.Debug("Unmarshalling json.")
	var keywords map[string]json.RawMessage
	err := json.Unmarshal(byteValue, &keywords)
	if err != nil {
		return nil, fmt.Errorf("%s %w", help.NewErrorStackTraceString("failed to unmarshal bytearray to json"), err)
	}
	return keywords, nil
}

func UnmarshalToKeywords(fileName string) (map[string]json.RawMessage, error) {
	errorMessage := "failed to unmarshal keywords"
	file, err := OpenKeywordsFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("%s %w", help.NewErrorStackTraceString(errorMessage), err)
	}
	defer file.Close()

	byteValue, err := ReadFileBytes(file)
	if err != nil {
		return nil, fmt.Errorf("%s %w", help.NewErrorStackTraceString(errorMessage), err)
	}
	keywords, err := UnmarshalJson(byteValue)
	if err != nil {
		return nil, fmt.Errorf("%s %w", help.NewErrorStackTraceString(errorMessage), err)
	}
	return keywords, nil
}
