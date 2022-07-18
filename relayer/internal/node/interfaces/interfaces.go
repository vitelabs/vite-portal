package interfaces

import (
	"github.com/vitelabs/vite-portal/internal/generics"
	"github.com/vitelabs/vite-portal/internal/node/types"
)

type ServiceI interface {
	GetNodes() generics.GenericPage[types.Node]
}

type StoreI interface {
	Get(chain string, id string) (types.Node, bool)
	Upsert(n types.Node) error
	UpsertMany(nodes []types.Node) error
	Remove(chain string, id string) error
	Count(chain string) int
	Clear()
	Close()
}