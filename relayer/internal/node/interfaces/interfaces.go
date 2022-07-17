package interfaces

import (
	"github.com/vitelabs/vite-portal/internal/generics"
	"github.com/vitelabs/vite-portal/internal/node/types"
)

type ServiceI interface {
	GetNodes() generics.GenericPage[types.Node]
}

type StoreI interface {
	GetById(id string) (types.Node, bool)
	GetAllByChain(c string) []types.Node
	Upsert(n types.Node) error
	UpsertMany(nodes []types.Node) error
	Remove(id string) error
	Close()
}