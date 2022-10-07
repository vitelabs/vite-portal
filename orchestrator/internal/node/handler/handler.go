package handler

import (
	"time"

	nodestore "github.com/vitelabs/vite-portal/orchestrator/internal/node/store"
	"github.com/vitelabs/vite-portal/orchestrator/internal/types"
)

type Handler struct {
	nodeStore   *nodestore.MemoryStore
	statusStore *nodestore.StatusStore
	timeout     time.Duration
}

func NewHandler(cfg types.Config, nodeStore *nodestore.MemoryStore, statusStore *nodestore.StatusStore) *Handler {
	return &Handler{
		nodeStore:   nodeStore,
		statusStore: statusStore,
		timeout:     time.Duration(cfg.RpcTimeout) * time.Millisecond,
	}
}
