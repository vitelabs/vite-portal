package service

import (
	"github.com/vitelabs/vite-portal/orchestrator/internal/relayer/types"
	"github.com/vitelabs/vite-portal/shared/pkg/generics"
)

func (s *Service) Get(offset, limit int) (generics.GenericPage[types.Relayer], error) {
	return s.store.GetPaginated(offset, limit)
}
