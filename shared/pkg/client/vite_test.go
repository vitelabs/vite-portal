package client

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vitelabs/vite-portal/shared/pkg/logger"
	"github.com/vitelabs/vite-portal/shared/pkg/types"
)

var defaultUrl = "https://node.vite.net/gvite"

func TestSendError(t *testing.T) {
	c := NewViteClient(defaultUrl)
	require.Equal(t, defaultUrl, c.url)
	require.Equal(t, uint64(0), c.requestId)

	response := types.RpcResponse[any]{}
	err := c.Send("test1234", nil, &response)
	require.Error(t, err)
	require.NotNil(t, response.Error)

	require.Equal(t, uint64(1), c.requestId)
	require.Equal(t, int(-32601), response.Error.Code)
	require.Equal(t, "The method test1234_ does not exist/is not available", response.Error.Message)
	require.Equal(t, "error code: -32601, message: The method test1234_ does not exist/is not available", err.Error())
}

func TestGetSnapshotChainHeight(t *testing.T) {
	c := NewViteClient(defaultUrl)
	require.Equal(t, defaultUrl, c.url)
	require.Equal(t, uint64(0), c.requestId)

	r, err := c.GetSnapshotChainHeight()
	require.NoError(t, err)
	require.NotNil(t, r)
	logger.Logger().Info().Msg(fmt.Sprintf("GetSnapshotChainHeight: %d", r))

	require.Equal(t, uint64(1), c.requestId)
	require.Greater(t, r, int64(0))
}

func TestGetLatestAccountBlock(t *testing.T) {
	c := NewViteClient(defaultUrl)

	r, err := c.GetLatestAccountBlock("vite_0000000000000000000000000000000000000006e82b8ba657")
	require.NoError(t, err)
	require.NotNil(t, r)
	logger.Logger().Info().Msg(fmt.Sprintf("GetLatestAccountBlock: %#v", r))

	height, err := strconv.ParseInt(r.Height, 10, 64)
	require.NoError(t, err)
	require.Equal(t, 4, r.BlockType)
	require.Greater(t, height, int64(0))
}
