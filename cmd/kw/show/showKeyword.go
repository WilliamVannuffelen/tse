package kwshow

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

var ShowCmd = &cobra.Command{
	Use:           "show",
	Short:         "Show keyword from the keyword storage",
	Long:          "Show keyword in the keyword storage",
	SilenceUsage:  true,
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {
		values := getFlagValues(cmd)
		log.Debug(fmt.Sprintf("Values: %v", values))

		setDefaultOutputFormat(values, appConfig)

		//TODO: keywords filepath should be configurable
		var keywordsFilePath = "./keywords/keywords.json"

		keywordValues, err := keywords.MatchAndExtractKeywords(keywordsFilePath, values["keyword"], "addKeyword")
		if err != nil {
			log.Error(err)
			return
		}
		if len(keywordValues) == 0 {
			fmt.Println("No keyword found for: ", values["keyword"])
		} else {
			keywordValues["keyword"] = values["keyword"]
			if values["output"] == "pp" {
				keywords.PrettyPrintKeyword(keywordValues)
			} else {
				keywords.PrintKeywordAsJSON(keywordValues)
				if err != nil {
					log.Error(err)
				}
			}
		}
	},
}

func init() {
	ShowCmd.Flags().BoolP("help", "h", false, "Display this help message")
	ShowCmd.Flags().StringP("keyword", "k", "", "Keyword to show.")
	ShowCmd.Flags().StringP("output", "o", "", "Output format. Options: json, pp (pretty print).")
}

func getFlagValues(cmd *cobra.Command) map[string]string {
	values := make(map[string]string)
	stringFlags := []string{"keyword"}
	for _, flag := range stringFlags {
		value, _ := cmd.Flags().GetString(flag)
		values[flag] = value
		log.Debug(fmt.Sprintf("%s: %s", flag, value))
	}
	return values
}

func setDefaultOutputFormat(values map[string]string, appConfig config.Config) {
	if values["output"] == "" {
		values["output"] = appConfig.Keywords.DefaultOutputFormat
	}
}
