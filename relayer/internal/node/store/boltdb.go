package store

import (
	"github.com/vitelabs/vite-portal/internal/node/interfaces"
	"github.com/vitelabs/vite-portal/internal/node/types"
)

type BoltdbStore struct {
	interfaces.StoreI
	db map[string]map[string]types.Node
}

func NewBoltdbStore() BoltdbStore {
	return BoltdbStore{
	}
}

// ---
// Implement "StoreI" interface

func (s BoltdbStore) GetById(id string) (n types.Node, found bool) {
	return types.Node{}, false
}

func (s BoltdbStore) GetAllByChain(c string) []types.Node {
	// https://github.com/boltdb/bolt#iterating-over-keys
	// Use cursor to get random nodes
	return []types.Node{}
}

func (s BoltdbStore) Upsert(n types.Node) error {
	return nil
}

func (s BoltdbStore) UpsertMany(nodes []types.Node) error {
	return nil
}

func (s BoltdbStore) Remove(id string) error {
	return nil
}

func (s BoltdbStore) Clear() error {
	return nil
}

func (s BoltdbStore) Count() int {
	return 0
}

func (s BoltdbStore) Close() {

}