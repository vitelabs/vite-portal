package types

import (
	"errors"
	"fmt"
)

type RpcRequest struct {
	Id      int    `json:"id"`
	Jsonrpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  []any  `json:"params"`
}

type RpcResponse[T any] struct {
	Id      int       `json:"id"`
	Jsonrpc string    `json:"jsonrpc"`
	Result  T         `json:"result"`
	Error   *RpcError `json:"error,omitempty"`
}

func (t *RpcResponse[T]) GetError() error {
	if t.Error != nil {
		return errors.New(fmt.Sprintf("error code: %d, message: %s", t.Error.Code, t.Error.Message))
	}
	return nil
}

type RpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type RpcAppInfoResponse struct {
	Id      string `json:"id"`
	Version string `json:"version"`
	Name    string `json:"name"`
}

type RpcViteNodeInfoResponse struct {
	ID                    string  `json:"id"`
	Name                  string  `json:"name"`
	NetID                 int     `json:"netId"`
	Version               int     `json:"version"`
	Port                  int     `json:"port"`
	FilePort              int     `json:"filePort"`
	Address               string  `json:"address"`
	AddressSignature      string  `json:"addressSignature"`
	PeerCount             int     `json:"peerCount"`
	Height                uint64  `json:"height"`
	Nodes                 int     `json:"nodes"`
	Latency               []int64 `json:"latency"` // [0,1,12,24]
	BroadCheckFailedRatio float32 `json:"broadCheckFailedRatio"`
}

type RpcViteProcessInfoResponse struct {
	BuildVersion  string `json:"build_version"`
	CommitVersion string `json:"commit_version"`
	NodeName      string `json:"nodeName"`
	RewardAddress string `json:"rewardAddress"`
	HTTPPort      int    `json:"httpPort"`
	WSPort        int    `json:"wsPort"`
	Pid           int    `json:"pid"`
}

type RpcViteRuntimeInfoResponse struct {
	ReqId              string                        `json:"reqId,omitempty"`
	PeersNum           int                           `json:"peersNum"`
	SnapshotPendingNum int                           `json:"snapshotPendingNum"`
	AccountPendingNum  string                        `json:"accountPendingNum"`
	LatestSnapshot     RpcViteLatestSnapshotResponse `json:"latestSnapshot"`
	UpdateTime         int                           `json:"updateTime"`
	DelayTime          []int                         `json:"delayTime"`
	Producer           string                        `json:"producer,omitempty"`
	SignData           string                        `json:"signData"`
}

type RpcViteLatestSnapshotResponse struct {
	Hash   string `json:"Hash"`
	Height int    `json:"Height"`
	Time   int    `json:"Time"`
}

type RpcViteLatestAccountBlockResponse struct {
	BlockType int    `json:"blockType"`
	Height    string `json:"height"`
}
