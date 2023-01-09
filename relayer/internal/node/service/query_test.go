package service

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	nodestore "github.com/vitelabs/vite-portal/relayer/internal/node/store"
	"github.com/vitelabs/vite-portal/relayer/internal/util/testutil"
)

func TestGetNodes_Empty(t *testing.T) {
	t.Parallel()
	store := nodestore.NewMemoryStore()
	service := NewService(store)

	result, err := service.GetNodes("chain1", 0, 0)

	assert.Nil(t, err)
	assert.Equal(t, 0, result.Offset)
	assert.Equal(t, 0, result.Limit)
	assert.Equal(t, 0, result.Total)
	assert.Equal(t, 0, len(result.Entries))
}

func TestGetNodes(t *testing.T) {
	t.Parallel()
	store := nodestore.NewMemoryStore()
	service := NewService(store)

	chain := "chain1"
	total := 5
	limit := 2

	for i := 0; i < total; i++ {
		node := testutil.NewNode(chain)
		node.Id = fmt.Sprintf("%d", i)
		service.store.Upsert(node)
	}

	result, err := service.GetNodes(chain, 0, limit)

	assert.Nil(t, err)
	assert.Equal(t, 0, result.Offset)
	assert.Equal(t, limit, result.Limit)
	assert.Equal(t, total, result.Total)
	assert.Equal(t, 2, len(result.Entries))
	assert.Equal(t, "0", result.Entries[0].Id)
	assert.Equal(t, "1", result.Entries[1].Id)

	result, err = service.GetNodes(chain, 2, limit)

	assert.Nil(t, err)
	assert.Equal(t, 2, result.Offset)
	assert.Equal(t, limit, result.Limit)
	assert.Equal(t, total, result.Total)
	assert.Equal(t, 2, len(result.Entries))
	assert.Equal(t, "2", result.Entries[0].Id)
	assert.Equal(t, "3", result.Entries[1].Id)

	result, err = service.GetNodes(chain, 4, limit)

	assert.Nil(t, err)
	assert.Equal(t, 4, result.Offset)
	assert.Equal(t, limit, result.Limit)
	assert.Equal(t, total, result.Total)
	assert.Equal(t, 1, len(result.Entries))
	assert.Equal(t, "4", result.Entries[0].Id)
}
