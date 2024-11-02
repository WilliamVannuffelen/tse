package kwlist

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	logger "github.com/williamvannuffelen/go_zaplogger_iso8601"
	"github.com/williamvannuffelen/tse/config"
	help "github.com/williamvannuffelen/tse/helpers"
	"github.com/williamvannuffelen/tse/keywords"
	"sort"
)

var log logger.Logger

func SetLogger(l logger.Logger) {
	log = l
}

var appConfig = config.InitConfig()

var ListCmd = &cobra.Command{
	Use:           "list",
	Short:         "List keywords in the keyword storage",
	Long:          "List keywords in the keyword storage",
	SilenceUsage:  true,
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {
		values := getFlagValues(cmd)
		log.Debug(fmt.Sprintf("Values: %v", values))

		setDefaultOutputFormat(values, appConfig)

		//TODO: keywords filepath should be configurable
		var keywordsFilePath = "./keywords/keywords.json"
		keywordsMap, err := keywords.UnmarshalToKeywords(keywordsFilePath)
		if err != nil {
			log.Error(err)
			return
		}
		if values["output"] == "pp" {
			printKeywords(keywordsMap, values)
		} else {
			printKeywordsAsJSON(keywordsMap)
			if err != nil {
				log.Error(err)
			}
		}
	},
}

func init() {
	ListCmd.Flags().BoolP("help", "h", false, "Display this help message")
	ListCmd.Flags().BoolP("keys", "k", true, "Display keys.")
	ListCmd.Flags().BoolP("all", "A", false, "Display everything.")
	ListCmd.Flags().BoolP("description", "d", false, "Display description.")
	ListCmd.Flags().BoolP("project", "p", false, "Display project.")
	ListCmd.Flags().BoolP("jira-ref", "j", false, "Display jira reference.")
	ListCmd.Flags().BoolP("app-ref", "a", false, "Display app reference.")
	ListCmd.Flags().StringP("output", "o", "", "Output format. Options: json, pp (pretty print).")
}

func getFlagValues(cmd *cobra.Command) map[string]interface{} {
	boolFlags := []string{"keys", "all", "description", "project", "jira-ref", "app-ref"}
	values := make(map[string]interface{})
	for _, flag := range boolFlags {
		value, _ := cmd.Flags().GetBool(flag)
		values[flag] = value
		log.Debug(fmt.Sprintf("%s: %t", flag, value))
	}

	stringFlags := []string{"output"}
	for _, flag := range stringFlags {
		value, _ := cmd.Flags().GetString(flag)
		values[flag] = value
		log.Debug(fmt.Sprintf("%s: %s", flag, value))
	}

	return values
}

func setDefaultOutputFormat(values map[string]interface{}, appConfig config.Config) {
	if values["output"] == "" {
		values["output"] = appConfig.Keywords.DefaultOutputFormat
	}
}

func printKeywordsAsJSON(keywords map[string]keywords.Keyword) error {
	jsonData, err := json.MarshalIndent(keywords, "", "  ")
	if err != nil {
		return fmt.Errorf("%s %w", help.NewErrorStackTraceString("failed to set active sheet"), err)
	}

	fmt.Println(string(jsonData))
	return nil
}

func printKeywords(keywords map[string]keywords.Keyword, values map[string]interface{}) {
	fmt.Println("Keywords:")
	keys := make([]string, 0, len(keywords))
	for key := range keywords {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		keyword := keywords[key]
		keywordDetails := make(map[string]interface{})

		all, _ := values["all"].(bool)
		description, _ := values["description"].(bool)
		project, _ := values["project"].(bool)
		jiraRef, _ := values["jira-ref"].(bool)
		appRef, _ := values["app-ref"].(bool)

		if all {
			if keyword.Description != "" {
				keywordDetails["description"] = keyword.Description
			}
			if keyword.Project != "" {
				keywordDetails["project"] = keyword.Project
			}
			if keyword.JiraRef != "" {
				keywordDetails["jira-ref"] = keyword.JiraRef
			}
			if keyword.AppRef != "" {
				keywordDetails["app-ref"] = keyword.AppRef
			}
		} else {
			if description && keyword.Description != "" {
				keywordDetails["description"] = keyword.Description
			}
			if project && keyword.Project != "" {
				keywordDetails["project"] = keyword.Project
			}
			if jiraRef && keyword.JiraRef != "" {
				keywordDetails["jira-ref"] = keyword.JiraRef
			}
			if appRef && keyword.AppRef != "" {
				keywordDetails["app-ref"] = keyword.AppRef
			}
		}

		if len(keywordDetails) > 0 {
			color.New(color.FgBlue, color.Bold).Printf("%s: ", key)
			color.New(color.FgHiYellow).Println("{")

			subKeys := make([]string, 0, len(keywordDetails))
			for subKey := range keywordDetails {
				subKeys = append(subKeys, subKey)
			}
			sort.Strings(subKeys)

			for _, subKey := range subKeys {
				if subKey == "description" {
					color.New(color.FgBlue, color.Bold).Printf("  \"%s\": ", subKey)
				} else if subKey == "project" {
					color.New(color.FgCyan).Printf("  \"%s\": ", subKey)
				} else if subKey == "jira-ref" {
					color.New(color.FgMagenta).Printf("  \"%s\": ", subKey)
				} else if subKey == "app-ref" {
					color.New(color.FgYellow).Printf("  \"%s\": ", subKey)
				}
				color.New(color.FgGreen).Printf("\"%v\",\n", keywordDetails[subKey])
			}
			color.New(color.FgHiYellow).Println("}")
		} else {
			color.New(color.FgBlue, color.Bold).Println(key)
		}
	}
}
