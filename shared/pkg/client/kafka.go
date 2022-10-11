package client

import (
	"context"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/vitelabs/vite-portal/shared/pkg/logger"
	sharedtypes "github.com/vitelabs/vite-portal/shared/pkg/types"
)

type KafkaClient struct {
	stopped   bool
	partition int
	timeout   time.Duration
	config    sharedtypes.KafkaConfig
	dialer    *kafka.Dialer
	conn      *kafka.Conn
}

func NewKafkaClient(timeout time.Duration, cfg sharedtypes.KafkaConfig) *KafkaClient {
	if cfg.KeyStoreLocation != "" {
		logger.Logger().Fatal().Msg("TLS not implemented yet")
	}
	dialer := &kafka.Dialer{
		Timeout:   timeout,
		DualStack: true,
		// TLS: tlsConfig(),
	}
	return &KafkaClient{
		stopped:   false,
		partition: 0,
		timeout:   timeout,
		config:    cfg,
		dialer:    dialer,
	}
}

func (c *KafkaClient) Start() {
	c.stopped = false
	c.connect()
}

func (c *KafkaClient) connect() {
	if c.stopped {
		return
	}
	conn, err := c.dialer.DialLeader(context.Background(), "tcp", c.config.Servers, c.config.Topic, c.partition)
	if err != nil {
		logger.Logger().Error().Err(err).Msg("trying to connect to kafka")
		time.Sleep(10 * time.Second)
		c.connect()
		return
	}
	c.conn = conn
}

func (c *KafkaClient) Stop() {
	if !c.stopped {
		c.stopped = true
		if err := c.conn.Close(); err != nil {
			logger.Logger().Error().Err(err).Msg("failed to close connection")
		}
	}
}

func (c *KafkaClient) Write(msg string) {
	c.conn.SetWriteDeadline(time.Now().Add(c.timeout))
	m := kafka.Message{
		Value: []byte(msg),
	}
	_, err := c.conn.WriteMessages(m)
	if err != nil {
		logger.Logger().Error().Err(err).Msg("failed to write message")
	}
}

func (c *KafkaClient) Read() []string {
	c.conn.SetReadDeadline(time.Now().Add(c.timeout))
	batch := c.conn.ReadBatch(10e3, 1e6) // fetch 10KB min, 1MB max
	var messages []string
	b := make([]byte, 10e3) // 10KB max per message
	for {
		n, err := batch.Read(b)
		if err != nil {
			break
		}
		m := string(b[:n])
		messages = append(messages, m)
	}
	return messages
}
