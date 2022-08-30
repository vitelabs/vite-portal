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
