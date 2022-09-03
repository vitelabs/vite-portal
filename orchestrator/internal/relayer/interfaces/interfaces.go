package interfaces

import (
	"github.com/vitelabs/vite-portal/orchestrator/internal/relayer/types"
)

type StoreI interface {
	Clear()
	Close()
	Count() int
	GetByIndex(index int) (r types.Relayer, found bool)
	GetById(id string) (r types.Relayer, found bool)
	Upsert(r types.Relayer) error
	Remove(id string) error
}
