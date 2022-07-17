package store

import (
	"github.com/vitelabs/vite-portal/internal/node/interfaces"
	"github.com/vitelabs/vite-portal/internal/node/types"
)

type KvStore struct {
	interfaces.StoreI
	db map[string]map[string]types.Node
}

func NewKvStore() KvStore {
	return KvStore{
	}
}

// ---
// Implement "StoreI" interface

func (s KvStore) GetById(id string) (n types.Node, found bool) {
	return types.Node{}, false
}

func (s KvStore) GetAllByChain(c string) []types.Node {
	return []types.Node{}
}

func (s KvStore) Upsert(n types.Node) error {
	return nil
}

func (s KvStore) UpsertMany(nodes []types.Node) error {
	return nil
}

func (s KvStore) Remove(id string) error {
	return nil
}

func (s KvStore) Clear() error {
	return nil
}

func (s KvStore) Count() int {
	return 0
}

func (s KvStore) Close() {

}