package orchestrator

import (
	urlutil "net/url"
	"strconv"
	"time"

	"github.com/vitelabs/vite-portal/relayer/internal/orchestrator/client"
	"github.com/vitelabs/vite-portal/shared/pkg/generics"
	"github.com/vitelabs/vite-portal/shared/pkg/handler"
	"github.com/vitelabs/vite-portal/shared/pkg/logger"
	"github.com/vitelabs/vite-portal/shared/pkg/rpc"
	"github.com/vitelabs/vite-portal/shared/pkg/ws"
)

type Orchestrator struct {
	relayerId string
	stopped   bool
	status    ws.ConnectionStatus
	client    *client.Client
	ps        *handler.Pubsub
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
		ps:        handler.NewPubsub(),
	}
}

func (o *Orchestrator) GetStatus() ws.ConnectionStatus {
	return o.status
}

func (o *Orchestrator) Start(apis []rpc.API) {
	o.stopped = false
	o.connect()

	o.client.RegisterNames(apis)

	go func() {
		<-o.client.Closed()
		o.setStatus(ws.Disconnected)
		if !o.stopped {
			time.Sleep(10 * time.Second)
			o.Start(apis)
		}
	}()
}

func (o *Orchestrator) Stop() {
	if !o.stopped {
		o.stopped = true
		o.client.Close()
	}
}

func (o *Orchestrator) SubscribeStatusChange() <-chan string {
	return o.ps.Subscribe()
}

func (o *Orchestrator) connect() {
	if o.stopped {
		return
	}
	o.setStatus(ws.Connecting)
	err := o.client.Connect(o.relayerId)
	if err != nil {
		logger.Logger().Error().Err(err).Msg("trying to connect to orchestrator")
		time.Sleep(10 * time.Second)
		o.connect()
		return
	}
	o.setStatus(ws.Connected)
}

func (o *Orchestrator) setStatus(newStatus ws.ConnectionStatus) {
	if o.status != newStatus {
		logger.Logger().Info().
			Int64("before", int64(o.status)).
			Int64("after", int64(newStatus)).
			Msg("connection status changed")
		o.status = newStatus
		o.ps.Publish(strconv.FormatInt(int64(newStatus), 10))
	}
}

func (o *Orchestrator) GetChains() []string {
	return make([]string, 0)
}

func (o *Orchestrator) GetNodes(chain string, offset int, limit int) (generics.GenericPage[string], error) {
	result := *generics.NewGenericPage[string]()
	return result, nil
}
