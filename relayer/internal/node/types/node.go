package types

import (
	"errors"
	"fmt"

	"github.com/vitelabs/vite-portal/internal/logger"
)

type Node struct {
	Id            string `json:"id"`
	Chain         string `json:"chain"`
	IpAddress     string `json:"ipAddress"`
	RewardAddress string `json:"rewardAddress"`
}

func (n *Node) IsValid() bool {
	return n != nil && n.Id != "" && n.Chain != "" && n.IpAddress != ""
}

func (n *Node) Validate() error {
	if !n.IsValid() {
		err := errors.New("node is invalid")
		logger.Logger().Error().Err(err).Str("node", fmt.Sprintf("%#v", n))
		return err
	}
	return nil
}

type GetNodesParams struct {
	Chain  string `json:"chain"`
	Offset int    `json:"offset"`
	Limit  int    `json:"limit"`
}
