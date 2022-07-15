package interfaces

import (
	"github.com/vitelabs/vite-portal/internal/core/generics"
	"github.com/vitelabs/vite-portal/internal/nodes/types"
)
type NodeI interface {
	GetName() string
	GetIpAddress() string
	GetRewardAddress() string
}

type KeeperI interface {
	GetNodes() generics.GenericPage[types.Node]
}