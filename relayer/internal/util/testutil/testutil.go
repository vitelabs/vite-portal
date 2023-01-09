package testutil

import (
	"testing"

	"github.com/stretchr/testify/require"
	nodeinterfaces "github.com/vitelabs/vite-portal/relayer/internal/node/interfaces"
	nodetypes "github.com/vitelabs/vite-portal/relayer/internal/node/types"
	roottypes "github.com/vitelabs/vite-portal/relayer/internal/types"
	"github.com/vitelabs/vite-portal/shared/pkg/util/idutil"
)

func PutNodes(t *testing.T, s nodeinterfaces.ServiceI, chain string, count int) {
	for i := 0; i < count; i++ {
		err := s.PutNode(NewNode(chain))
		require.NoError(t, err)
	}
}

func NewNode(chain string) nodetypes.Node {
	return nodetypes.Node{
		Id: idutil.NewGuid(),
		Chain: chain,
		RpcHttpUrl: roottypes.DefaultRpcHttpUrl,
		RpcWsUrl: roottypes.DefaultRpcWsUrl,
	}
}