package collections

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		entry nameObjectEntry[testEntry]
		success bool
	}{
		{
			name: "Test empty name",
			entry: *newNameObjectEntry("", *new(testEntry)),
			success: false,
		},
		{
			name: "Test empty object",
			entry: *newNameObjectEntry("1", *new(testEntry)),
			success: false,
		},
		{
			name: "Test int object",
			entry: *newNameObjectEntry("1", testEntry{id: "1", val: "1234"}),
			success: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := NewNameObjectCollection[testEntry]()

			count := c.Count()

			c.Add(tc.entry.name, tc.entry.obj)

			if !tc.success {
				assert.Equal(t, count, c.Count())
				assert.Empty(t, c.Get(tc.entry.name))
			} else {
				assert.Equal(t, count+1, c.Count())
				assert.Equal(t, tc.entry.obj, c.Get(tc.entry.name))
			}
		})
	}
}

type testEntry struct {
	id string
	val string
}

func TestGet(t *testing.T) {
	t.Parallel()
	c := NewNameObjectCollection[testEntry]()
	actual1 := c.Get("")
	assert.Empty(t, actual1)

	actual2 := c.Get("test1234")
	assert.Empty(t, actual2)
}

func TestGetByIndex(t *testing.T) {
	t.Parallel()
	c := NewNameObjectCollection[testEntry]()
	entry1 := testEntry{id: "1", val: "test1"}
	c.Add(entry1.id, entry1)
	entry2 := testEntry{id: "2", val: "test2"}
	c.Add(entry2.id, entry2)
	
	actual := c.GetByIndex(-1)
	assert.Empty(t, actual)

	actual = c.GetByIndex(2)
	assert.Empty(t, actual)

	actualEntry := c.GetByIndex(0)
	assert.Equal(t, entry1.id, actualEntry.id)

	actualEntry = c.GetByIndex(1)
	assert.Equal(t, entry2.id, actualEntry.id)
}

func TestRemove(t *testing.T) {
	t.Parallel()
	c := NewNameObjectCollection[testEntry]()
	entry1 := testEntry{id: "1", val: "test1"}
	entry2 := testEntry{id: "2", val: "test2"}

	c.Add(entry1.id, entry1)
	assert.Equal(t, 1, c.Count())
	c.Remove(entry1.id)
	assert.Equal(t, 0, c.Count())
	c.Add(entry1.id, entry1)
	c.Add(entry1.id, entry1)
	c.Add(entry2.id, entry2)
	assert.Equal(t, 2, c.Count())
	c.Remove(entry1.id)
	assert.Equal(t, 1, c.Count())
	assert.Empty(t, c.Get(entry1.id))
	assert.NotNil(t, c.Get(entry2.id))
}

func TestSet(t *testing.T) {
	t.Parallel()
	c := NewNameObjectCollection[testEntry]()
	entry1 := testEntry{id: "1", val: "test1"}
	entry2 := testEntry{id: "2", val: "test2"}

	c.Set(entry1.id, entry1)
	assert.Equal(t, 1, c.Count())

	c.Set(entry2.id, entry2)
	assert.Equal(t, 2, c.Count())

	c.Set(entry2.id, entry2)
	assert.Equal(t, 2, c.Count())
}

func TestUpdate(t *testing.T) {
	t.Parallel()
	c := NewNameObjectCollection[testEntry]()
	entry1 := testEntry{id: "1", val: "test1"}
	entry2 := testEntry{id: "2", val: "test2"}

	c.Add(entry1.id, entry1)
	c.Add(entry2.id, entry2)
	assert.Equal(t, 2, c.Count())

	entry1Before := c.Get(entry1.id)
	assert.Equal(t, entry1, entry1Before)
	assert.Equal(t, c.entriesMap[entry1.id].obj, entry1Before)
	assert.Equal(t, c.entriesSlice[0].obj, entry1Before)

	entry1.val = "test1.1"
	assert.NotEqual(t, entry1, entry1Before)

	// Check if object is unchanged
	entry1Unchanged := c.Get(entry1.id)
	assert.NotEqual(t, entry1, entry1Unchanged)
	assert.NotEqual(t, c.entriesMap[entry1.id].obj, entry1)
	assert.NotEqual(t, c.entriesSlice[0].obj, entry1)

	entry1Before.val = "test1.2"
	assert.NotEqual(t, c.entriesMap[entry1.id].obj, entry1Before)
	assert.NotEqual(t, c.entriesSlice[0].obj, entry1Before)

	c.Set(entry1.id, entry1)
	// Check if object has changed
	entry1Changed := c.Get(entry1.id)
	assert.Equal(t, entry1, entry1Changed)
	// Check if object has been updated in map and array
	assert.Equal(t, c.entriesMap[entry1.id].obj, entry1)
	assert.Equal(t, c.entriesSlice[0].obj, entry1)
	// Check if object of entry2 is unchanged
	assert.NotEqual(t, c.entriesMap[entry2.id].obj, entry1)
	assert.NotEqual(t, c.entriesSlice[1].obj, entry1)
}

func TestEnumerator(t *testing.T) {
	t.Parallel()
	c := NewNameObjectCollection[testEntry]()
	entry1 := testEntry{id: "1", val: "test1"}
	entry2 := testEntry{id: "2", val: "test2"}
	c.Add(entry1.id, entry1)
	c.Add(entry2.id, entry2)

	e := c.GetEnumerator()

	c.Remove(entry1.id)
	assert.Equal(t, 1, c.Count())

	current, found := e.Current()
	assert.False(t, found)
	assert.Empty(t, current)

	hasNext := e.MoveNext()
	assert.True(t, hasNext)
	current, found  = e.Current()
	assert.True(t, found)
	assert.Equal(t, entry1, current)

	hasNext = e.MoveNext()
	assert.True(t, hasNext)
	current, found  = e.Current()
	assert.True(t, found)
	assert.Equal(t, entry2, current)

	current.val = "test2.1"
	assert.NotEqual(t, entry2, current)
	entry2After := c.Get(entry2.id)
	assert.Equal(t, entry2, entry2After)

	hasNext = e.MoveNext()
	assert.False(t, hasNext)

	// Reset
	e.Reset()
	hasNext = e.MoveNext()
	assert.True(t, hasNext)
	current, found  = e.Current()
	assert.True(t, found)
	assert.Equal(t, entry1, current)
}
