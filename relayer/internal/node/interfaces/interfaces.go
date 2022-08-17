package interfaces

import (
	"github.com/vitelabs/vite-portal/relayer/internal/generics"
	"github.com/vitelabs/vite-portal/relayer/internal/node/types"
)

type ServiceI interface {
	GetChains() []string
	GetNodeCount(chain string) int
	GetNodes(chain string, offset, limit int) (generics.GenericPage[types.Node], error)
	GetNode(id string) (n types.Node, found bool)
	GetNodeByIndex(chain string, index int) (n types.Node, found bool)
	PutNode(n types.Node) error
	DeleteNode(id string) error
	LastActivityTimestamp(chain string, a types.NodeActivity) int64
}

type StoreI interface {
	GetChains() []string
	Get(chain string, id string) (types.Node, bool)
	GetByIndex(chain string, index int) (n types.Node, found bool)
	GetById(id string) (n types.Node, found bool)
	Upsert(n types.Node) error
	UpsertMany(nodes []types.Node) error
	Remove(chain string, id string) error
	Count(chain string) int
	Clear()
	Close()
}