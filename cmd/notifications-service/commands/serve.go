package commands

import (
	"github.com/lancewf/notifications-service/pkg"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
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
		server := pkg.New(8080)

		server.Start()
	},
}

func init() {
	RootCmd.AddCommand(serveCmd)
}