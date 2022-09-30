package collections

type EnumeratorI[T any] interface {
	MoveNext() bool
	Current() (found bool, curr T)
	Reset()
}
