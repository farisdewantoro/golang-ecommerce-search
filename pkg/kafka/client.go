package kafka

import (
	"log"

	"github.com/Shopify/sarama"
)

type Config struct {
	Brokers []string
	GroupID string
}

type Producer struct {
	producer sarama.SyncProducer
}

type Consumer struct {
	consumer sarama.Consumer
	groupID  string
}

func NewProducer(cfg *Config) (*Producer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true

	producer, err := sarama.NewSyncProducer(cfg.Brokers, config)
	if err != nil {
		return nil, err
	}

	return &Producer{
		producer: producer,
	}, nil
}

func NewConsumer(cfg *Config) (*Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer(cfg.Brokers, config)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		consumer: consumer,
		groupID:  cfg.GroupID,
	}, nil
}

func (p *Producer) SendMessage(topic string, value string) error {
	message := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(value),
	}
	log.Println("Sending message to topic:", topic)
	_, _, err := p.producer.SendMessage(message)
	return err
}

func (p *Producer) Close() error {
	return p.producer.Close()
}

func (c *Consumer) ConsumePartition(topic string, partition int32, offset int64) (sarama.PartitionConsumer, error) {
	return c.consumer.ConsumePartition(topic, partition, offset)
}

func (c *Consumer) Close() error {
	return c.consumer.Close()
}
