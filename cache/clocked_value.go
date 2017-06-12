package cache

import "github.com/la0rg/test_tasks/vector_clock"

type ClockedValue struct {
	*CacheValue
	*vector_clock.VC
}
