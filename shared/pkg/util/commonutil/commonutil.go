package commonutil

import (
	"reflect"
	"time"
)

func WaitFor[T any](timeout time.Duration, c chan T, checkFn func(result T) bool) {
	for {
		select {
		case res := <-c:
			if checkFn(res) {
				return
			}
		case <-time.After(timeout):
			return
		}
	}
}

func IsZero(value any) bool {
	return reflect.ValueOf(value).IsZero()
}

func CheckPagination(offset, limit int) (int, int) {
	if offset < 0 {
		offset = 0
	}
	if limit <= 0 || limit > 1000 {
		limit = 1000
	}
	return offset, limit
}

// IsEmpty gets whether the specified object is considered empty or not.
func IsEmpty(object interface{}) bool {
	// get nil case out of the way
	if object == nil {
		return true
	}

	objValue := reflect.ValueOf(object)

	switch objValue.Kind() {
	// collection types are empty when they have no element
	case reflect.Chan, reflect.Map, reflect.Slice:
		return objValue.Len() == 0
	// pointers are empty if nil or if the value they point to is empty
	case reflect.Ptr:
		if objValue.IsNil() {
			return true
		}
		deref := objValue.Elem().Interface()
		return IsEmpty(deref)
	// for all other types, compare against the zero value
	// array types are empty when they match their zero-initialized state
	default:
		zero := reflect.Zero(objValue.Type())
		return reflect.DeepEqual(object, zero.Interface())
	}
}