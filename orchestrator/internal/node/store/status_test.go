package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetProcessedSet(t *testing.T) {
	s := NewStatusStore()
	set := *s.GetProcessedSet("test1")
	assert.NotEmpty(t, set)
	assert.Equal(t, 0, set.Cardinality())
	set.Add("id1")
	assert.Equal(t, 1, set.Cardinality())
	set.Add("id2")
	assert.Equal(t, 2, set.Cardinality())

	set = *s.GetProcessedSet("test2")
	assert.Equal(t, 0, set.Cardinality())
	set.Add("id1")
	assert.Equal(t, 1, set.Cardinality())

	set = *s.GetProcessedSet("test1")
	assert.Equal(t, 2, set.Cardinality())
	assert.True(t, set.Contains("id1"))
	set.Remove("id1")
	assert.Equal(t, 1, set.Cardinality())
	assert.False(t, set.Contains("id1"))

	set = *s.GetProcessedSet("test2")
	assert.Equal(t, 1, set.Cardinality())
	assert.True(t, set.Contains("id1"))
}