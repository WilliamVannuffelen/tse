package config

import (
	"github.com/spf13/viper"
	"log"
	"os"
)

type Config struct {
	File    FileConfig    `mapstructure:"File"`
	Project ProjectConfig `mapstructure:"Project"`
	JiraRef JiraRefConfig `mapstructure:"JiraRef"`
	AppRef  AppRefConfig  `mapstructure:"AppRef"`
}

type FileConfig struct {
	TargetFilePath  string `mapstructure:"targetFilePath"`
	TargetSheetName string `mapstructure:"targetSheetName"`
}

type ProjectConfig struct {
	DefaultProjectName string `mapstructure:"defaultProjectName"`
}

type JiraRefConfig struct {
	DefaultValue    string `mapstructure:"defaultValue"`
	SetDefaultValue bool   `mapstructure:"setDefaultValue"`
}

type AppRefConfig struct {
	DefaultValue    string `mapstructure:"defaultValue"`
	SetDefaultValue bool   `mapstructure:"setDefaultValue"`
}

var cfgFile string

func InitConfig(debugEnabled bool) {
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
		log.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		log.Println("Failed to read config file:", err)
	}

	if debugEnabled {
		log.Println("Raw config data:", viper.AllSettings())
	}

	err := viper.Unmarshal(&config)
	if err != nil {
		log.Println("Failed to unmarshal config:", err)
	}
	if debugEnabled {
		log.Println("Unmarshalled config file.")
	}
	if debugEnabled {
		log.Printf("Config: %+v\n", config)
	}

	if config.File.TargetFilePath == "" {
		log.Println("Warning: TargetFilePath is empty")
	}
	if config.File.TargetSheetName == "" {
		log.Println("Warning: TargetSheetName is empty")
	}
	if config.Project.DefaultProjectName == "" {
		log.Println("Warning: DefaultProjectName is empty")
	}
	if config.JiraRef.DefaultValue == "" {
		log.Println("Warning: JiraRef DefaultValue is empty")
	}
	if config.AppRef.DefaultValue == "" {
		log.Println("Warning: AppRef DefaultValue is empty")
	}
}
