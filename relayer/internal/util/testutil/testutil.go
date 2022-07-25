package testutil

import (
	"testing"

	"github.com/stretchr/testify/require"
	nodeinterfaces "github.com/vitelabs/vite-portal/internal/node/interfaces"
	nodetypes "github.com/vitelabs/vite-portal/internal/node/types"
	"github.com/vitelabs/vite-portal/internal/util/idutil"
)

func PutNodes(t *testing.T, s nodeinterfaces.ServiceI, chain string, count int) {
	for i := 0; i < count; i++ {
		err := s.PutNode(nodetypes.Node{
			Id:        idutil.NewGuid(),
			Chain:     chain,
			IpAddress: "0.0.0.0",
		})
		require.NoError(t, err)
	}
}