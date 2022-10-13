package interfaces

import (
	nodestore "github.com/vitelabs/vite-portal/orchestrator/internal/node/store"
	relayerstore "github.com/vitelabs/vite-portal/orchestrator/internal/relayer/store"
)

type ContextI interface {
	GetRelayerStore() *relayerstore.MemoryStore
	GetChainContext(chain string) (ChainContextI, error)
}

type ChainContextI interface {
	GetNodeStore() *nodestore.MemoryStore
	GetStatusStore() *nodestore.StatusStore
}