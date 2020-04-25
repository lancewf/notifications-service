package commands

import (
	"github.com/lancewf/notifications-service/pkg"
	"github.com/lancewf/notifications-service/pkg/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Config defines the available configuration options
type Config struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Launches Notifications services",
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("Starting Notification Service")

		conf, err := configFromViper()
		if err != nil {
			log.WithFields(log.Fields{
				"error": err.Error(),
			}).Fatal("Failed to load config")
		}

		pkg.New(conf).Start()
	},
}

func init() {
	RootCmd.AddCommand(serveCmd)
}

func configFromViper() (*config.NotificationsConfig, error) {
	cfg := &config.NotificationsConfig{}
	if err := viper.Unmarshal(cfg); err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Fatal("Failed to marshal config options to server config")
	}

	return cfg, nil
}
