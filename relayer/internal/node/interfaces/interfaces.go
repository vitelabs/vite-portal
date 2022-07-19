package interfaces

import (
	"github.com/vitelabs/vite-portal/internal/generics"
	"github.com/vitelabs/vite-portal/internal/node/types"
)

type ServiceI interface {
	GetChains() []string
	GetNodes(chain string, offset, limit int) (generics.GenericPage[types.Node], error)
}

type StoreI interface {
	GetChains() []string
	Get(chain string, id string) (types.Node, bool)
	GetByIndex(chain string, index int) (types.Node, bool)
	Upsert(n types.Node) error
	UpsertMany(nodes []types.Node) error
	Remove(chain string, id string) error
	Count(chain string) int
	Clear()
	Close()
}