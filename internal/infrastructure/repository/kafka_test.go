package repository

import (
	"anyompt/config"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewKafkaRepository_KafkaDisabled(t *testing.T) {
	cfg := config.Config{KafkaEnabled: false}
	repo, err := NewKafkaRepository(cfg)
	assert.NoError(t, err)
	assert.Nil(t, repo)
}

func TestKafkaRepository_Produce_NilProducer(t *testing.T) {
	repo := &KafkaRepository{}
	ctx := context.Background()
	message := []byte("test message")

	err := repo.Produce(ctx, message)
	assert.NoError(t, err) // Should return nil error and log a warning
}

func TestKafkaRepository_Close_NilProducer(t *testing.T) {
	repo := &KafkaRepository{}
	// Should not panic or cause an error
	repo.Close()
}
