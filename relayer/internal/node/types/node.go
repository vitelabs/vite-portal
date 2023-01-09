package types

import (
	"errors"
	"fmt"

	"github.com/vitelabs/vite-portal/shared/pkg/logger"
)

type Node struct {
	Id         string `json:"id"`
	Chain      string `json:"chain"`
	RpcHttpUrl string `json:"rpcHttpUrl"`
	RpcWsUrl   string `json:"rpcWsUrl"`
}

func (n *Node) IsValid() bool {
	return n != nil && n.Id != "" && n.Chain != "" && n.RpcHttpUrl != "" && n.RpcWsUrl != ""
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

type GetNodesParams struct {
	Chain  string `json:"chain"`
	Offset int    `json:"offset,string,omitempty"`
	Limit  int    `json:"limit,string,omitempty"`
}

type NodeActivity int

const (
	Put    NodeActivity = 0
	Delete NodeActivity = 1
)

type NodeActivityEntry struct {
	Put    int64
	Delete int64
}
