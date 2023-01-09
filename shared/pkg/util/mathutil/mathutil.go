package mathutil

import "github.com/vitelabs/vite-portal/shared/pkg/types/constraints"

func Min[T constraints.Ordered](a, b T) T {
	if a < b {
			return a
	}
	return b
}

// Abs returns the absolute value of x.
func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}