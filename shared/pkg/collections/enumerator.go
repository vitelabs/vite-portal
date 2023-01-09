package collections

type EnumeratorI[T any] interface {
	MoveNext() bool
	Current() (curr T, found bool)
	Reset()
}
