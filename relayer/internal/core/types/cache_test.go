package types

import (
	"testing"

	"github.com/stretchr/testify/require"
	roottypes "github.com/vitelabs/vite-portal/internal/types"
	"github.com/vitelabs/vite-portal/internal/util/idutil"
)

func TestCache(t *testing.T) {
	capacity := 10
	roottypes.GlobalSessionCache = roottypes.NewCache(capacity)

	s := Session{
		Header: SessionHeader{
			IpAddress: idutil.NewGuid(),
			Chain: "chain1",
		},
	}
	existing, found := GetSession(s.Header)
	require.False(t, found)
	require.Empty(t, existing)

	SetSession(s)
	existing, found = GetSession(s.Header)
	require.True(t, found)
	require.NotEmpty(t, existing)
	require.Equal(t, s.Header.Chain, existing.Header.Chain)

	ClearSessions()
	existing, found = GetSession(s.Header)
	require.False(t, found)
	require.Empty(t, existing)
}