package service

import (
	"github.com/vitelabs/vite-portal/internal/node/types"
)

func (s Service) PutNode(n types.Node) error {
	err := n.Validate()
	if err != nil {
		return err
	}

	return s.store.Upsert(n)
}

func (s Service) DeleteNode(id string) error {
	n, found := s.store.GetById(id)
	if !found {
		return nil
	}
	return s.store.Remove(n.Chain, n.Id)
}