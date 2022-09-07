package types

import (
	"errors"
	"fmt"

	"github.com/vitelabs/vite-portal/shared/pkg/logger"
	"github.com/vitelabs/vite-portal/shared/pkg/rpc"
	sharedtypes "github.com/vitelabs/vite-portal/shared/pkg/types"
)

type Relayer struct {
	Id            string               `json:"id"`
	Version       string               `json:"version"`
	Transport     string               `json:"transport"`
	RemoteAddress string               `json:"remoteAddress"`
	RpcClient     *rpc.Client          `json:"-"`
	HTTPInfo      sharedtypes.HTTPInfo `json:"httpInfo"`
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
