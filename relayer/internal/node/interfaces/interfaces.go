package interfaces

import (
	"github.com/vitelabs/vite-portal/internal/core/generics"
	"github.com/vitelabs/vite-portal/internal/node/types"
)
type NodeI interface {
	GetName() string
	GetIpAddress() string
	GetRewardAddress() string
}

type ServiceI interface {
	GetNodes() generics.GenericPage[types.Node]
}