package kafka

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
	"go.uber.org/zap"
)

// kafkaWriter interface abstracts kafka.Writer for testability
type kafkaWriter interface {
	WriteMessages(ctx context.Context, msgs ...kafka.Message) error
	Close() error
	Stats() kafka.WriterStats
}

type Producer struct {
	writer kafkaWriter
	logger *zap.Logger
	topic  string // Store topic separately for logging
}

type ProducerConfig struct {
	Brokers  []string
	Topic    string
	Logger   *zap.Logger
	Username string
	Password string
	UseTLS   bool
}

type Message struct {
	Key   string
	Value interface{}
}

func NewProducer(cfg ProducerConfig) *Producer {
	dialer := &kafka.Dialer{
		Timeout:   10 * time.Second,
		DualStack: true,
	}

	// Configuration SASL/SSL for Confluent Cloud
	if cfg.UseTLS && cfg.Username != "" && cfg.Password != "" {
		dialer.SASLMechanism = plain.Mechanism{
			Username: cfg.Username,
			Password: cfg.Password,
		}
		dialer.TLS = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
	}

	// Create transport with dialer if TLS is enabled
	var transport *kafka.Transport
	if cfg.UseTLS && cfg.Username != "" && cfg.Password != "" {
		transport = &kafka.Transport{
			Dial: func(ctx context.Context, network, addr string) (net.Conn, error) {
				conn, err := dialer.DialContext(ctx, network, addr)
				if err != nil {
					return nil, err
				}
				return conn, nil
			},
		}
	}

	writer := &kafka.Writer{
		Addr:         kafka.TCP(cfg.Brokers...),
		Topic:        cfg.Topic,
		Balancer:     &kafka.LeastBytes{},
		MaxAttempts:  3,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
		RequiredAcks: kafka.RequireOne,
		Async:        false,
	}

	if transport != nil {
		writer.Transport = transport
	}

	return &Producer{
		writer: writer,
		logger: cfg.Logger,
		topic:  cfg.Topic,
	}
}

// Publish sends a message to Kafka with retries
func (p *Producer) Publish(ctx context.Context, key string, value interface{}) error {
	valueBytes, err := json.Marshal(value)
	if err != nil {
		if p.logger != nil {
			p.logger.Error("Failed to marshal message",
				zap.String("key", key),
				zap.Error(err),
			)
		}
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	msg := kafka.Message{
		Key:   []byte(key),
		Value: valueBytes,
		Time:  time.Now(),
	}

	if p.logger != nil {
		p.logger.Debug("Publishing message to Kafka",
			zap.String("topic", p.topic),
			zap.String("key", key),
		)
	}

	err = p.writer.WriteMessages(ctx, msg)
	if err != nil {
		if p.logger != nil {
			p.logger.Error("Failed to publish message",
				zap.String("topic", p.topic),
				zap.String("key", key),
				zap.Error(err),
			)
		}
		return fmt.Errorf("failed to publish message: %w", err)
	}

	if p.logger != nil {
		p.logger.Info("Message published successfully",
			zap.String("topic", p.topic),
			zap.String("key", key),
		)
	}

	return nil
}

// PublishBatch sends multiple messages in a batch
func (p *Producer) PublishBatch(ctx context.Context, messages []Message) error {
	kafkaMessages := make([]kafka.Message, len(messages))

	for i, msg := range messages {
		valueBytes, err := json.Marshal(msg.Value)
		if err != nil {
			if p.logger != nil {
				p.logger.Error("Failed to marshal batch message",
					zap.Int("index", i),
					zap.Error(err),
				)
			}
			return fmt.Errorf("failed to marshal batch message at index %d: %w", i, err)
		}

		kafkaMessages[i] = kafka.Message{
			Key:   []byte(msg.Key),
			Value: valueBytes,
			Time:  time.Now(),
		}
	}

	err := p.writer.WriteMessages(ctx, kafkaMessages...)
	if err != nil {
		if p.logger != nil {
			p.logger.Error("Failed to publish batch",
				zap.Int("count", len(messages)),
				zap.Error(err),
			)
		}
		return fmt.Errorf("failed to publish batch: %w", err)
	}

	if p.logger != nil {
		p.logger.Info("Batch published successfully",
			zap.Int("count", len(messages)),
		)
	}

	return nil
}

// Close gracefully shuts down the producer
func (p *Producer) Close() error {
	if p.logger != nil {
		p.logger.Info("Closing Kafka producer")
	}
	return p.writer.Close()
}

// Stats returns producer statistics
func (p *Producer) Stats() kafka.WriterStats {
	return p.writer.Stats()
}
