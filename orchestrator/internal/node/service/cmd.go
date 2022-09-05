package service

import (
	"errors"
	"time"

	"github.com/vitelabs/vite-portal/shared/pkg/rpc"
)

func (s *Service) HandleConnect(timeout time.Duration, c *rpc.Client, peerInfo rpc.PeerInfo) (id string, err error) {
	return "", errors.New("not implemented yet")
}