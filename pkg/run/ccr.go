package run

import (
	"github.com/buger/jsonparser"
)

type Run struct {
	Message string
	Status  string
}

func (run Run) SendNotification() bool {
	return run.Status == "failure" && run.Message == "run_converge"
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
