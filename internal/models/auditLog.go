package models

import (
	"audit/internal/protos"
	"encoding/json"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type AuditLog = protos.AuditLog

func UnmarshalEvent(data []byte) (map[string]interface{}, error) {
	var event map[string]interface{}
	err := json.Unmarshal(data, &event)
	return event, err
}

func NewAuditLog(event map[string]interface{}) *AuditLog {
	return &AuditLog{
		EventType:      event["event_type"].(string),
		AnonymizedUser: event["user"].(string),
		Timestamp:      timestamppb.New(time.Now()),
		Data:           event["data"].(string),
	}
}

func ToProto(a *AuditLog) ([]byte, error) {
	return proto.Marshal(a)
}

func AuditLogFromProto(data []byte) (*AuditLog, error) {
	a := &AuditLog{}
	err := proto.Unmarshal(data, a)
	return a, err
}
