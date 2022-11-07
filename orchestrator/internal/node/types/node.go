package types

import (
	"errors"
	"fmt"

	sharedinterfaces "github.com/vitelabs/vite-portal/shared/pkg/interfaces"
	"github.com/vitelabs/vite-portal/shared/pkg/logger"
	sharedtypes "github.com/vitelabs/vite-portal/shared/pkg/types"
)

type Node struct {
	Id            string                      `json:"id"`
	Name          string                      `json:"name,omitempty"`
	Chain         string                      `json:"chain,omitempty"`
	Version       string                      `json:"version,omitempty"`
	Commit        string                      `json:"commit,omitempty"`
	RewardAddress string                      `json:"rewardAddress,omitempty"`
	Transport     string                      `json:"transport,omitempty"`
	RemoteAddress string                      `json:"remoteAddress,omitempty"`
	ClientIp      string                      `json:"clientIp,omitempty"`
	HTTPort       int                         `json:"httpPort,omitempty"`
	WSPort        int                         `json:"wsPort,omitempty"`
	Status        int                         `json:"status,omitempty"`
	LastUpdate    sharedtypes.Int64           `json:"lastUpdate"`
	DelayTime     sharedtypes.Int64           `json:"delayTime"`
	LastBlock     sharedtypes.ChainBlock      `json:"lastBlock"`
	HTTPInfo      sharedtypes.HTTPInfo        `json:"httpInfo"`
	RpcClient     sharedinterfaces.RpcClientI `json:"-"`
}

func (n *Node) IsValid() bool {
	return n != nil && n.Id != "" && n.Chain != "" && n.RpcClient != nil && n.ClientIp != ""
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
