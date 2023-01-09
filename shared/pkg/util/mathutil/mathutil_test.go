package mathutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMin(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		a int
		b int
		expected int
	}{
		{
			name: "Test a greater than b",
			a: 2,
			b: 1,
			expected: 1,
		},
		{
			name: "Test b greater than a",
			a: 1,
			b: 2,
			expected: 1,
		},
		{
			name: "Test a equals b",
			a: 0,
			b: 0,
			expected: 0,
		},
		{
			name: "Test negative value",
			a: -1,
			b: 0,
			expected: -1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := Min(tc.a, tc.b)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestAbs(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		v int
		expected int
	}{
		{
			name: "Test 0",
			v: 0,
			expected: 0,
		},
		{
			name: "Test 1",
			v: 1,
			expected: 1,
		},
		{
			name: "Test -1",
			v: -1,
			expected: 1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := Abs(tc.v)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
