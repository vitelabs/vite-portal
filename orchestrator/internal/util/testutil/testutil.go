package testutil

import (
	nodetypes "github.com/vitelabs/vite-portal/orchestrator/internal/node/types"
	"github.com/vitelabs/vite-portal/shared/pkg/rpc"
	"github.com/vitelabs/vite-portal/shared/pkg/util/idutil"
)

func NewNode(chain string) nodetypes.Node {
	return nodetypes.Node{
		Id: idutil.NewGuid(),
		Chain: chain,
		RpcClient: &rpc.Client{},
		PeerInfo: rpc.PeerInfo{
			RemoteAddr: idutil.NewGuid(),
		},
	}
}