package queue

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/IBM/sarama"

	"github.com/fungicibus/order/internal/types"
)

func (s *service) EnqueueOrder(ctx context.Context, order types.Order) error {
	value, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("failed to marshal order message: %w", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: s.cfg.TopicOrderCreated,
		Key:   sarama.StringEncoder(order.Id),
		Value: sarama.ByteEncoder(value),
	}

	_, _, err = s.producer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("failed to send message to Kafka: %w", err)
	}

	return nil
}
