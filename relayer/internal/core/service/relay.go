package service

import (
	coretypes "github.com/vitelabs/vite-portal/internal/core/types"
	"github.com/vitelabs/vite-portal/internal/logger"
	roottypes "github.com/vitelabs/vite-portal/internal/types"
)

// HandleRelay handles a read/write request to one or multiple nodes
func (s *Service) HandleRelay(relay coretypes.Relay) (*coretypes.RelayResponse, roottypes.Error) {
	response, err := relay.Execute()
	if err != nil {
		logger.Logger().Error().Err(err).Msg("could not execute relay")
	}
	res := &coretypes.RelayResponse{
		Response: response,
	}
	// TODO: track relay time and add to metrics
	return res, nil
}
