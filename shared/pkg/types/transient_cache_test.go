package types

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/vitelabs/vite-portal/shared/pkg/util/idutil"
)

func TestTransientCache(t *testing.T) {
	t.Parallel()
	sessionMaxDuration := int64(1000)
	capacity := 10
	store := NewTransientCache[TestItem](capacity)

	s := createItem()
	existing, found := store.Get(s.Key, sessionMaxDuration)
	require.False(t, found)
	require.Empty(t, existing)

	store.Set(s.Key, s.Timestamp, s)
	existing, found = store.Get(s.Key, sessionMaxDuration)
	require.True(t, found)
	require.NotEmpty(t, existing)
	require.Equal(t, s.Value, existing.Value)

	store.Clear()
	existing, found = store.Get(s.Key, sessionMaxDuration)
	require.False(t, found)
	require.Empty(t, existing)
}

func TestExpired(t *testing.T) {
	t.Parallel()
	sessionMaxDuration := int64(50)
	capacity := 10
	store := NewTransientCache[TestItem](capacity)

	s := createItem()
	existing, found := store.Get(s.Key, sessionMaxDuration)
	require.False(t, found)
	require.Empty(t, existing)

	store.Set(s.Key, s.Timestamp, s)
	existing, found = store.Get(s.Key, sessionMaxDuration)
	require.True(t, found)
	require.NotEmpty(t, existing)
	require.Equal(t, s.Value, existing.Value)

	require.Less(t, time.Now().UnixMilli() - sessionMaxDuration, s.Timestamp)
	time.Sleep(time.Duration(sessionMaxDuration) * time.Millisecond + time.Millisecond)
	require.Greater(t, time.Now().UnixMilli() - sessionMaxDuration, s.Timestamp)

	existing, found = store.Get(s.Key, sessionMaxDuration)
	require.False(t, found)
	require.Empty(t, existing)
}

type TestItem struct {
	Key string
	Timestamp int64
	Value string
}

func createItem() TestItem {
	return TestItem{
		Key: idutil.NewGuid(),
		Timestamp: time.Now().UnixMilli(),
		Value: idutil.NewGuid(),
	}
}
