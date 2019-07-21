package run

import (
	"github.com/buger/jsonparser"
)

type Run struct {
	Message string
	Status  string
}

func (run Run) HasNotificationToSend() bool {
	return run.Status == "failure" && run.Message == "run_converge"
}

func (run Run) WebHookMessage() string {
	return "Failed Chef Client Run Report!"
}

func (run Run) IFTTTWebHookMessage() string {
	return "{\"value1\" : \"CCR\", \"value2\" : \"node name\", \"value3\" : \"resource failed\"}"
}

func (run Run) SlackWebhookMessage() string {
	return "{\"text\": \"Failed Chef Client Run Report!\"}"
}

func ParseRun(message []byte) Run {
	return Run{
		Message: getStringIfExists("message_type", message),
		Status:  getStringIfExists("status", message),
	}
}

func getStringIfExists(fieldName string, body []byte) string {
	fieldString, err := jsonparser.GetString(body, fieldName)
	if err == nil {
		return fieldString
	}
	return ""
}
