package commonutil

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWaitFor(t *testing.T) {
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
	WaitFor(timeout, c, func(result string) bool {
		return result == "test2"
	})
	endTime := time.Now()
	assert.Less(t, endTime.UnixNano()-startTime.UnixNano(), timeout.Nanoseconds())
}

func TestWaitForTimeout(t *testing.T) {
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
	WaitFor(timeout, c, func(result string) bool {
		return result == "test2"
	})
	endTime := time.Now()
	assert.Greater(t, endTime.UnixNano()-startTime.UnixNano(), timeout.Nanoseconds())
}