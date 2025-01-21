package kafka

import (
	"fmt"
	"strings"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Producer struct {
	producer *kafka.Producer
}

func New(address []string) (*Producer, error) {
	conf := &kafka.ConfigMap{
		"bootstrap.servers": strings.Join(address, ","),
	}
	p, err := kafka.NewProducer(conf)
	if err != nil {
		return nil, err
	}
	return &Producer{producer: p}, nil
}

func (p *Producer) Produce(message []byte, topic string) error {
	kafkaMsg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          message,
		Key:            nil,
	}
	kafkaCh := make(chan kafka.Event)
	if err := p.producer.Produce(kafkaMsg, kafkaCh); err != nil {
		return err
	}
	e := <-kafkaCh
	switch ev := e.(type) {
	case *kafka.Message:
		return nil
	case kafka.Error:
		return ev
	default:
		return fmt.Errorf("unknown event type")
	}
}

func (p *Producer) Close() {
	p.producer.Flush(2000)
	p.producer.Close()
}
