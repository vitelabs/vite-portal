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
var defaultTimeout = time.Duration(1000) * time.Millisecond

func TestSendError(t *testing.T) {
	t.Parallel()
	c := NewViteClient(defaultUrl, defaultTimeout)
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
	c := NewViteClient(defaultUrl, defaultTimeout)
	require.Equal(t, defaultUrl, c.url)
	require.Equal(t, uint64(0), c.requestId)

	h, err := c.GetSnapshotChainHeight()
	require.NoError(t, err)
	logger.Logger().Info().Msg(fmt.Sprintf("GetSnapshotChainHeight: %d", h))

	require.Equal(t, uint64(1), c.requestId)
	require.Greater(t, h, int64(0))
}

func TestGetSnapshotChainHeight_Timeout(t *testing.T) {
	t.Parallel()
	c := NewViteClient(defaultUrl, time.Duration(1) * time.Millisecond)
	require.Equal(t, defaultUrl, c.url)
	require.Equal(t, uint64(0), c.requestId)

	h, err := c.GetSnapshotChainHeight()
	require.Error(t, err)
	require.Equal(t, fmt.Sprintf("POST %s giving up after 1 attempt(s): context deadline exceeded", c.url), err.Error())
	logger.Logger().Info().Msg(fmt.Sprintf("GetSnapshotChainHeight: %d", h))

	require.Equal(t, uint64(1), c.requestId)
	require.Equal(t, int64(0), h)
}

func TestGetLatestAccountBlock(t *testing.T) {
	t.Parallel()
	c := NewViteClient(defaultUrl, defaultTimeout)

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
	c := NewViteClient(url, 0)
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
