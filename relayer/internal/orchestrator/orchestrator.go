package orchestrator

import (
	urlutil "net/url"
	"time"

	"github.com/vitelabs/vite-portal/relayer/internal/orchestrator/client"
	"github.com/vitelabs/vite-portal/shared/pkg/logger"
	"github.com/vitelabs/vite-portal/shared/pkg/rpc"
	"github.com/vitelabs/vite-portal/shared/pkg/ws"
)

type Orchestrator struct {
	relayerId string
	stopped   bool
	status    ws.ConnectionStatus
	client    *client.Client
}

func NewOrchestrator(relayerId, url, jwtSecret string, timeout, jwtExpiryTimeout time.Duration) *Orchestrator {
	u, e := urlutil.Parse(url)
	if e != nil {
		logger.Logger().Error().Err(e).Msg("orchestrator URL parse failed")
	}
	if u.Scheme != "ws" && u.Scheme != "wss" {
		logger.Logger().Error().Str("url", url).Msg("orchestrator URL does not match WebSocket protocol")
	}
	return &Orchestrator{
		relayerId: relayerId,
		stopped:   false,
		status:    ws.Unknown,
		client:    client.NewClient(url, jwtSecret, timeout, jwtExpiryTimeout),
	}
}

func (o *Orchestrator) GetStatus() ws.ConnectionStatus {
	return o.status
}

func (o *Orchestrator) Start(s *rpc.Server) {
	o.stopped = false
	o.connect(s)

	codec := rpc.NewFuncCodec(o.client.Conn, o.client.Conn.WriteJSON, o.client.Conn.ReadJSON)
	go s.ServeCodec(codec, 0, nil, nil)

	go func() {
		<-codec.Closed()
		o.setStatus(ws.Disconnected)
		if !o.stopped {
			time.Sleep(10 * time.Second)
			o.Start(s)
		}
	}()
}

func (o *Orchestrator) Stop() {
	if !o.stopped {
		o.stopped = true
		conn := o.client.Conn
		if conn != nil {
			conn.Close()
		}
	}
}

func (o *Orchestrator) connect(s *rpc.Server) {
	if o.stopped {
		return
	}
	o.setStatus(ws.Connecting)
	err := o.client.Connect(o.relayerId)
	if err != nil {
		logger.Logger().Error().Err(err).Msg("trying to connect to orchestrator")
		time.Sleep(10 * time.Second)
		o.connect(s)
		return
	}
	o.setStatus(ws.Connected)
	// TODO: notify subscribers
}

func (o *Orchestrator) setStatus(newStatus ws.ConnectionStatus) {
	if o.status != newStatus {
		logger.Logger().Info().
			Int64("before", int64(o.status)).
			Int64("after", int64(newStatus)).
			Msg("connection status changed")
		o.status = newStatus
	}
}
