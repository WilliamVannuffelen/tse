package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/williamvannuffelen/tse/keywords"
	//help "github.com/williamvannuffelen/tse/helpers"
)

// process flags - see if valid combination + valid values
// if valid, append to correct json file

var addKeywordCmd = &cobra.Command{
	Use:           "addKeyword",
	Short:         "Add keyword",
	Long:          "Adds a keywords to the keyword storage.",
	SilenceUsage:  true,
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {
		//TODO: keywords filepath should be configurable
		var keywordsFilePath = "./keywords/keywords.json"
		flags := []string{"description", "jira-ref", "project", "app-ref", "keyword"}
		values := make(map[string]string)
		for _, flag := range flags {
			value, _ := cmd.Flags().GetString(flag)
			values[flag] = value
			log.Debug(fmt.Sprintf("%s: %s", flag, value))
		}

		err := keywords.ValidateFlags(values)
		if err != nil {
			log.Error(err)
			return
		}
		keywords.SetDefaultValues(values)
		log.Debug("Processed values: ", values)

		keywordValues, err := keywords.MatchAndExtractKeywords(keywordsFilePath, values["keyword"], "addKeyword")
		if err != nil {
			log.Error(err)
			return
		}
		if keywordValues == nil {
			log.Debug("No keyword found for: ", values["keyword"])
			keywordsMap, err := keywords.UnmarshalToKeywords(keywordsFilePath)
			if err != nil {
				log.Error(err)
				return
			}
			updatedKeywords, err := keywords.AddNewKeyword(values, keywordsMap)
			if err != nil {
				log.Error(err)
				return
			}
			err = keywords.WriteKeywordsToFile(keywordsFilePath, updatedKeywords)
			if err != nil {
				log.Error(err)
				return
			}
		}
		// invoke function to updae keyword values

		log.Debug("Keyword values: ", keywordValues)
	},
}

func init() {
	rootCmd.AddCommand(addKeywordCmd)
	addKeywordCmd.Flags().BoolP("help", "h", false, "Display this help message")
	addKeywordCmd.Flags().StringP("jira-ref", "j", "", "Jira reference of the timesheet entry. Will default to the value set in config.yaml if setting default is not disabled.")
	addKeywordCmd.Flags().StringP("project", "p", "", "Project of the timesheet entry. Will default to the value set in config.yaml")
	addKeywordCmd.Flags().StringP("description", "d", "", "Description of the timesheet entry.")
	addKeywordCmd.Flags().StringP("app-ref", "a", "", "App reference of the timesheet entry. Will default to the value set in config.yaml if setting default is not disabled.")
	addKeywordCmd.Flags().StringP("keyword", "k", "", "Keyword of the timesheet entry. Used to source full description, project, jira-ref and app-ref for known tasks.")
}
