package kafka

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/segmentio/kafka-go"
)

func TestNewConsumer(t *testing.T) {
	consumer, err := NewConsumer("localhost:9092", "test-topic", "test-group")
	if err != nil {
		t.Fatalf("Failed to create consumer: %v", err)
	}

	kafkaConsumer, ok := consumer.(*KafkaConsumer)
	if !ok {
		t.Fatal("Failed to cast consumer to KafkaConsumer")
	}

	if kafkaConsumer.reader == nil {
		t.Error("Consumer reader is nil")
	}
}

func TestReadMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConsumer := NewMockConsumer(ctrl)

	expectedMsg := kafka.Message{
		Topic: "test-topic",
		Value: []byte("test message"),
	}

	mockConsumer.EXPECT().
		ReadMessage(gomock.Any()).
		Return(expectedMsg, nil)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	msg, err := mockConsumer.ReadMessage(ctx)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if string(msg.Value) != "test message" {
		t.Errorf("Expected message 'test message', got '%s'", string(msg.Value))
	}
}

func TestClose(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConsumer := NewMockConsumer(ctrl)

	mockConsumer.EXPECT().
		Close().
		Return(nil)

	err := mockConsumer.Close()
	if err != nil {
		t.Errorf("Unexpected error on Close(): %v", err)
	}
}
