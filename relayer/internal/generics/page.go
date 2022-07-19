package generics

type GenericPage[T any] struct {
	Entries []T `json:"entries"` //	The page of items for this page. This will be an empty array if there are no results.
	Limit   int `json:"limit"`   // The limit used for this page of items.
	Offset  int `json:"offset"`  // The offset used for this page of items.
	Total   int `json:"total"`   // The total number of items in the collection.
}

func NewGenericPage[T any]() *GenericPage[T] {
	return &GenericPage[T]{
		Entries: make([]T, 0),
		Limit:   0,
		Offset:  0,
		Total:   0,
	}
}
