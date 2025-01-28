package models

import (
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
	if time.Since(log.Timestamp) > time.Second {
		t.Errorf("Timestamp is not recent")
	}
	if log.Data != "sensitive" {
		t.Errorf("Expected Data 'sensitive', got %v", log.Data)
	}
}
