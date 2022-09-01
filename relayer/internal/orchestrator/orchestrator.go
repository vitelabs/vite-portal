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
		status:        ws.Unknown,
		client:        client.NewClient(url, timeout),
	}
}

func (o *Orchestrator) GetStatus() ws.ConnectionStatus {
	return o.status
}

func (o *Orchestrator) Start(s *rpc.Server) {
	// TODO: start/stop properly
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

	codec := rpc.NewFuncCodec(o.client.Conn, o.client.Conn.WriteJSON, o.client.Conn.ReadJSON)
	go s.ServeCodec(codec, 0, nil, nil)

	go func() {
		<-codec.Closed()
		o.setStatus(ws.Disconnected)
	}()
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
