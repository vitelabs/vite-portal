package store

import (
	"errors"

	"github.com/vitelabs/vite-portal/internal/collections"
	"github.com/vitelabs/vite-portal/internal/logger"
	"github.com/vitelabs/vite-portal/internal/node/interfaces"
	"github.com/vitelabs/vite-portal/internal/node/types"
	"github.com/vitelabs/vite-portal/internal/util/jsonutil"
)

type KvStore struct {
	interfaces.StoreI
	db map[string]collections.NameObjectCollectionI
}

func NewKvStore() *KvStore {
	s := &KvStore{}
	s.Clear()
	return s
}

// ---
// Implement "StoreI" interface

func (s *KvStore) Get(chain string, id string) (n types.Node, found bool) {
	// Assign default return values
	n = *new(types.Node)
	found = false

	if chain == "" || id == "" || s.db[chain] == nil {
		return
	}

	node := s.db[chain].Get(id)
	if node == nil {
		return
	}

	return node.(types.Node), true
}

func (s *KvStore) Upsert(n types.Node) error {
	err := validateNode(n)
	if err != nil {
		return err
	}

	c := s.initChain(n.Chain)
	
	c.Add(n.Id, n)

	return nil
}

func (s *KvStore) Remove(chain string, id string) error {
	if chain == "" || id == "" || s.db[chain] == nil {
		return nil
	}

	s.db[chain].Remove(id)

	return nil
}

func (s *KvStore) Count(chain string) int {
	if s.db[chain] == nil {
		return 0
	}

	return s.db[chain].Count()
}

func (s *KvStore) Clear() {
	s.db = map[string]collections.NameObjectCollectionI{}
}

func (s *KvStore) Close() {

}

func (s *KvStore) initChain(chain string) (c collections.NameObjectCollectionI) {
	if s.db[chain] == nil {
		s.db[chain] = collections.NewNameObjectCollection()
	}

	return s.db[chain]
}

func validateNode(n types.Node) error {
	if !n.IsValid() {
		err := errors.New("Trying to insert invalid node")
		logger.Logger().Error().Err(err).Str("node", jsonutil.ToString(n))
		return err
	}
	return nil
}
