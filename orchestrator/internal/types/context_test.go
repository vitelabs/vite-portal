package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	sharedtypes "github.com/vitelabs/vite-portal/shared/pkg/types"
)

func TestProcessedSet(t *testing.T) {
	cfg := NewDefaultConfig()
	require.NoError(t, cfg.Validate())
	c := NewContext(cfg)
	store1, err := c.GetStatusStore(sharedtypes.DefaultSupportedChains[0].Name)
	require.NoError(t, err)
	store2, err := c.GetStatusStore(sharedtypes.DefaultSupportedChains[1].Name)
	require.NoError(t, err)
	set := *store1.ProcessedSet
	assert.NotEmpty(t, set)
	assert.Equal(t, 0, set.Cardinality())
	set.Add("id1")
	assert.Equal(t, 1, set.Cardinality())
	set.Add("id2")
	assert.Equal(t, 2, set.Cardinality())

	set = *store2.ProcessedSet
	assert.Equal(t, 0, set.Cardinality())
	set.Add("id1")
	assert.Equal(t, 1, set.Cardinality())

	set = *store1.ProcessedSet
	assert.Equal(t, 2, set.Cardinality())
	assert.True(t, set.Contains("id1"))
	set.Remove("id1")
	assert.Equal(t, 1, set.Cardinality())
	assert.False(t, set.Contains("id1"))

	set = *store2.ProcessedSet
	assert.Equal(t, 1, set.Cardinality())
	assert.True(t, set.Contains("id1"))
}