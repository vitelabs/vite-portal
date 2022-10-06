package store

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/vitelabs/vite-portal/shared/pkg/client"
	"github.com/vitelabs/vite-portal/shared/pkg/util/testutil"
)

func TestGetGlobalHeight(t *testing.T) {
	t.Parallel()
	start := time.Now().UnixMilli()
	client := client.NewViteClient(testutil.DefaultViteMainNodeUrl)
	store := NewStatusStore(client)
	require.Equal(t, int64(0), store.globalHeight)
	require.Equal(t, int64(0), store.lastUpdate)

	height := store.GetGlobalHeight()
	lastHeight := store.globalHeight
	lastUpdate := store.lastUpdate
	require.Greater(t, height, int64(0))
	require.Equal(t, height, lastHeight)
	require.GreaterOrEqual(t, lastUpdate, start)

	time.Sleep(time.Millisecond * 5)
	height = store.GetGlobalHeight()
	require.Equal(t, lastHeight, store.globalHeight)
	require.Equal(t, lastUpdate, store.lastUpdate)

	time.Sleep(time.Second * 2)
	height = store.GetGlobalHeight()
	require.Greater(t, store.globalHeight, lastHeight)
	require.Greater(t, store.lastUpdate, lastUpdate)
}