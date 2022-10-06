package client

import (
	"fmt"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/stretchr/testify/require"
	"github.com/vitelabs/vite-portal/shared/pkg/logger"
	"github.com/vitelabs/vite-portal/shared/pkg/types"
	"github.com/vitelabs/vite-portal/shared/pkg/util/testutil"
)

var defaultUrl = testutil.DefaultViteMainNodeUrl

func TestSendError(t *testing.T) {
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
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

func TestRetry(t *testing.T) {
	t.Parallel()
	start := time.Now()
	url := "http://localhost:1234"
	c := NewViteClient(url)
	retryCount := 0
	c.client.RequestLogHook = func(logger retryablehttp.Logger, request *http.Request, count int) {
		retryCount++
	}
	require.Equal(t, url, c.url)
	require.Equal(t, uint64(0), c.requestId)
	require.Equal(t, 0, retryCount)

	response := types.RpcResponse[any]{}
	err := c.Send("test1234", nil, &response)
	require.Error(t, err)
	logger.Logger().Info().Msg(fmt.Sprintf("%#v", err))

	elapsed := time.Since(start)
	require.Equal(t, 4, retryCount)
	require.GreaterOrEqual(t, elapsed, 3 * time.Second)
	require.LessOrEqual(t, elapsed, 4 * time.Second)
}
