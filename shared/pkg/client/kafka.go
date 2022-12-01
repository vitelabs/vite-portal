package client

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"
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
	tlsConfig := newTLSConfig(cfg)
	dialer := &kafka.Dialer{
		Timeout:   timeout,
		DualStack: true,
		TLS:       tlsConfig,
	}
	servers := strings.Split(cfg.Servers, ",")
	writer := &kafka.Writer{
		Addr:                   kafka.TCP(servers...),
		Topic:                  topic.Topic,
		AllowAutoTopicCreation: true,
		ReadTimeout:            timeout,
		WriteTimeout:           timeout,
		Transport: &kafka.Transport{
			TLS: tlsConfig,
		},
	}
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: servers,
		Topic:   topic.Topic,
		// Set GroupId if offsets should be managed by the broker
		// GroupID: topic.GroupId,
		Dialer: dialer,
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

func (c *KafkaClient) Read(offset int64, limit int, timeout time.Duration) ([]string, error) {
	start := time.Now()
	c.mutex.Lock()
	defer c.mutex.Unlock()

	offsetBefore := c.reader.Offset()
	c.reader.SetOffset(offset)
	logger.Logger().Debug().Int64("offset_before", offsetBefore).Int64("offset", offset).Int("limit", limit).Msg("read kafka messages started")

	if timeout.Milliseconds() > c.timeout.Milliseconds() {
		timeout = c.timeout
	}

	var messages []string
	var err error
	for {
		ctx, cancelFn := context.WithTimeout(context.Background(), timeout)
		defer cancelFn()
		m, err := c.reader.ReadMessage(ctx)
		if err != nil {
			if !errors.Is(err, context.DeadlineExceeded) {
				logger.Logger().Error().Err(err).Msg("failed to read message")
			}
			break
		}
		messages = append(messages, string(m.Value))
		if len(messages) >= limit {
			break
		}
	}
	elapsed := time.Since(start)
	logger.Logger().Debug().Int64("elapsed", elapsed.Milliseconds()).Int("count", len(messages)).Msg("read kafka messages ended")
	return messages, err
}

func (c *KafkaClient) Write(msgs ...string) {
	if len(msgs) == 0 {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	m := make([]kafka.Message, len(msgs))
	for i := 0; i < len(msgs); i++ {
		m[i] = kafka.Message{
			Value: []byte(msgs[i]),
		}
	}
	err := c.writer.WriteMessages(ctx, m...)
	if err != nil {
		logger.Logger().Error().Err(err).Msg("failed to write message")
	}
}

func newTLSConfig(cfg sharedtypes.KafkaServerConfig) *tls.Config {
	if cfg.CertLocation == "" || cfg.CertKeyLocation == "" || cfg.CertPoolLocation == "" {
		return nil
	}

	certPEM, _ := ioutil.ReadFile(cfg.CertLocation)
	keyPEM, _ := ioutil.ReadFile(cfg.CertKeyLocation)
	certificate, err := tls.X509KeyPair([]byte(certPEM), []byte(keyPEM))
	if err != nil {
		logger.Logger().Fatal().Err(err).Msg("failed to read certs")
	}

	caPEM, _ := ioutil.ReadFile(cfg.CertKeyLocation)
	caCertPool := x509.NewCertPool()
	if ok := caCertPool.AppendCertsFromPEM([]byte(caPEM)); !ok {
		logger.Logger().Fatal().Err(err).Msg("failed to read cert pool")
	}

	return &tls.Config{
		Certificates:       []tls.Certificate{certificate},
		RootCAs:            caCertPool,
		InsecureSkipVerify: true,
	}
}
