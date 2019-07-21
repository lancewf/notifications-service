package pkg

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/lancewf/notifications-service/pkg/config"
	"github.com/lancewf/notifications-service/pkg/inspec"
	"github.com/lancewf/notifications-service/pkg/run"
	log "github.com/sirupsen/logrus"
)

type NotificationReport interface {
	HasNotificationToSend() bool
	WebHookMessage() string
	IFTTTWebHookMessage() string
	SlackWebhookMessage() string
}

// Server
type Server struct {
	config *config.NotificationsConfig
}

func New(config *config.NotificationsConfig) Server {
	log.Infof("Created server with %v", config)
	return Server{
		config: config,
	}
}

func (server Server) Start() {
	mutexCCR := sync.Mutex{}
	mutexInspec := sync.Mutex{}

	http.HandleFunc("/ccr_runs", func(w http.ResponseWriter, r *http.Request) {
		mutexCCR.Lock()
		defer mutexCCR.Unlock()

		if r.Method == "POST" {
			log.Info("run POST")
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				log.WithError(err).Warn("Could not read body")
				return
			}

			run := run.ParseRun(body)

			server.SendNotification(run)

			w.Header().Set("Content-Type", "application/javascript")
			w.Write([]byte("{}"))
		} else if r.Method == "GET" {

			log.Info("run GET")

			w.Header().Set("Content-Type", "application/javascript")
			w.Write([]byte("{}"))
		} else {
			log.Info("Unhandled")
		}
	})

	http.HandleFunc("/inspec_reports", func(w http.ResponseWriter, r *http.Request) {
		mutexInspec.Lock()
		defer mutexInspec.Unlock()

		if r.Method == "POST" {
			log.Info("inspec POST")
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				log.WithError(err).Warn("Could not read body")
				return
			}

			report := inspec.ParseReport(body)

			server.SendNotification(report)

			w.Header().Set("Content-Type", "application/javascript")
			w.Write([]byte("{}"))
		} else if r.Method == "GET" {

			log.Info("inspec GET")

			w.Header().Set("Content-Type", "application/javascript")
			w.Write([]byte("{}"))
		} else {
			log.Info("Unhandled")
		}

	})
	http.ListenAndServe(fmt.Sprintf(":%d", server.config.Service.Port), nil)
}

func (server Server) SendNotification(report NotificationReport) {
	if report.HasNotificationToSend() {

		if server.config.Webhook.URL != "" {
			log.Infof("Send webhook alert")
			_, err := http.Post(server.config.Webhook.URL,
				"application/json", bytes.NewBuffer([]byte(report.WebHookMessage())))
			if err != nil {
				log.Error("Failed to send report")
			}
		}

		if server.config.IFTTTWebhook.URL != "" {
			log.Infof("Send IFTTT webhook alert")
			_, err := http.Post(server.config.IFTTTWebhook.URL,
				"application/json", bytes.NewBuffer([]byte(report.IFTTTWebHookMessage())))
			if err != nil {
				log.Error("Failed to send report")
			}
		}

		if server.config.SlackWebhook.URL != "" {
			log.Infof("Send Slack webhook alert")
			_, err := http.Post(server.config.SlackWebhook.URL,
				"application/json", bytes.NewBuffer([]byte(report.SlackWebhookMessage())))
			if err != nil {
				log.Error("Failed to send report")
			}
		}
	}
}
