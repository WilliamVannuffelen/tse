package config

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os"
)

type Config struct {
	General  GeneralConfig `mapstructure:"General"`
	File     FileConfig    `mapstructure:"File"`
	Project  ProjectConfig `mapstructure:"Project"`
	JiraRef  JiraRefConfig `mapstructure:"JiraRef"`
	AppRef   AppRefConfig  `mapstructure:"AppRef"`
	Keywords Keywords      `mapstructure:"Keywords"`
}

type GeneralConfig struct {
	DebugEnabled          bool   `mapstructure:"debugEnabled"`
	LogLevel              string `mapstructure:"logLevel"`
	SilenceConfigWarnings bool   `mapstructure:"silenceConfigWarnings"`
}

type FileConfig struct {
	TargetFilePath    string `mapstructure:"targetFilePath"`
	TargetSheetName   string `mapstructure:"targetSheetName"`
	TemplateSheetName string `mapstructure:"templateSheetName"`
}

type ProjectConfig struct {
	DefaultProjectName string `mapstructure:"defaultProjectName"`
}

type JiraRefConfig struct {
	DefaultValue                  string `mapstructure:"defaultValue"`
	SetDefaultValue               bool   `mapstructure:"setDefaultValue"`
	SetDefaultValueForNewKeywords bool   `mapstructure:"setDefaultValueForNewKeywords"`
}

type AppRefConfig struct {
	DefaultValue                  string `mapstructure:"defaultValue"`
	SetDefaultValue               bool   `mapstructure:"setDefaultValue"`
	SetDefaultValueForNewKeywords bool   `mapstructure:"setDefaultValueForNewKeywords"`
}

type Keywords struct {
	DefaultOutputFormat string `mapstructure:"defaultOutputFormat"`
}

var cfgFile string

func InitConfig() Config {
	var config Config

	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath("config")
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	if err := viper.ReadInConfig(); err == nil {
	} else {
		log.Println("Failed to read config file:", err)
	}

	err := viper.Unmarshal(&config)
	if err != nil {
		log.Println("Failed to unmarshal config:", err)
	}
	debugEnabled := config.General.DebugEnabled
	if debugEnabled {
		log.Println("Running in debug mode. Will print verbose messages. Set debugMode to false in config file to disable.")
		log.Println("Unmarshalled config file.")
		log.Printf("Config: %+v\n", config)
	}

	if !config.General.SilenceConfigWarnings {
		warnings := []struct {
			value      string
			field      string
			defaultMsg string
		}{
			{config.General.LogLevel, "LogLevel", "Will use 'info'."},
			{config.File.TargetFilePath, "TargetFilePath", ""},
			{config.File.TargetSheetName, "TargetSheetName", ""},
			{config.File.TemplateSheetName, "TemplateSheetName", "Will use sheet at index 0."},
			{config.Project.DefaultProjectName, "DefaultProjectName", ""},
			{config.JiraRef.DefaultValue, "JiraRef DefaultValue", ""},
			{config.AppRef.DefaultValue, "AppRef DefaultValue", ""},
			{config.Keywords.DefaultOutputFormat, "Keywords DefaultOutputFormat", "Will use 'pp'."},
		}

		for _, w := range warnings {
			if w.value == "" {
				message := fmt.Sprintf("Warning: %s is empty. %s Set 'silenceConfigWarnings: true' in config file to suppress this warning.", w.field, w.defaultMsg)
				log.Println(message)
			}
		}
	}
	return config
}
