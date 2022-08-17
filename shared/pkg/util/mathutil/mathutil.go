package mathutil

import "github.com/vitelabs/vite-portal/shared/pkg/types/constraints"

func Min[T constraints.Ordered](a, b T) T {
	if a < b {
			return a
	}
	return b
}