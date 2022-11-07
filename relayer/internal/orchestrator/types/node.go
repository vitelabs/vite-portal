package types

import (
	sharedtypes "github.com/vitelabs/vite-portal/shared/pkg/types"
)

type Node struct {
	Id            string                 `json:"id"`
	Name          string                 `json:"name,omitempty"`
	Chain         string                 `json:"chain,omitempty"`
	Version       string                 `json:"version,omitempty"`
	Commit        string                 `json:"commit,omitempty"`
	RewardAddress string                 `json:"rewardAddress,omitempty"`
	Transport     string                 `json:"transport,omitempty"`
	RemoteAddress string                 `json:"remoteAddress,omitempty"`
	ClientIp      string                 `json:"clientIp,omitempty"`
	HTTPort       int                    `json:"httpPort,omitempty"`
	WSPort        int                    `json:"wsPort,omitempty"`
	Status        int                    `json:"status,omitempty"`
	LastUpdate    sharedtypes.Int64      `json:"lastUpdate"`
	DelayTime     sharedtypes.Int64      `json:"delayTime"`
	LastBlock     sharedtypes.ChainBlock `json:"lastBlock"`
	HTTPInfo      sharedtypes.HTTPInfo   `json:"httpInfo"`
}
