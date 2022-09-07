package types

import (
	"errors"
	"fmt"

	"github.com/vitelabs/vite-portal/shared/pkg/logger"
	"github.com/vitelabs/vite-portal/shared/pkg/rpc"
	sharedtypes "github.com/vitelabs/vite-portal/shared/pkg/types"
)

type Node struct {
	Id            string               `json:"id"`
	Name          string               `json:"name"`
	Chain         string               `json:"chain"`
	Version       string               `json:"version"`
	Commit        string               `json:"commit"`
	RewardAddress string               `json:"rewardAddress"`
	Transport     string               `json:"transport"`
	RemoteAddress string               `json:"remoteAddress"`
	HTTPInfo      sharedtypes.HTTPInfo `json:"httpInfo"`
	RpcClient     *rpc.Client          `json:"-"`
}

func (n *Node) IsValid() bool {
	return n != nil && n.Id != "" && n.Chain != "" && n.RpcClient != nil && n.RemoteAddress != ""
}

func (n *Node) Validate() error {
	if !n.IsValid() {
		msg := "node is invalid"
		err := errors.New(msg)
		logger.Logger().Error().Err(err).Str("node", fmt.Sprintf("%#v", n)).Msg(msg)
		return err
	}
	return nil
}
