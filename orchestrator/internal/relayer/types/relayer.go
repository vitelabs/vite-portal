package types

import (
	"errors"
	"fmt"

	"github.com/vitelabs/vite-portal/shared/pkg/logger"
	"github.com/vitelabs/vite-portal/shared/pkg/rpc"
)

type Relayer struct {
	Id        string
	Version   string
	RpcClient *rpc.Client
	PeerInfo  rpc.PeerInfo
}

func (r *Relayer) IsValid() bool {
	return r != nil && r.Id != "" && r.RpcClient != nil
}

func (r *Relayer) Validate() error {
	if !r.IsValid() {
		msg := "relayer is invalid"
		err := errors.New(msg)
		logger.Logger().Error().Err(err).Str("relayer", fmt.Sprintf("%#v", r)).Msg(msg)
		return err
	}
	return nil
}
