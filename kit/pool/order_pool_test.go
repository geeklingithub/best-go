package pool

import (
	"fmt"
	"testing"
	"time"
)

func TestNewOrderPool(t *testing.T) {
	pool := NewOrderPool(100, 100)

	for i := 0; i < 10000; i++ {
		index := i

		pool.submitTask(index, func() {
			fmt.Println("count", index, index%100)
		})
	}

	time.Sleep(time.Minute)
}
