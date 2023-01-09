package generics

import (
	mapset "github.com/deckarep/golang-set/v2"
)

func FilterDuplicates[T comparable](vals ...T) []T {
	s := mapset.NewSet[T]()
	for _, val := range vals {
		s.Add(val)
	}
	return s.ToSlice()
}