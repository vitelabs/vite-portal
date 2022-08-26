package store

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vitelabs/vite-portal/relayer/internal/core/types"
)

func TestMockCollector(t *testing.T) {
	t.Parallel()
	c := NewHttpCollector("", "", 0)
	err := c.DispatchRelayResult(types.RelayResult{})
	require.Error(t, err)
}