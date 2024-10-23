package cmd

import (
	"github.com/spf13/cobra"
)

// var log logger.Logger

// func SetLogger(l logger.Logger) {
// 	log = l
// }

// addTimeSheetEntryCmd represents the addTimeSheetEntry command
var addTimeSheetEntryCmd = &cobra.Command{
	Use:           "addTimeSheetEntry",
	Short:         "Add timesheet entry",
	Long:          "Adds a timesheet entry to the timesheet.",
	SilenceUsage:  true,
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("addTimeSheetEntry called")
	},
}

func init() {
	rootCmd.AddCommand(addTimeSheetEntryCmd)
	addTimeSheetEntryCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	addTimeSheetEntryCmd.Flags().BoolP("help", "h", false, "Help message for help")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addTimeSheetEntryCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addTimeSheetEntryCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
