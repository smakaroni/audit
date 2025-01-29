package models

import (
	"google.golang.org/protobuf/types/known/timestamppb"
	"testing"
	"time"
)

func TestUnmarshalEvent(t *testing.T) {
	jsonData := []byte(`{"event_type": "test", "user": "john", "data": "sensitive"}`)
	event, err := UnmarshalEvent(jsonData)
	if err != nil {
		t.Fatalf("Failed to unmarshal event: %v", err)
	}

	if event["event_type"] != "test" {
		t.Errorf("Expected event_type 'test', got %v", event["event_type"])
	}
	if event["user"] != "john" {
		t.Errorf("Expected user 'john', got %v", event["user"])
	}
	if event["data"] != "sensitive" {
		t.Errorf("Expected data 'sensitive', got %v", event["data"])
	}
}

func TestNewAuditLog(t *testing.T) {
	event := map[string]interface{}{
		"event_type": "test",
		"user":       "john",
		"data":       "sensitive",
	}

	log := NewAuditLog(event)

	if log.EventType != "test" {
		t.Errorf("Expected EventType 'test', got %v", log.EventType)
	}
	if log.AnonymizedUser != "john" {
		t.Errorf("Expected AnonymizedUser 'john', got %v", log.AnonymizedUser)
	}
	if log.Timestamp == nil {
		t.Errorf("Timestamp is nil")
	} else if time.Since(log.Timestamp.AsTime()) > time.Second {
		t.Errorf("Timestamp is not recent")
	}
	if log.Data != "sensitive" {
		t.Errorf("Expected Data 'sensitive', got %v", log.Data)
	}
}

func TestToProtoAndFromProto(t *testing.T) {
	originalLog := &AuditLog{
		EventType:      "test",
		AnonymizedUser: "john",
		Timestamp:      timestamppb.Now(),
		Data:           "sensitive",
	}

	// Convert to protobuf
	protoData, err := ToProto(originalLog)
	if err != nil {
		t.Fatalf("Failed to convert to protobuf: %v", err)
	}

	// Convert back from protobuf
	restoredLog, err := AuditLogFromProto(protoData)
	if err != nil {
		t.Fatalf("Failed to convert from protobuf: %v", err)
	}

	// Compare original and restored
	if originalLog.EventType != restoredLog.EventType {
		t.Errorf("EventType mismatch. Expected %s, got %s", originalLog.EventType, restoredLog.EventType)
	}
	if originalLog.AnonymizedUser != restoredLog.AnonymizedUser {
		t.Errorf("AnonymizedUser mismatch. Expected %s, got %s", originalLog.AnonymizedUser, restoredLog.AnonymizedUser)
	}
	if !originalLog.Timestamp.AsTime().Equal(restoredLog.Timestamp.AsTime()) {
		t.Errorf("Timestamp mismatch. Expected %v, got %v", originalLog.Timestamp.AsTime(), restoredLog.Timestamp.AsTime())
	}
	if originalLog.Data != restoredLog.Data {
		t.Errorf("Data mismatch. Expected %s, got %s", originalLog.Data, restoredLog.Data)
	}
}
