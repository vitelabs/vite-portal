package generics

type GenericPage[T any] struct {
	Result []T `json:"result"`
	Total  int `json:"total_pages"`
	Page   int `json:"page"`
}
