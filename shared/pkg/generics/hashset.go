package generics

import (
	g "github.com/zyedidia/generic"
	"github.com/zyedidia/generic/hashset"
)

// TODO: replace with built-in "Of" when new generic version is released
// https://github.com/zyedidia/generic/tree/master/hashset#func-of
func HashsetOf[K comparable](capacity uint64, equals g.EqualsFn[K], hash g.HashFn[K], vals ...K) *hashset.Set[K] {
	s := hashset.New(capacity, equals, hash)
	for _, val := range vals {
		s.Put(val)
	}
	return s
}