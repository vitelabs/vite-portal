package interfaces

import (
	"github.com/vitelabs/vite-portal/orchestrator/internal/relayer/types"
	"github.com/vitelabs/vite-portal/shared/pkg/generics"
)

type StoreI interface {
	Clear()
	Close()
	Count() int
	GetByIndex(index int) (r types.Relayer, found bool)
	GetById(id string) (r types.Relayer, found bool)
	GetPaginated(offset, limit int) (generics.GenericPage[types.Relayer], error)
	Upsert(r types.Relayer) error
	Remove(id string) error
}
