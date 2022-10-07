package handler

import (
	"sync"
	"time"

	nodestore "github.com/vitelabs/vite-portal/orchestrator/internal/node/store"
	"github.com/vitelabs/vite-portal/orchestrator/internal/types"
	sharedclients "github.com/vitelabs/vite-portal/shared/pkg/client"
)

type Handler struct {
	client      *sharedclients.ViteClient
	nodeStore   *nodestore.MemoryStore
	statusStore *nodestore.StatusStore
	timeout     time.Duration
	heightLock  sync.Mutex
}

func NewHandler(cfg types.Config, client *sharedclients.ViteClient, nodeStore *nodestore.MemoryStore, statusStore *nodestore.StatusStore) *Handler {
	return &Handler{
		client:      client,
		nodeStore:   nodeStore,
		statusStore: statusStore,
		timeout:     time.Duration(cfg.RpcTimeout) * time.Millisecond,
		heightLock:  sync.Mutex{},
	}
}
