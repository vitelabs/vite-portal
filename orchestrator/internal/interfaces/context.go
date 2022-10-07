package interfaces

import (
	nodestore "github.com/vitelabs/vite-portal/orchestrator/internal/node/store"
	relayerstore "github.com/vitelabs/vite-portal/orchestrator/internal/relayer/store"
)

type ContextI interface {
	GetRelayerStore() *relayerstore.MemoryStore
	GetNodeStore(chain string) (*nodestore.MemoryStore, error)
	GetStatusStore(chain string) (*nodestore.StatusStore, error)
}