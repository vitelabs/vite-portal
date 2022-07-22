package service

import "github.com/vitelabs/vite-portal/internal/node/types"

func (k Service) PutNode(n types.Node) error {
	return k.store.Upsert(n)
}

func (k Service) DeleteNode(id string) error {
	n, found := k.store.GetById(id)
	if !found {
		return nil
	}
	return k.store.Remove(n.Chain, n.Id)
}
