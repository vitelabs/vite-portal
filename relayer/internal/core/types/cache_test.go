package types

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	roottypes "github.com/vitelabs/vite-portal/internal/types"
	"github.com/vitelabs/vite-portal/internal/util/idutil"
)

func TestCache(t *testing.T) {
	sessionMaxDuration := int64(1000)
	capacity := 10
	roottypes.GlobalSessionCache = roottypes.NewCache(capacity)

	s := createSession()
	existing, found := GetSession(s.Header, sessionMaxDuration)
	require.False(t, found)
	require.Empty(t, existing)

	SetSession(s)
	existing, found = GetSession(s.Header, sessionMaxDuration)
	require.True(t, found)
	require.NotEmpty(t, existing)
	require.Equal(t, s.Header.Chain, existing.Header.Chain)

	ClearSessions()
	existing, found = GetSession(s.Header, sessionMaxDuration)
	require.False(t, found)
	require.Empty(t, existing)
}

func TestExpired(t *testing.T) {
	sessionMaxDuration := int64(50)
	capacity := 10
	roottypes.GlobalSessionCache = roottypes.NewCache(capacity)

	s := createSession()
	existing, found := GetSession(s.Header, sessionMaxDuration)
	require.False(t, found)
	require.Empty(t, existing)

	SetSession(s)
	existing, found = GetSession(s.Header, sessionMaxDuration)
	require.True(t, found)
	require.NotEmpty(t, existing)
	require.Equal(t, s.Header.Chain, existing.Header.Chain)

	require.Less(t, time.Now().UnixMilli() - sessionMaxDuration, s.Timestamp)
	time.Sleep(time.Duration(sessionMaxDuration) * time.Millisecond + time.Millisecond)
	require.Greater(t, time.Now().UnixMilli() - sessionMaxDuration, s.Timestamp)

	existing, found = GetSession(s.Header, sessionMaxDuration)
	require.False(t, found)
	require.Empty(t, existing)
}

func createSession() Session {
	return Session{
		Timestamp: time.Now().UnixMilli(),
		Header: SessionHeader{
			IpAddress: idutil.NewGuid(),
			Chain:     "chain1",
		},
	}
}
