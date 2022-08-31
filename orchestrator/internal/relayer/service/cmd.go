package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/vitelabs/vite-portal/shared/pkg/logger"
	"github.com/vitelabs/vite-portal/shared/pkg/rpc"
)

func (s *Service) Handle(timeout time.Duration, c *rpc.Client, peerInfo rpc.PeerInfo) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	var resp interface{}
	err := c.CallContext(ctx, &resp, "core_getAppInfo")
	if err != nil {
		msg := "calling context failed"
		logger.Logger().Error().Err(err).Msg(msg)
		return errors.New(msg)
	}
	logger.Logger().Debug().Str("resp", fmt.Sprintf("%#v", resp)).Msg("onconnect result")
	return nil
}
