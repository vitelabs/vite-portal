package handler

import (
	"sync"
	"time"

	"github.com/vitelabs/vite-portal/shared/pkg/client"
	"github.com/vitelabs/vite-portal/shared/pkg/types"
	"github.com/vitelabs/vite-portal/shared/pkg/util/jsonutil"
)

type KafkaHandler struct {
	closed        bool
	mutex         sync.Mutex
	timeout       time.Duration
	defaultClient *client.KafkaClient
	rpcClient     *client.KafkaClient
}

func NewKafkaHandler(timeout time.Duration, cfg types.KafkaConfig) *KafkaHandler {
	return &KafkaHandler{
		closed:        false,
		timeout:       timeout,
		defaultClient: client.NewKafkaClient(timeout, cfg.Server, cfg.DefaultTopic),
		rpcClient:     client.NewKafkaClient(timeout, cfg.Server, cfg.RpcTopic),
	}
}

func (h *KafkaHandler) Close() {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	if h.closed {
		return
	}
	h.closed = true
	h.defaultClient.Close()
	h.rpcClient.Close()
}

func (h *KafkaHandler) ReadDefault(offset, limit, timeout int) ([]string, error) {
	return h.defaultClient.Read(int64(offset), limit, time.Duration(timeout)*time.Millisecond)
}

func (h *KafkaHandler) ReadRpc(offset, limit, timeout int) ([]string, error) {
	return h.rpcClient.Read(int64(offset), limit, time.Duration(timeout)*time.Millisecond)
}

func (h *KafkaHandler) WriteDefault(msg types.KafkaNodeOnlineStatus) {
	h.defaultClient.Write(jsonutil.ToString(msg))
}

func (h *KafkaHandler) WriteRpc(msg any) {
	h.rpcClient.Write(jsonutil.ToString(msg))
}
