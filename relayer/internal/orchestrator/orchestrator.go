package orchestrator

import (
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
	u, e := urlutil.Parse(url)
	if e != nil {
		logger.Logger().Error().Err(e).Msg("orchestrator URL parse failed")
	}
	if u.Scheme != "ws" && u.Scheme != "wss" {
		logger.Logger().Error().Msg("orchestrator URL does not match WebSocket protocol")
	}
	return &Orchestrator{
		StatusChanged: make(chan ws.ConnectionStatus),
		status:        ws.Unknown,
		client:        client.NewClient(url, timeout),
	}
}

func (o *Orchestrator) GetStatus() ws.ConnectionStatus {
	return o.status
}

func (o *Orchestrator) Start() {
	o.setStatus(ws.Connecting)
	err := o.client.Connect()
	if err != nil {
		logger.Logger().Error().Err(err).Msg("trying to connect to orchestrator")
		// TODO: use exponential backoff strategy
		time.Sleep(1 * time.Second)
		// o.init()
		return
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
		select {
		case o.StatusChanged <- newStatus:
		default:
		}
	}
}
