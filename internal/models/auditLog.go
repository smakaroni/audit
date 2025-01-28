package models

import (
	"encoding/json"
	"time"
)

type AuditLog struct {
	EventType      string
	AnonymizedUser string
	Timestamp      time.Time
	Data           string
}

func UnmarshalEvent(data []byte) (map[string]interface{}, error) {
	var event map[string]interface{}
	err := json.Unmarshal(data, &event)
	return event, err
}

func NewAuditLog(event map[string]interface{}) *AuditLog {
	return &AuditLog{
		EventType:      event["event_type"].(string),
		AnonymizedUser: event["user"].(string),
		Timestamp:      time.Now(),
		Data:           event["data"].(string),
	}
}
