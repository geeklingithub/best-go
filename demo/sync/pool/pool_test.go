package pool

import (
	"fmt"
	"runtime/debug"
	"sync"
	"testing"
)

func TestPoolNew(t *testing.T) {
	// disable GC so we can control when it happens.
	defer debug.SetGCPercent(debug.SetGCPercent(-1))

	i := 0
	p := sync.Pool{
		New: func() any {
			i++
			return i
		},
	}

	for i := 0; i < 100; i++ {
		p.Put(-1)
	}
	fmt.Println()
}
