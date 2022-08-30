package client

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/vitelabs/vite-portal/relayer/internal/types"
)

func TestConnect(t *testing.T) {
	t.Skip()
	c := NewClient("ws://localhost:57331/", time.Duration(types.DefaultRpcTimeout) * time.Millisecond)
	err := c.Connect()
	require.NoError(t, err)
}