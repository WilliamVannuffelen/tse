package keywords

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	_ "github.com/williamvannuffelen/tse/config"
	help "github.com/williamvannuffelen/tse/helpers"
)

type Keyword struct {
	JiraRef     string `json:"jiraRef,omitempty"`
	Project     string `json:"project,omitempty"`
	Description string `json:"description,omitempty"`
	AppRef      string `json:"appRef,omitempty"`
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

// func UnmarshalJson(byteValue []byte) (map[string]json.RawMessage, error) {
func UnmarshalJson(byteValue []byte) (map[string]Keyword, error) {
	log.Debug("Unmarshalling json.")
	keywords := make(map[string]Keyword) //make(map[string]json.RawMessage)
	err := json.Unmarshal(byteValue, &keywords)
	if len(keywords) == 0 {
		return nil, fmt.Errorf("%s %w", help.NewErrorStackTraceString("failed to unmarshal bytearray to json"), err)
	}
	return keywords, nil
}

// func UnmarshalToKeywords(fileName string) (map[string]json.RawMessage, error) {
func UnmarshalToKeywords(fileName string) (map[string]Keyword, error) {
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
	log.Debug("Done unmarshalling keywords")
	return keywords, nil
}
