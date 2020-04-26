package inspec

import (
	"encoding/json"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
)

type Manager struct {
	FailedLastRun              bool
	LastFailedNotificationDate time.Time
}

func (manager *Manager) ParseReport(rawReport []byte, minImpact float32) Report {
	report := ParseReport(rawReport, minImpact)
	report.NotificationToSend = false
	log.Infof("minImpact %f", minImpact)

	failed := len(report.failedProfiles(minImpact)) > 0

	// Start of failure
	// Send failure notification
	if failed && !manager.FailedLastRun {
		log.Info("Start of failure")
		report.NotificationToSend = true
		manager.FailedLastRun = true
		manager.LastFailedNotificationDate = time.Now()
		report.WebHookMessage = "Failed InSpec Report"
		report.IFTTTWebHookMessage = generateFailedIFTTTWebHookMessage(report)
		report.SlackWebhookMessage = generateFailedSlackWebhookMessage(report)
		log.Infof("FailedLastRun %t", manager.FailedLastRun)
		return report
	}

	hoursAgo := time.Now().Add(time.Hour * -12)

	// Send Reminder
	if failed && manager.FailedLastRun && hoursAgo.After(manager.LastFailedNotificationDate) {
		log.Info("Reminder")
		report.NotificationToSend = true
		manager.LastFailedNotificationDate = time.Now()
		report.WebHookMessage = "Reminder of Failed InSpec Report"
		report.IFTTTWebHookMessage = generateReminderOfFailedIFTTTWebHookMessage(report)
		report.SlackWebhookMessage = generateReminderOfFailedSlackWebhookMessage(report)
		return report
	}

	// Send recovery notification
	if manager.FailedLastRun && !failed {
		log.Info("Recovery")
		report.NotificationToSend = true
		manager.FailedLastRun = false
		report.WebHookMessage = "Recovered Failed InSpec Report"
		report.IFTTTWebHookMessage = generateRecoveryOfFailedIFTTTWebHookMessage(report)
		report.SlackWebhookMessage = generateRecoveryOfFailedSlackWebhookMessage(report)
		return report
	}

	return report
}

func generateFailedIFTTTWebHookMessage(report Report) string {
	failedProfilesName := "<None>"
	failedControlInfo := "<None>"
	failedProfiles := report.failedProfiles(0.0)
	if len(failedProfiles) > 0 {
		failedProfile := failedProfiles[0]
		failedProfilesName = failedProfile.Name

		failedControls := failedProfile.failedControls(0.0)
		if len(failedControls) > 0 {
			failedControl := failedControls[0]
			failedControlInfo = fmt.Sprintf("%s:%s", failedControl.ID, failedControl.Title)
		}
	}

	msg := IFTTTMessage{
		Value1: failedProfilesName,
		Value2: failedControlInfo,
		Value3: fmt.Sprintf("%d", report.numberOfFailedTests()),
	}

	JSONRaw, err := json.Marshal(msg)
	if err != nil {
		log.Errorf("Error parsing message %v", err)
		return ""
	}

	return string(JSONRaw)
}

func generateReminderOfFailedIFTTTWebHookMessage(report Report) string {
	failedProfilesName := "<None>"
	failedControlInfo := "<None>"
	failedProfiles := report.failedProfiles(0.0)
	if len(failedProfiles) > 0 {
		failedProfile := failedProfiles[0]
		failedProfilesName = failedProfile.Name

		failedControls := failedProfile.failedControls(0.0)
		if len(failedControls) > 0 {
			failedControl := failedControls[0]
			failedControlInfo = fmt.Sprintf("%s:%s", failedControl.ID, failedControl.Title)
		}
	}

	msg := IFTTTMessage{
		Value1: "Reinder of " + failedProfilesName,
		Value2: failedControlInfo,
		Value3: fmt.Sprintf("%d", report.numberOfFailedTests()),
	}

	JSONRaw, err := json.Marshal(msg)
	if err != nil {
		log.Errorf("Error parsing message %v", err)
		return ""
	}

	return string(JSONRaw)
}

func generateRecoveryOfFailedIFTTTWebHookMessage(report Report) string {
	msg := IFTTTMessage{
		Value1: "Recovery of Failed Node",
		Value2: "",
		Value3: "",
	}

	JSONRaw, err := json.Marshal(msg)
	if err != nil {
		log.Errorf("Error parsing message %v", err)
		return ""
	}

	return string(JSONRaw)
}

func generateFailedSlackWebhookMessage(report Report) string {
	failedProfilesName := "<None>"
	failedControlInfo := "<None>"
	failedProfiles := report.failedProfiles(0.0)
	if len(failedProfiles) > 0 {
		failedProfile := failedProfiles[0]
		failedProfilesName = failedProfile.Name

		failedControls := failedProfile.failedControls(0.0)
		if len(failedControls) > 0 {
			failedControl := failedControls[0]
			failedControlInfo = fmt.Sprintf("%s:%s", failedControl.ID, failedControl.Title)
		}
	}

	msg := SlackMessage{
		Username: "Notification Service",
		Text:     fmt.Sprintf("InSpec found a critical control failure on node %q", report.NodeName),
		IconURL:  "https://docs.chef.io/_static/chef_logo_v2.png",
		Attachments: []SlackAttachment{
			{
				Text: fmt.Sprintf("%d tests failed. Rerun the test locally for full details.",
					report.numberOfFailedTests()),
				Color: "warning",
				Fields: []SlackField{
					{
						Title: "Control ID::Title",
						Value: failedControlInfo,
						Short: false,
					},
					{
						Title: "Profile",
						Value: failedProfilesName,
						Short: false,
					},
					{
						Title: "Node",
						Value: report.NodeName,
						Short: false,
					},
					{
						Title: "Highest Failed Impact",
						Value: fmt.Sprintf("%.1f", report.highestFailedImpact()),
						Short: false,
					},
				},
			},
		},
	}

	JSONRaw, err := json.Marshal(msg)
	if err != nil {
		log.Errorf("Error parsing message %v", err)
		return fmt.Sprintf("{\"text\": \"%s\"}", "Error")
	}

	return string(JSONRaw)
}

func generateReminderOfFailedSlackWebhookMessage(report Report) string {
	failedProfilesName := "<None>"
	failedControlInfo := "<None>"
	failedProfiles := report.failedProfiles(0.0)
	if len(failedProfiles) > 0 {
		failedProfile := failedProfiles[0]
		failedProfilesName = failedProfile.Name

		failedControls := failedProfile.failedControls(0.0)
		if len(failedControls) > 0 {
			failedControl := failedControls[0]
			failedControlInfo = fmt.Sprintf("%s:%s", failedControl.ID, failedControl.Title)
		}
	}

	msg := SlackMessage{
		Username: "Notification Service",
		Text:     fmt.Sprintf("Reminder of InSpec found a critical control failure on node %q", report.NodeName),
		IconURL:  "https://docs.chef.io/_static/chef_logo_v2.png",
		Attachments: []SlackAttachment{
			{
				Text: fmt.Sprintf("%d tests failed. Rerun the test locally for full details.",
					report.numberOfFailedTests()),
				Color: "warning",
				Fields: []SlackField{
					{
						Title: "Control ID::Title",
						Value: failedControlInfo,
						Short: false,
					},
					{
						Title: "Profile",
						Value: failedProfilesName,
						Short: false,
					},
					{
						Title: "Node",
						Value: report.NodeName,
						Short: false,
					},
					{
						Title: "Highest Failed Impact",
						Value: fmt.Sprintf("%.1f", report.highestFailedImpact()),
						Short: false,
					},
				},
			},
		},
	}

	JSONRaw, err := json.Marshal(msg)
	if err != nil {
		log.Errorf("Error parsing message %v", err)
		return fmt.Sprintf("{\"text\": \"%s\"}", "Error")
	}

	return string(JSONRaw)
}

func generateRecoveryOfFailedSlackWebhookMessage(report Report) string {
	msg := SlackMessage{
		Username: "Notification Service",
		Text:     fmt.Sprintf("Recovery of InSpec critical control failure on node %q", report.NodeName),
		IconURL:  "https://docs.chef.io/_static/chef_logo_v2.png",
		Attachments: []SlackAttachment{
			{
				Text:   "Recovered Failure",
				Color:  "#2eb886",
				Fields: []SlackField{},
			},
		},
	}

	JSONRaw, err := json.Marshal(msg)
	if err != nil {
		log.Errorf("Error parsing message %v", err)
		return fmt.Sprintf("{\"text\": \"%s\"}", "Error")
	}

	return string(JSONRaw)
}
