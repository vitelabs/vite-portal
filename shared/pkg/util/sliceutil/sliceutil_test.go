package sliceutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContains(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		s []string
		e string
		expected bool
	}{
		{
			name: "Test empty s and e",
			s: []string{},
			e: "",
			expected: false,
		},
		{
			name: "Test empty s",
			s: []string{},
			e: "test",
			expected: false,
		},
		{
			name: "Test no match",
			s: []string{"test1"},
			e: "test",
			expected: false,
		},
		{
			name: "Test case sensitive",
			s: []string{"test"},
			e: "TEST",
			expected: false,
		},
		{
			name: "Test default",
			s: []string{"test1", "test2"},
			e: "test2",
			expected: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := Contains(tc.s, tc.e)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestRemoveAt(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		input []int
		index int
		expected []int
	}{
		{
			name: "Test empty slice",
			input: []int{},
			index: 1,
			expected: []int{},
		},
		{
			name: "Test negative index",
			input: []int{1},
			index: -1,
			expected: []int{1},
		},
		{
			name: "Test out of range 1",
			input: []int{1},
			index: 1,
			expected: []int{1},
		},
		{
			name: "Test out of range 2",
			input: []int{1, 2},
			index: 2,
			expected: []int{1, 2},
		},
		{
			name: "Test remove first",
			input: []int{1, 2, 3},
			index: 0,
			expected: []int{3, 2},
		},
		{
			name: "Test remove middle",
			input: []int{1, 2, 3},
			index: 1,
			expected: []int{1, 3},
		},
		{
			name: "Test remove last",
			input: []int{1, 2, 3},
			index: 2,
			expected: []int{1, 2},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := RemoveAt(tc.input, tc.index)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
