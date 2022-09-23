package client

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/vitelabs/vite-portal/relayer/internal/types"
	"github.com/vitelabs/vite-portal/shared/pkg/util/idutil"
)

func TestConnect(t *testing.T) {
	t.Skip()
	c := NewClient("ws://localhost:57331/", types.DefaultJwtSecret, time.Duration(types.DefaultRpcTimeout) * time.Millisecond)
	err := c.Connect(idutil.NewGuid())
	require.NoError(t, err)
}