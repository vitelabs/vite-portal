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