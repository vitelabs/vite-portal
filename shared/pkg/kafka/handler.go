package kafka

import (
	"sync"
	"time"

	"github.com/vitelabs/vite-portal/shared/pkg/types"
	"github.com/vitelabs/vite-portal/shared/pkg/util/jsonutil"
)

type Handler struct {
	closed        bool
	mutex         sync.Mutex
	timeout       time.Duration
	defaultClient *Client
	rpcClient     *Client
}

func NewHandler(timeout time.Duration, cfg types.KafkaConfig) *Handler {
	return &Handler{
		closed:        false,
		timeout:       timeout,
		defaultClient: NewClient(timeout, cfg.Server, cfg.DefaultTopic),
		rpcClient:     NewClient(timeout, cfg.Server, cfg.RpcTopic),
	}
}

func (h *Handler) Close() {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	if h.closed {
		return
	}
	h.closed = true
	h.defaultClient.Close()
	h.rpcClient.Close()
}

func (h *Handler) ReadDefault(offset, limit, timeout int) ([]string, error) {
	return h.defaultClient.Read(int64(offset), limit, time.Duration(timeout)*time.Millisecond)
}

func (h *Handler) ReadRpc(offset, limit, timeout int) ([]string, error) {
	return h.rpcClient.Read(int64(offset), limit, time.Duration(timeout)*time.Millisecond)
}

func (h *Handler) WriteDefault(msg types.KafkaNodeOnlineStatus) {
	h.defaultClient.Write(jsonutil.ToString(msg))
}

func (h *Handler) WriteRpc(msg any) {
	h.rpcClient.Write(jsonutil.ToString(msg))
}
