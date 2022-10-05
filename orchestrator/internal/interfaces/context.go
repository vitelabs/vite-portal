package interfaces

import (
	nodestore "github.com/vitelabs/vite-portal/orchestrator/internal/node/store"
	relayerstore "github.com/vitelabs/vite-portal/orchestrator/internal/relayer/store"
)

type ContextI interface {
	GetNodeStore() *nodestore.MemoryStore
	GetRelayerStore() *relayerstore.MemoryStore
	GetStatusStore(chain string) *nodestore.StatusStore
}