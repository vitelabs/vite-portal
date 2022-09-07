package interfaces

import (
	"github.com/vitelabs/vite-portal/orchestrator/internal/node/types"
)

type StoreI interface {
	Clear()
	Close()
	Count(chain string) int
	GetChains() []string
	Get(chain string, id string) (types.Node, bool)
	GetByIndex(chain string, index int) (n types.Node, found bool)
	GetById(id string) (n types.Node, found bool)
	Add(n types.Node) error
	Remove(chain string, id string) error
}
