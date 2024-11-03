package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	logger "github.com/williamvannuffelen/go_zaplogger_iso8601"
	"github.com/williamvannuffelen/tse/cmd/add"
	"github.com/williamvannuffelen/tse/cmd/kw"
	"github.com/williamvannuffelen/tse/cmd/kw/add"
	"github.com/williamvannuffelen/tse/cmd/kw/list"
	"github.com/williamvannuffelen/tse/cmd/kw/show"
	"github.com/williamvannuffelen/tse/cmd/show"
	"github.com/williamvannuffelen/tse/config"
	"github.com/williamvannuffelen/tse/excel"
	"github.com/williamvannuffelen/tse/keywords"
	"github.com/williamvannuffelen/tse/workitem"
)

var (
	cfgFile   string
	appConfig config.Config
	log       logger.Logger
	err       error
	rootCmd   = &cobra.Command{
		Use:           "tse",
		Short:         "Time Sheet Entry",
		Long:          "Time Sheet Entry is a CLI tool to manage time sheet entries.",
		SilenceUsage:  true,
		SilenceErrors: true,
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.Flags().BoolP("help", "h", true, "Display this help message")
	rootCmd.PersistentFlags().BoolP("debug", "x", false, "Enable debug logging")
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file path (default is ./config/config.yaml)")

	appConfig = config.InitConfig()
	log, err = logger.InitLogger("log.txt", appConfig.General.LogLevel) // TODO: add log path to config
	if err != nil {
		log.Warn(err)
	}
	workitem.SetLogger(log)
	excel.SetLogger(log)
	keywords.SetLogger(log)
	add.SetLogger(log)
	show.SetLogger(log)
	kw.SetLogger(log)
	kwlist.SetLogger(log)
	kwadd.SetLogger(log)
	kwshow.SetLogger(log)
	log.Debug("Logger init done from root.go")

	rootCmd.AddCommand(add.AddTimeSheetEntryCmd)
	rootCmd.AddCommand(show.ShowTimeSheetEntryCmd)
	rootCmd.AddCommand(kw.KwCmd)
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {

		viper.AddConfigPath("./config/config.yaml")
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	if err := viper.ReadInConfig(); err == nil {
		debug, err := rootCmd.Flags().GetBool("debug")
		if err != nil {
			log.Error(err)
		}
		if debug {
			fmt.Println("Using config file:", viper.ConfigFileUsed())
		}
	}
}
