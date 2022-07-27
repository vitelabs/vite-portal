package service

import (
	"testing"

	"github.com/stretchr/testify/require"
	coretypes "github.com/vitelabs/vite-portal/internal/core/types"
	nodeservice "github.com/vitelabs/vite-portal/internal/node/service"
	nodestore "github.com/vitelabs/vite-portal/internal/node/store"
	"github.com/vitelabs/vite-portal/internal/types"
)

func TestGetActualNodes(t *testing.T) {
	types.GlobalSessionCache = types.NewCache(10)
	store := nodestore.NewMemoryStore()
	nodesvc := nodeservice.NewService(store)
	svc := NewService(nodesvc)
	r := svc.getActualNodes(coretypes.Session{})
	require.Equal(t, 0, len(r))
}