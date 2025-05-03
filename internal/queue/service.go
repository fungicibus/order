package queue

import (
	"fmt"

	"github.com/IBM/sarama"

	"github.com/fungicibus/order/config"
	"github.com/fungicibus/order/internal/logger"
)

type service struct {
	cfg config.Kafka
	log *logger.Logger

	_        sarama.ConsumerGroup
	producer sarama.SyncProducer
}

func New(cfg config.Kafka, log *logger.Logger, clientID string) (*service, error) {
	svc := &service{
		cfg: cfg,
		log: log,
	}
	sarama.Logger = logger.WithSource(log, "sarama")

	config := sarama.NewConfig()
	config.Version = sarama.V4_0_0_0
	config.ClientID = clientID

	config.Producer.Partitioner = sarama.NewHashPartitioner
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true

	producer, err := sarama.NewSyncProducer(cfg.Brokers, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create producer: %w", err)
	}
	svc.producer = producer

	return svc, nil
}

func (s *service) PingBrokers() error {
	for _, brokerAddr := range s.cfg.Brokers {
		broker := sarama.NewBroker(brokerAddr)
		err := broker.Open(nil)
		if err != nil {
			return fmt.Errorf("failed to open broker connection: %w", err)
		}

		request := sarama.MetadataRequest{Topics: s.getTopics()}
		_, err = broker.GetMetadata(&request)
		if err != nil {
			_ = broker.Close()
			return fmt.Errorf("failed to get metadata: %w", err)
		}

		if err = broker.Close(); err != nil {
			return fmt.Errorf("failed to close broker connection: %w", err)
		}
	}

	return nil
}

func (s *service) getTopics() []string {
	return []string{
		s.cfg.TopicOrderCreated,
		s.cfg.TopicOrderConfirmed,
		s.cfg.TopicOrderRejected,
	}
}
