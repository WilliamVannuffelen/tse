package cmd

import (
	"github.com/spf13/cobra"
	logger "github.com/williamvannuffelen/go_zaplogger_iso8601"
	"os"
)

var log logger.Logger

func SetLogger(l logger.Logger) {
	log = l
}

var rootCmd = &cobra.Command{
	Use:           "tse",
	Short:         "Time Sheet Entry",
	Long:          "Time Sheet Entry is a CLI tool to manage time sheet entries.",
	SilenceUsage:  true,
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("rootCmd called. This does nothing.")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	log.Info("f")
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.Flags().BoolP("help", "h", false, "Help message for help")
	rootCmd.PersistentFlags().String("debug", "d", "Enable debug logging")
}
