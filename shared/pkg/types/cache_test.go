package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	capacity := 10
	testCache := NewCache(capacity)

	require.Equal(t, capacity, testCache.Capacity())
	require.Equal(t, 0, testCache.Len())

	e := testEntry{
		value: "test1234",
	}
	key := "1"
	testCache.Add(key, e)
	testCache.Add(key, e)
	require.Equal(t, 1, testCache.Len())

	actual1 := getTestEntry(t, testCache, key)
	require.Equal(t, e, actual1)

	e.value = "1234"
	require.NotEqual(t, e.value, actual1.value)
	actual2 := getTestEntry(t, testCache, key)
	require.NotEqual(t, e, actual2)

	testCache.Add("2", testEntry{
		value: "4321test",
	})
	require.Equal(t, 2, testCache.Len())
	testCache.Remove("2")
	getTestEntry(t, testCache, key)
	require.Equal(t, 1, testCache.Len())

	testCache.Purge()
	val, ok := testCache.Get(key)
	require.False(t, ok)
	require.Empty(t, val)
	require.Equal(t, capacity, testCache.Capacity())
	require.Equal(t, 0, testCache.Len())
}

type testEntry struct {
	value string
}

func getTestEntry(t *testing.T, cache *Cache, key string) testEntry {
	val, ok := cache.Get(key)
	require.True(t, ok)

	actual, _ := val.(testEntry)
	return actual
}