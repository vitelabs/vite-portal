package service

import (
	"time"

	"github.com/vitelabs/vite-portal/internal/core/types"
	ct "github.com/vitelabs/vite-portal/internal/core/types"
	rt "github.com/vitelabs/vite-portal/internal/types"
)

// HandleDispatch handles the session information for a client request
func (s Service) HandleDispatch(header types.SessionHeader) (*ct.DispatchResponse, rt.Error) {
	header.RequestTime = time.Now().UnixMilli()
	err := header.ValidateHeader()
	if err != nil {
		return nil, err
	}
	return nil, nil
}