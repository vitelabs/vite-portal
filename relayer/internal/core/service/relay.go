package service

import (
	coretypes "github.com/vitelabs/vite-portal/internal/core/types"
	roottypes "github.com/vitelabs/vite-portal/internal/types"
)

// HandleRelay handles a read/write request to one or multiple nodes
func (s Service) HandleRelay(relay coretypes.Relay) (*coretypes.RelayResponse, roottypes.Error) {
	res := &coretypes.RelayResponse{
		Response: "",
	}
	return res, nil
}
