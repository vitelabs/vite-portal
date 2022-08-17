package sliceutil

// Contains reports whether e is present in s.
func Contains[T comparable](s []T, e T) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}

// RemoveAt removes the element at the specified index without preserving order.
func RemoveAt[T any](s []T, i int) []T {
	// Source: https://github.com/golang/go/wiki/SliceTricks#delete-without-preserving-order
	if i < 0 || i+1 > len(s) || len(s) == 0 {
		return s
	}
	// Replace the element to delete with the one at the end of the slice
	s[i] = s[len(s)-1]
	// Avoid memory leak by assigning a zero value before removing it from the slice
	s[len(s)-1] = *new(T)
	return s[:len(s)-1]
}
