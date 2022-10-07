package handler

import (
	nodestore "github.com/vitelabs/vite-portal/orchestrator/internal/node/store"
	"github.com/vitelabs/vite-portal/orchestrator/internal/types"
)

type Handler struct {
	config      types.Config
	nodeStore   *nodestore.MemoryStore
	statusStore *nodestore.StatusStore
}

func NewHandler(cfg types.Config, nodeStore *nodestore.MemoryStore, statusStore *nodestore.StatusStore) *Handler {
	return &Handler{
		config:      cfg,
		nodeStore:   nodeStore,
		statusStore: statusStore,
	}
}
