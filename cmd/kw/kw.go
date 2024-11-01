package kw

import (
	"github.com/spf13/cobra"
	"github.com/williamvannuffelen/tse/cmd/kw/list"
	//"github.com/williamvannuffelen/tse/cmd"
	logger "github.com/williamvannuffelen/go_zaplogger_iso8601"
)

var log logger.Logger

func SetLogger(l logger.Logger) {
	log = l
}

var KwCmd = &cobra.Command{
	Use:           "kw",
	Short:         "Keyword commands",
	Long:          "Used for keyword commands such as 'kw list'.",
	SilenceUsage:  true,
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {
		log.Debug("fooss!")
	},
}

func init() {
	KwCmd.AddCommand(kwlist.ListCmd)
}
