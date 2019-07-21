package util

import (
	"github.com/buger/jsonparser"
)

func GetStringIfExists(fieldName string, body []byte) string {
	fieldString, err := jsonparser.GetString(body, fieldName)
	if err == nil {
		return fieldString
	}
	return ""
}

// Returns 0 if field does not exist
func GetFloat32IfExists(fieldName string, body []byte) float32 {
	fieldFloat64, err := jsonparser.GetFloat(body, fieldName)
	if err == nil {
		return float32(fieldFloat64)
	}
	return 0
}
