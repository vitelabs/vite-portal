package interfaces

import (
	nodestore "github.com/vitelabs/vite-portal/orchestrator/internal/node/store"
	relayerstore "github.com/vitelabs/vite-portal/orchestrator/internal/relayer/store"
)

type ContextI interface {
	GetNodeStore(chain string) (*nodestore.MemoryStore, error)
	GetRelayerStore() *relayerstore.MemoryStore
	GetStatusStore(chain string) (*nodestore.StatusStore, error)
}