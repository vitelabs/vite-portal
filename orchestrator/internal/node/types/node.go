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
	Name          string                      `json:"name"`
	Chain         string                      `json:"chain"`
	Version       string                      `json:"version"`
	Commit        string                      `json:"commit"`
	RewardAddress string                      `json:"rewardAddress"`
	Transport     string                      `json:"transport"`
	RemoteAddress string                      `json:"remoteAddress"`
	ClientIp      string                      `json:"clientIp"`
	Status        int                         `json:"status"`
	LastUpdate    sharedtypes.Int64           `json:"lastUpdate"`
	DelayTime     sharedtypes.Int64           `json:"delayTime"`
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
