package cryptoutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vitelabs/vite-portal/relayer/internal/generics"
	g "github.com/zyedidia/generic"
)

func TestUniqueRandomIntCases(t *testing.T) {
	tests := []struct {
		name     string
		max      int
		n        int
		expected int
	}{
		{
			name:     "Test max=0 and n=0",
			max:      0,
			n:        0,
			expected: 0,
		},
		{
			name:     "Test max=1 and n=0",
			max:      1,
			n:        0,
			expected: 0,
		},
		{
			name:     "Test max=0 and n=1",
			max:      0,
			n:        1,
			expected: 0,
		},
		{
			name:     "Test max=1 and n=2",
			max:      1,
			n:        2,
			expected: 0,
		},
		{
			name:     "Test max=1 and n=1",
			max:      1,
			n:        1,
			expected: 1,
		},
		{
			name:     "Test max=100 and n=100",
			max:      100,
			n:        100,
			expected: 100,
		},
		{
			name:     "Test max=100 and n=10",
			max:      100,
			n:        10,
			expected: 10,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			for i := 0; i < 1000; i++ {
				r := UniqueRandomInt(tc.max, tc.n)
				assert.Equal(t, tc.expected, len(r))
				// assert random numbers are unique
				s := generics.HashsetOf(uint64(len(r)), g.Equals[int], g.HashInt, r...)
				assert.Equal(t, tc.expected, s.Size())	
			}
		})
	}
}
