package kwlist

import (
	"fmt"
	"github.com/spf13/cobra"
	logger "github.com/williamvannuffelen/go_zaplogger_iso8601"
)

var log logger.Logger

func SetLogger(l logger.Logger) {
	log = l
}

var ListCmd = &cobra.Command{
	Use:           "list",
	Short:         "List keywords in the keyword storage",
	Long:          "List keywords in the keyword storage",
	SilenceUsage:  true,
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {
		log.Debug("foossssssss!!!!!")
		fmt.Println("foossssssss!!!!!")
	},
}

func init() {
	//
}
