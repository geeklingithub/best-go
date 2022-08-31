package pool

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestNewOrderPool(t *testing.T) {
	pool := NewOrderPool(100, 100, 5*time.Minute)

	for i := 0; i < 10000; i++ {
		index := i

		_, err := pool.SubmitTask(index, func(ctx context.Context) {
			fmt.Println("count", index, index%100)
		})
		if err != nil {
			return
		}

	}

	time.Sleep(time.Minute)
}
