package kwadd

import (
	"fmt"

	"github.com/spf13/cobra"
	logger "github.com/williamvannuffelen/go_zaplogger_iso8601"
	"github.com/williamvannuffelen/tse/config"
	"github.com/williamvannuffelen/tse/keywords"
)

var log logger.Logger

func SetLogger(l logger.Logger) {
	log = l
}

var appConfig = config.InitConfig()

var AddKeywordCmd = &cobra.Command{
	Use:           "add",
	Short:         "Add keyword",
	Long:          "Adds a keywords to the keyword storage.",
	SilenceUsage:  true,
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {
		//TODO: keywords filepath should be configurable
		var keywordsFilePath = "./keywords/keywords.json"
		values := getFlagValues(cmd)
		if err := validateAndSetDefaults(values); err != nil {
			log.Error(err)
			return
		}

		setDefaultOutputFormat(values, appConfig)

		kwValues := make(map[string]string)
		for k, v := range values {
			if k != "output" {
				kwValues[k] = v
			}
		}

		keywordValues, err := keywords.MatchAndExtractKeywords(keywordsFilePath, values["keyword"], "addKeyword")
		if err != nil {
			log.Error(err)
			return
		}

		keywordsMap, err := keywords.UnmarshalToKeywords(keywordsFilePath)
		if err != nil {
			log.Error(err)
			return
		}

		if keywordValues == nil {
			log.Debug("No keyword found for: ", values["keyword"])
			if err := addNewKeyword(kwValues, keywordsMap, keywordsFilePath); err != nil {
				log.Error(err)
			}
			if values["output"] == "pp" {
				fmt.Println("Keyword added:")
			}
		} else {
			log.Debug("Existing keyword found for: ", values["keyword"])
			if err := updateExistingKeyword(kwValues, keywordsMap, keywordsFilePath); err != nil {
				log.Error(err)
			}
			if values["output"] == "pp" {
				fmt.Println("Keyword updated:")
			}
		}
		if values["output"] == "pp" {
			keywords.PrettyPrintKeyword(kwValues)
		} else {
			err := keywords.PrintKeywordAsJSON(kwValues)
			if err != nil {
				log.Error(err)
			}
		}
	},
}

func init() {
	AddKeywordCmd.Flags().BoolP("help", "h", false, "Display this help message")
	AddKeywordCmd.Flags().StringP("jira-ref", "j", "", "Jira reference of the timesheet entry. Will default to the value set in config.yaml if setting default is not disabled.")
	AddKeywordCmd.Flags().StringP("project", "p", "", "Project of the timesheet entry. Will default to the value set in config.yaml")
	AddKeywordCmd.Flags().StringP("description", "d", "", "Description of the timesheet entry.")
	AddKeywordCmd.Flags().StringP("app-ref", "a", "", "App reference of the timesheet entry. Will default to the value set in config.yaml if setting default is not disabled.")
	AddKeywordCmd.Flags().StringP("keyword", "k", "", "Keyword of the timesheet entry. Used to source full description, project, jira-ref and app-ref for known tasks.")
	AddKeywordCmd.Flags().StringP("output", "o", "", "Output format. Options: json, pp (pretty print). Default: pp")
}

func getFlagValues(cmd *cobra.Command) map[string]string {
	flags := []string{"description", "jira-ref", "project", "app-ref", "keyword", "output"}
	values := make(map[string]string)
	for _, flag := range flags {
		value, _ := cmd.Flags().GetString(flag)
		values[flag] = value
		log.Debug(fmt.Sprintf("%s: %s", flag, value))
	}
	return values
}

func validateAndSetDefaults(values map[string]string) error {
	if err := keywords.ValidateFlags(values); err != nil {
		return err
	}
	keywords.SetDefaultValues(values)
	log.Debug("Processed values: ", values)
	return nil
}

func setDefaultOutputFormat(values map[string]string, appConfig config.Config) {
	if values["output"] == "" {
		values["output"] = appConfig.Keywords.DefaultOutputFormat
	}
}

func addNewKeyword(values map[string]string, keywordsMap map[string]keywords.Keyword, keywordsFilePath string) error {
	updatedKeywords, err := keywords.AddNewKeyword(values, keywordsMap)
	if err != nil {
		return err
	}
	if err := keywords.WriteKeywordsToFile(keywordsFilePath, updatedKeywords); err != nil {
		return err
	}
	return nil
}

func updateExistingKeyword(values map[string]string, keywordsMap map[string]keywords.Keyword, keywordsFilePath string) error {
	updatedKeywords, err := keywords.UpdateKeyword(values, keywordsMap)
	if err != nil {
		return err
	}
	if err := keywords.WriteKeywordsToFile(keywordsFilePath, updatedKeywords); err != nil {
		return err
	}
	return nil
}
