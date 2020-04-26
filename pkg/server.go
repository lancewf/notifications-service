package pkg

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/lancewf/notifications-service/pkg/config"
	"github.com/lancewf/notifications-service/pkg/inspec"
	log "github.com/sirupsen/logrus"
)

type NotificationReport interface {
	HasNotificationToSend() bool
	GetWebHookMessage() string
	GetIFTTTWebHookMessage() string
	GetSlackWebhookMessage() string
}

type Server struct {
	config  *config.NotificationsConfig
	manager *inspec.Manager
}

func New(config *config.NotificationsConfig) Server {
	log.Infof("Created server with %v", config)
	return Server{
		config:  config,
		manager: &inspec.Manager{},
	}
}

func (server Server) Start() {
	log.Infof("server.config.Inspec.MinImpact %f", server.config.Inspec.MinImpact)
	mutexInspec := sync.Mutex{}

	http.HandleFunc("/inspec_reports/data-collector/v0/", func(w http.ResponseWriter, r *http.Request) {
		mutexInspec.Lock()
		defer mutexInspec.Unlock()

		if r.Method == "POST" {
			log.Info("inspec POST")
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				log.WithError(err).Warn("Could not read body")
				return
			}
			server.forwardToAutomate(body)

			report := server.manager.ParseReport(body, server.config.Inspec.MinImpact)
			server.sendNotification(report)

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

func (server Server) forwardToAutomate(body []byte) {
	if server.config.Automate.EnableForwarding {
		log.Infof("Forwarding to Automate")
		request, err := http.NewRequest("POST", server.config.Automate.URL, bytes.NewBuffer(body))
		if err != nil {
			log.Errorf("Failed to forward to Automate %v", err)
			return
		}
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("api-token", server.config.Automate.Token)

		response, err := http.DefaultClient.Do(request)
		if err != nil {
			log.Errorf("Failed to forward to Automate %v", err)
		}
		log.Infof("Automate foward response status %q", response.Status)
	}
}

func (server Server) sendNotification(report NotificationReport) {
	if report.HasNotificationToSend() {

		if server.config.Webhook.URL != "" {
			log.Infof("Send webhook alert")
			_, err := http.Post(server.config.Webhook.URL,
				"application/json", bytes.NewBuffer([]byte(report.GetWebHookMessage())))
			if err != nil {
				log.Error("Failed to send report")
			}
		}

		if server.config.IFTTTWebhook.URL != "" {
			log.Infof("Send IFTTT webhook alert")
			_, err := http.Post(server.config.IFTTTWebhook.URL,
				"application/json", bytes.NewBuffer([]byte(report.GetIFTTTWebHookMessage())))
			if err != nil {
				log.Error("Failed to send report")
			}
		}

		if server.config.SlackWebhook.URL != "" {
			log.Infof("Send Slack webhook alert")
			_, err := http.Post(server.config.SlackWebhook.URL,
				"application/json", bytes.NewBuffer([]byte(report.GetSlackWebhookMessage())))
			if err != nil {
				log.Error("Failed to send report")
			}
		}
	}
}
