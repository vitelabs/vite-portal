package handler

import (
	"sync"
	"time"

	"github.com/vitelabs/vite-portal/orchestrator/internal/interfaces"
	nodestore "github.com/vitelabs/vite-portal/orchestrator/internal/node/store"
	"github.com/vitelabs/vite-portal/orchestrator/internal/types"
	sharedclients "github.com/vitelabs/vite-portal/shared/pkg/client"
	sharedkafka "github.com/vitelabs/vite-portal/shared/pkg/kafka"
)

type Handler struct {
	client      *sharedclients.ViteClient
	kafka       *sharedkafka.Handler
	nodeStore   *nodestore.MemoryStore
	statusStore *nodestore.StatusStore
	timeout     time.Duration
	heightLock  sync.Mutex
}

func NewHandler(cfg types.Config, client *sharedclients.ViteClient, kafka *sharedkafka.Handler, ctx interfaces.ChainContextI) *Handler {
	timeout := time.Duration(cfg.RpcTimeout) * time.Millisecond
	return &Handler{
		client:      client,
		kafka:       kafka,
		nodeStore:   ctx.GetNodeStore(),
		statusStore: ctx.GetStatusStore(),
		timeout:     timeout,
		heightLock:  sync.Mutex{},
	}
}
