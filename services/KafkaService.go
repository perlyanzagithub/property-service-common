package services

import (
	"github.com/segmentio/kafka-go"
	"golang.org/x/net/context"
	"log"
)

type KafkaService interface {
	SendMessage(ctx context.Context, key, value string) error
	ConsumeMessage(ctx context.Context, handler func(kafka.Message) error)
}

type kafkaService struct {
	reader *kafka.Reader
	writer *kafka.Writer
}

func NewKafkaService(reader *kafka.Reader, writer *kafka.Writer) KafkaService {
	return &kafkaService{reader: reader, writer: writer}

}
func (c *kafkaService) ConsumeMessage(ctx context.Context, handler func(kafka.Message) error) {
	for {
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			log.Printf("Error reading message from Kafka: %v", err)
			continue
		}
		if err := handler(msg); err != nil {
			log.Printf("Error handling message: %v", err)
		}
	}
}
func (c *kafkaService) SendMessage(ctx context.Context, key, value string) error {
	msg := kafka.Message{
		Key:   []byte(key),
		Value: []byte(value),
	}
	return c.writer.WriteMessages(ctx, msg)
}
