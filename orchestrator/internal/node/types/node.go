package types

import (
	"errors"
	"fmt"

	"github.com/vitelabs/vite-portal/shared/pkg/logger"
	"github.com/vitelabs/vite-portal/shared/pkg/rpc"
)

type Node struct {
	Id            string
	Name          string
	Chain         string
	Version       string
	Commit        string
	RewardAddress string
	RpcClient     *rpc.Client
	PeerInfo      rpc.PeerInfo
}

func (n *Node) IsValid() bool {
	return n != nil && n.Id != "" && n.Chain != "" && n.RpcClient != nil && n.PeerInfo.RemoteAddr != ""
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
