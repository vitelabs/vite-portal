package store

import (
	"errors"

	"github.com/tidwall/buntdb"

	"github.com/vitelabs/vite-portal/internal/logger"
	"github.com/vitelabs/vite-portal/internal/node/interfaces"
	"github.com/vitelabs/vite-portal/internal/node/types"
	"github.com/vitelabs/vite-portal/internal/util/jsonutil"
)

type BuntdbStore struct {
	interfaces.StoreI
	db *buntdb.DB
}

func NewBuntdbStore() BuntdbStore {
	db, err := buntdb.Open(":memory:") // Open a file that does not persist to disk
	db.CreateIndex("chain", "*", buntdb.IndexJSON("chain"))
	if err != nil {
		logger.Logger().Fatal().Err(err).Msg("Buntdb creation failed")
	}
	return BuntdbStore{
		db: db,
	}
}

// ---
// Implement "StoreI" interface

func (s BuntdbStore) GetById(id string) (n types.Node, found bool) {
	err := s.db.View(func(tx *buntdb.Tx) error {
		val, err := tx.Get(id)
		if err != nil {
			return err
		}
		err = jsonutil.FromByte([]byte(val), &n)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return types.Node{}, false
	}
	return n, true
}

func (s BuntdbStore) GetAllByChain(c string) []types.Node {
	return []types.Node{}
}

func (s BuntdbStore) Upsert(n types.Node) error {
	return s.UpsertMany([]types.Node{n})
}

func (s BuntdbStore) UpsertMany(nodes []types.Node) error {
	return s.db.Update(func(tx *buntdb.Tx) error {
		for _, n := range nodes {
			if n == (types.Node{}) {
				return errors.New("Empty node")
			}
			json, err :=  jsonutil.ToString(n)
			if err != nil {
				return err
			}
			_, _, err = tx.Set(n.Id, json, nil)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (s BuntdbStore) Remove(id string) error {
	return s.db.Update(func(tx *buntdb.Tx) error {
		_, err := tx.Delete(id)
		return err
	})
}

func (s BuntdbStore) Clear() error {
	return s.db.Update(func(tx *buntdb.Tx) error {
		return tx.DeleteAll()
	})
}

func (s BuntdbStore) Count() int {
	var res int
	s.db.View(func(tx *buntdb.Tx) error {
		len, err := tx.Len()
		res = len
		return err
	})
	return res
}

func (s BuntdbStore) Close() {
	s.db.Close()
}