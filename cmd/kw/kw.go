package kw

import (
	"github.com/spf13/cobra"
	"github.com/williamvannuffelen/tse/cmd/kw/add"
	"github.com/williamvannuffelen/tse/cmd/kw/list"
	//"github.com/williamvannuffelen/tse/cmd"
	logger "github.com/williamvannuffelen/go_zaplogger_iso8601"
)

var log logger.Logger

func SetLogger(l logger.Logger) {
	log = l
}

var KwCmd = &cobra.Command{
	Use:           "keyword",
	Aliases:       []string{"kw"},
	Short:         "Alias: kw - keyword commands: list, add, show",
	Long:          "Used for keyword commands such as 'kw list'.",
	SilenceUsage:  true,
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("kw command called without subcommand")
	},
}

func init() {
	KwCmd.AddCommand(kwlist.ListCmd)
	KwCmd.AddCommand(kwadd.AddKeywordCmd)
}
