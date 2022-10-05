package store

import (
	mapset "github.com/deckarep/golang-set/v2"
)

type StatusStore struct {
	GlobalHeight int
	ProcessedSet *mapset.Set[string]
}

func NewStatusStore() *StatusStore {
	s := &StatusStore{
		GlobalHeight: 0,
	}
	set := mapset.NewSet[string]()
	s.ProcessedSet = &set
	return s
}
