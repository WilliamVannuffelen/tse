package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	//logger "github.com/williamvannuffelen/go_zaplogger_iso8601"
	"github.com/spf13/viper"
	_ "github.com/williamvannuffelen/tse/config"
)

var (
	cfgFile string
	rootCmd = &cobra.Command{
		Use:           "tse",
		Short:         "Time Sheet Entry",
		Long:          "Time Sheet Entry is a CLI tool to manage time sheet entries.",
		SilenceUsage:  true,
		SilenceErrors: true,
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.Flags().BoolP("help", "h", false, "Help message for help")
	rootCmd.PersistentFlags().String("debug", "d", "Enable debug logging")
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./config/config.yaml)")

	rootCmd.AddCommand(addTimeSheetEntryCmd)
	fmt.Println("no errors yet")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {

		// Search config in home directory with name ".cobra" (without extension).
		viper.AddConfigPath("./config/config.yaml")
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
