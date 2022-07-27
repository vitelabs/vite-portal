package service

import (
	"testing"

	"github.com/stretchr/testify/require"
	corestore "github.com/vitelabs/vite-portal/internal/core/store"
	coretypes "github.com/vitelabs/vite-portal/internal/core/types"
	nodeservice "github.com/vitelabs/vite-portal/internal/node/service"
	nodestore "github.com/vitelabs/vite-portal/internal/node/store"
)

func TestGetActualNodes(t *testing.T) {
	cache := corestore.NewCacheStore(10)
	store := nodestore.NewMemoryStore()
	nodesvc := nodeservice.NewService(store)
	svc := NewService(cache, nodesvc)
	r := svc.getActualNodes(coretypes.Session{})
	require.Equal(t, 0, len(r))
}