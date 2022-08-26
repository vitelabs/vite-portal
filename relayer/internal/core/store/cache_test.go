package store

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/vitelabs/vite-portal/relayer/internal/core/types"
	"github.com/vitelabs/vite-portal/shared/pkg/util/idutil"
)

func TestCacheStore(t *testing.T) {
	t.Parallel()
	sessionMaxDuration := int64(1000)
	capacity := 10
	store := NewCacheStore(capacity)

	s := createSession()
	existing, found := store.GetSession(s.Header, sessionMaxDuration)
	require.False(t, found)
	require.Empty(t, existing)

	store.SetSession(s)
	existing, found = store.GetSession(s.Header, sessionMaxDuration)
	require.True(t, found)
	require.NotEmpty(t, existing)
	require.Equal(t, s.Header.Chain, existing.Header.Chain)

	store.ClearSessions()
	existing, found = store.GetSession(s.Header, sessionMaxDuration)
	require.False(t, found)
	require.Empty(t, existing)
}

func TestExpired(t *testing.T) {
	t.Parallel()
	sessionMaxDuration := int64(50)
	capacity := 10
	store := NewCacheStore(capacity)

	s := createSession()
	existing, found := store.GetSession(s.Header, sessionMaxDuration)
	require.False(t, found)
	require.Empty(t, existing)

	store.SetSession(s)
	existing, found = store.GetSession(s.Header, sessionMaxDuration)
	require.True(t, found)
	require.NotEmpty(t, existing)
	require.Equal(t, s.Header.Chain, existing.Header.Chain)

	require.Less(t, time.Now().UnixMilli() - sessionMaxDuration, s.Timestamp)
	time.Sleep(time.Duration(sessionMaxDuration) * time.Millisecond + time.Millisecond)
	require.Greater(t, time.Now().UnixMilli() - sessionMaxDuration, s.Timestamp)

	existing, found = store.GetSession(s.Header, sessionMaxDuration)
	require.False(t, found)
	require.Empty(t, existing)
}

func createSession() types.Session {
	header := types.SessionHeader{
		IpAddress: idutil.NewGuid(),
		Chain:     "chain1",
	}
	return types.Session{
		Key: header.HashString(),
		Timestamp: time.Now().UnixMilli(),
		Header: header,
	}
}
