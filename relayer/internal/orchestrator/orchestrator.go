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
	StatusChanged chan ws.ConnectionStatus
	status        ws.ConnectionStatus
	client        *client.Client
}

func NewOrchestrator(url string, timeout time.Duration) *Orchestrator {
	return &Orchestrator{
		StatusChanged: make(chan ws.ConnectionStatus),
		status:        ws.Unknown,
		client:        client.NewClient(url, timeout),
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
	o.setStatus(ws.Connecting)
	err := o.client.Connect()
	if err != nil {
		// TODO: use exponential backoff strategy
		time.Sleep(1 * time.Second)
		o.init()
	}
	o.setStatus(ws.Connected)
	go o.handleMessages()
}

func (o *Orchestrator) handleMessages() {
	defer func() {
		o.client.Conn.Close()
		o.setStatus(ws.Disconnected)
	}()
	for {
		_, message, err := o.client.Conn.ReadMessage()
		if err != nil {
			break
		}
		logger.Logger().Info().Msg(fmt.Sprintf("message: %s", message))
	}
}

func (o *Orchestrator) setStatus(newStatus ws.ConnectionStatus) {
	if o.status != newStatus {
		o.status = newStatus
		o.StatusChanged <- newStatus
	}
}
