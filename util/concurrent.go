package util

import (
	"sync"
	"time"
)

// WaitWithTimeout wait on the waitgroup and then returns false
// if specified time is out then returns true
func WaitWithTimeout(wg *sync.WaitGroup, timeout time.Duration) bool {
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()

	select {
	case <-c:
		return false
	case <-time.After(timeout):
		return true
	}
}

// WaitOnChanWithTimeout waits to receieve count items from ch and returns false
// or wait untill timeout and returns true
func WaitOnChanWithTimeout(ch chan struct{}, count int, timeout time.Duration) bool {
	doneChannel := make(chan struct{})
	go func() {
		defer close(doneChannel)
		for i := 0; i < count; i++ {
			<-ch
		}
	}()

	select {
	case <-doneChannel:
		return false
	case <-time.After(timeout):
		return true
	}
}
