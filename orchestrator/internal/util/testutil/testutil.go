package testutil

import (
	"time"

	nodetypes "github.com/vitelabs/vite-portal/orchestrator/internal/node/types"
	"github.com/vitelabs/vite-portal/shared/pkg/rpc"
	sharedtypes "github.com/vitelabs/vite-portal/shared/pkg/types"
	"github.com/vitelabs/vite-portal/shared/pkg/util/idutil"
)

func NewNode(chain string) nodetypes.Node {
	return nodetypes.Node{
		Id:         idutil.NewGuid(),
		Chain:      chain,
		ClientIp:   idutil.NewGuid(),
		RpcClient:  &rpc.ClientMock{},
		LastUpdate: sharedtypes.Int64(time.Now().UnixMilli()),
	}
}
