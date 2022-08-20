package orchestrator

import (
	"errors"
	"fmt"
	urlutil "net/url"
	"time"

	"github.com/vitelabs/vite-portal/relayer/internal/orchestrator/client"
	"github.com/vitelabs/vite-portal/shared/pkg/logger"
	"github.com/vitelabs/vite-portal/shared/pkg/ws"
)

type Orchestrator struct {
	status ws.ConnectionStatus
	client *client.Client
}

func NewOrchestrator(url string, timeout time.Duration) *Orchestrator {
	return &Orchestrator{
		client: client.NewClient(url, timeout),
	}
}

func InitOrchestrator(url string, timeout time.Duration) (*Orchestrator, error) {
	u, e := urlutil.Parse(url)
	if e != nil {
		return nil, e
	}
	if u.Scheme != "ws" && u.Scheme != "wss" {
		return nil, errors.New("URL need to match WebSocket Protocol.")
	}
	orchestrator := NewOrchestrator(url, timeout)
	go orchestrator.init()
	return orchestrator, nil
}

func (o *Orchestrator) GetStatus() ws.ConnectionStatus {
	return o.status
}

func (o *Orchestrator) init() {
	o.status = ws.Connecting
	err := o.client.Connect()
	if err != nil {
		// TODO: use use exponential backoff strategy
		time.Sleep(1 * time.Second)
		o.init()
	}
	o.status = ws.Connected
	o.client.Conn.SetCloseHandler(func(code int, text string) error {
		o.status = ws.Disconnected
		logger.Logger().Error().Msg(fmt.Sprintf("orchestrator connection closed with code: %d and text: %s", code, text))
		return nil
	})
	go o.handleMessages()
}

func (o *Orchestrator) handleMessages() {
	for {
		_, message, err := o.client.Conn.ReadMessage()
		if err != nil {
			break
		}
		logger.Logger().Info().Msg(fmt.Sprintf("message: %s", message))
	}
}
