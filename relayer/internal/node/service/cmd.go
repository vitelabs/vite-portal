package service

import (
	"github.com/vitelabs/vite-portal/internal/node/types"
)

func (s *Service) PutNode(n types.Node) error {
	err := n.Validate()
	if err != nil {
		return err
	}
	err = s.store.Upsert(n)
	if err != nil {
		return err
	}
	s.updateLastActivityTimestamp(n.Chain, types.Put)
	return nil
}

func (s *Service) DeleteNode(id string) error {
	n, found := s.store.GetById(id)
	if !found {
		return nil
	}
	err := s.store.Remove(n.Chain, n.Id)
	if err != nil {
		return err
	}
	s.updateLastActivityTimestamp(n.Chain, types.Delete)
	return nil
}