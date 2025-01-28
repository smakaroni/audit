package kafka

import (
	"context"
	"github.com/segmentio/kafka-go"
)

//go:generate mockgen -destination=mock_consumer.go -package=kafka audit/internal/kafka Consumer

type Consumer interface {
	ReadMessage(ctx context.Context) (kafka.Message, error)
	Close() error
}

type KafkaConsumer struct {
	reader *kafka.Reader
}

func NewConsumer(broker, topic, groupID string) (Consumer, error) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{broker},
		Topic:     topic,
		GroupID:   groupID,
		Partition: 0,
	})

	return &KafkaConsumer{reader: r}, nil
}

func (c *KafkaConsumer) ReadMessage(ctx context.Context) (kafka.Message, error) {
	return c.reader.ReadMessage(ctx)
}

func (c *KafkaConsumer) Close() error {
	return c.reader.Close()
}
