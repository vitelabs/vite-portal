package interfaces

import (
	"github.com/vitelabs/vite-portal/orchestrator/internal/node/types"
)

type StoreI interface {
	Clear()
	Close()
	Count() int
	GetByIndex(index int) (n types.Node, found bool)
	GetById(id string) (n types.Node, found bool)
	Upsert(n types.Node) error
	Remove(id string) error
}
