package commonutil

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWaitForChan(t *testing.T) {
	t.Parallel()
	c := make(chan string)
	timeout := 100 * time.Millisecond
	startTime := time.Now()
	go func() {
		sleepDuration := time.Duration(timeout.Milliseconds() / 4 * int64(time.Millisecond))
		time.Sleep(sleepDuration)
		c<-"test1"
		time.Sleep(sleepDuration)
		c<-"test2"
	}()
	WaitForChan(timeout, c, func(result string) bool {
		return result == "test2"
	})
	endTime := time.Now()
	assert.Less(t, endTime.UnixNano()-startTime.UnixNano(), timeout.Nanoseconds())
}

func TestWaitForChanTimeout(t *testing.T) {
	t.Parallel()
	c := make(chan string)
	timeout := 100 * time.Millisecond
	startTime := time.Now()
	go func() {
		sleepDuration := time.Duration(timeout.Milliseconds() / 4 * int64(time.Millisecond))
		time.Sleep(sleepDuration)
		c<-"test1"
		time.Sleep(sleepDuration)
		c<-"test1"
	}()
	WaitForChan(timeout, c, func(result string) bool {
		return result == "test2"
	})
	endTime := time.Now()
	assert.Greater(t, endTime.UnixNano()-startTime.UnixNano(), timeout.Nanoseconds())
}

type testValue struct {
	id int
}

func TestIsZero(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		v testValue
		expected bool
	}{
		{
			name: "Test zero",
			v: testValue{},
			expected: true,
		},
		{
			name: "Test not zero",
			v: testValue{
				id: 1,
			},
			expected: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := IsZero(tc.v)
			assert.Equal(t, tc.expected, actual)
		})
	}
}