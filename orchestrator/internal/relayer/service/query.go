package service

import (
	"errors"

	"github.com/vitelabs/vite-portal/orchestrator/internal/relayer/types"
	"github.com/vitelabs/vite-portal/shared/pkg/generics"
	sharedtypes "github.com/vitelabs/vite-portal/shared/pkg/types"
	"github.com/vitelabs/vite-portal/shared/pkg/util/mathutil"
)

func (s *Service) Get(offset, limit int) (generics.GenericPage[types.RelayerEntity], error) {
	total := s.store.Count()
	result := *generics.NewGenericPage[types.RelayerEntity]()
	result.Offset = offset
	result.Limit = limit
	result.Total = total
	if offset >= total {
		return result, nil
	}
	result.Entries = make([]types.RelayerEntity, mathutil.Min(total-result.Offset, limit))
	count := mathutil.Min(result.Offset+result.Limit, total)
	current := 0
	for i := result.Offset; i < count; i++ {
		item, found := s.store.GetByIndex(i)
		if !found {
			return *generics.NewGenericPage[types.RelayerEntity](), errors.New("inconsistent state")
		}
		result.Entries[current] = types.RelayerEntity{
			Id:            item.Id,
			Version:       item.Version,
			Transport:     item.PeerInfo.Transport,
			RemoteAddress: item.PeerInfo.RemoteAddr,
			HttpInfo: sharedtypes.HttpInfo{
				Version:   item.PeerInfo.HTTP.Version,
				UserAgent: item.PeerInfo.HTTP.UserAgent,
				Origin:    item.PeerInfo.HTTP.Origin,
				Host:      item.PeerInfo.HTTP.Host,
			},
		}
		current++
	}
	return result, nil
}
