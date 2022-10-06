package store

import (
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/vitelabs/vite-portal/orchestrator/internal/node/types"
)

type StatusStore struct {
	GlobalHeight *types.ChainHeight
	ProcessedSet *mapset.Set[string]
}

func NewStatusStore() *StatusStore {
	s := &StatusStore{
		GlobalHeight: types.NewChainHeight(),
	}
	set := mapset.NewSet[string]()
	s.ProcessedSet = &set
	return s
}
