package util

import (
	"context"
)

// WaitOnChan waits to receieve count items from ch and returns with no error
// If ctx is done before enough messages were received returns an error
func WaitOnChan(ctx context.Context, done chan struct{}, count int) error {
	doneChannel := make(chan struct{})

	go func() {
		defer close(doneChannel)
		for i := 0; i < count; i++ {
			<-done
		}
	}()

	select {
	case <-doneChannel:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
