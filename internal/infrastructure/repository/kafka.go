package repository

import (
	"anyompt/internal/config"
	"anyompt/internal/domain"
	"context"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/rs/zerolog/log"
)

// KafkaRepository implements the ProducerRepository interface for Kafka.
type KafkaRepository struct {
	producer *kafka.Producer
	topic    string
}

// NewKafkaRepository creates a new instance of the Kafka repository.
func NewKafkaRepository(cfg config.Config) (domain.ProducerRepository, error) {
	if !cfg.KafkaEnabled {
		log.Info().Msg("Kafka is disabled. Skipping producer creation.")
		return nil, nil
	}

	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": cfg.KafkaBrokers})
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka producer: %w", err)
	}

	log.Info().Msgf("Successfully created Kafka producer for brokers: %s", cfg.KafkaBrokers)

	return &KafkaRepository{
		producer: p,
		topic:    cfg.KafkaTopic,
	}, nil
}

// Produce a message to a Kafka topic.
func (r *KafkaRepository) Produce(ctx context.Context, message []byte) error {
	if r.producer == nil {
		log.Warn().Msg("Kafka producer is not initialized; cannot send messages.")
		return nil
	}
	correlationID := ctx.Value("correlation_id").(string)
	routingID := ctx.Value("routing_id").(string)

	headers := []kafka.Header{
		{Key: "correlation_id", Value: []byte(correlationID)},
		{Key: "origin", Value: []byte("anyompt")},
	}
	deliveryChan := make(chan kafka.Event)
	err := r.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &r.topic, Partition: kafka.PartitionAny},
		Value:          message,
		Headers:        headers,
		Key:            []byte(routingID),
	}, deliveryChan)

	if err != nil {
		return fmt.Errorf("failed to produce message: %w", err)
	}

	e := <-deliveryChan
	m := e.(*kafka.Message)

	if m.TopicPartition.Error != nil {
		return fmt.Errorf("delivery failed: %v", m.TopicPartition.Error)
	}

	log.Debug().Msgf("Delivered message to topic %s [%d] at offset %v",
		*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)

	close(deliveryChan)
	return nil
}

// Close closes the Kafka producer.
func (r *KafkaRepository) Close() {
	if r.producer != nil {
		r.producer.Close()
	}
}
