package models

import (
	"audit/internal/protos"
	"encoding/json"
	"fmt"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"reflect"
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

func CompareAuditLogs(current *AuditLog, previous *AuditLog) (map[string]interface{}, error) {
	currentData, err := UnmarshalData(current)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling current data: %v", err)
	}

	previousData, err := UnmarshalData(previous)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling previous data: %v", err)
	}

	changes := make(map[string]interface{})

	for key, currentValue := range currentData {
		if previousValue, exists := previousData[key]; exists {
			if !reflect.DeepEqual(currentValue, previousValue) {
				changes[key] = map[string]interface{}{
					"old": previousValue,
					"new": currentValue,
				}
			}
		} else {
			changes[key] = map[string]interface{}{
				"old": nil,
				"new": currentValue,
			}
		}
	}

	for key, previousValue := range previousData {
		if _, exists := currentData[key]; !exists {
			changes[key] = map[string]interface{}{
				"old": previousValue,
				"new": nil,
			}
		}
	}

	return changes, nil
}

func UnmarshalData(a *AuditLog) (map[string]interface{}, error) {
	var data map[string]interface{}
	err := json.Unmarshal([]byte(a.Data), &data)
	return data, err
}
