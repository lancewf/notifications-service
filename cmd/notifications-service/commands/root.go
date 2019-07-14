package commands

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// RootCmd is the command runner.
var RootCmd = &cobra.Command{
	Use:   "notifications-service",
	Short: "Notifications Service",
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		log.Error(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// global config
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.notifications-service.toml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetConfigName(".notifications-service") // name of config file (without extension)
	viper.AddConfigPath("$HOME")                  // adding home directory as first search path
	viper.AddConfigPath(".")

	// override default config file if config is passed in via cli
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	}

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.WithFields(log.Fields{"file": viper.ConfigFileUsed()}).Info("Using config file")
	}

	// Override our config with any matching environment variables
	viper.AutomaticEnv()
}
