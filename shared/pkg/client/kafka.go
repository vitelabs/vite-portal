package client

import (
	"context"
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/vitelabs/vite-portal/shared/pkg/logger"
	sharedtypes "github.com/vitelabs/vite-portal/shared/pkg/types"
)

type KafkaClient struct {
	closed  bool
	mutex   sync.Mutex
	timeout time.Duration
	dialer  *kafka.Dialer
	writer  *kafka.Writer
	reader  *kafka.Reader
}

func NewKafkaClient(timeout time.Duration, cfg sharedtypes.KafkaServerConfig, topic sharedtypes.KafkaTopicConfig) *KafkaClient {
	if cfg.KeyStoreLocation != "" {
		logger.Logger().Fatal().Msg("TLS not implemented yet")
	}
	dialer := &kafka.Dialer{
		Timeout:   timeout,
		DualStack: true,
		// TLS: tlsConfig(),
	}
	servers := strings.Split(cfg.Servers, ",")
	writer := &kafka.Writer{
		Addr:                   kafka.TCP(servers...),
		Topic:                  topic.Topic,
		AllowAutoTopicCreation: true,
		ReadTimeout:            timeout,
		WriteTimeout:           timeout,
		Transport:              &kafka.Transport{
			// TLS: &tls.Config{},
		},
	}
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: servers,
		Topic:   topic.Topic,
		GroupID: topic.GroupId,
		Dialer:  dialer,
	})
	return &KafkaClient{
		closed:  false,
		timeout: timeout,
		dialer:  dialer,
		writer:  writer,
		reader:  reader,
	}
}

func (c *KafkaClient) Close() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.closed {
		return
	}
	c.closed = true
	c.writer.Close()
	c.reader.Close()
}

func (c *KafkaClient) Write(msg string) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	m := kafka.Message{
		Value: []byte(msg),
	}
	err := c.writer.WriteMessages(ctx, m)
	if err != nil && !errors.Is(err, context.DeadlineExceeded) {
		logger.Logger().Error().Err(err).Msg("failed to write message")
	}
}

func (c *KafkaClient) Read() []string {
	var messages []string
	for {
		ctx, cancelFn := context.WithTimeout(context.Background(), c.timeout)
		defer cancelFn()
		m, err := c.reader.ReadMessage(ctx)
		if err != nil {
			if !errors.Is(err, context.DeadlineExceeded) {
				logger.Logger().Error().Err(err).Msg("failed to read message")
			}
			break
		}
		messages = append(messages, string(m.Value))
	}
	return messages
}
